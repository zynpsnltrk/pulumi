// Copyright 2016-2018, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploy

import (
	"context"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/util/contract"
	"github.com/pulumi/pulumi/pkg/util/logging"
)

// planExecutor is a bag of state used to coordinate during the execution of a plan.
type planExecutor struct {
	context context.Context // The context used to cancel plan execution.

	preview bool   // True if this plan is a preview.
	events  Events // An optional event reporter.

	pendingNews *sync.Map // The set of resources awaiting outputs.
}

// apply applies a single step during plan execution.
func (w *planExecutor) apply(step Step) (resource.Status, error) {
	urn := step.URN()

	// If there is a pre-event, raise it.
	var eventctx interface{}
	if w.events != nil {
		var eventerr error
		eventctx, eventerr = w.events.OnResourceStepPre(step)
		if eventerr != nil {
			return resource.StatusOK, errors.Wrapf(eventerr, "pre-step event returned an error")
		}
	}

	// Apply the step.
	logging.V(9).Infof("Applying step %v on %v (preview %v)", step.Op(), urn, w.preview)
	status, err := step.Apply(w.preview)
	logging.V(9).Infof("Applied step %v on %v (preview %v): %v", step.Op(), urn, w.preview, err)

	// If there is no error, proceed to save the state; otherwise, go straight to the exit codepath.
	if err == nil {
		// If we have a state object, and this is a create or update, remember it, as we may need to update it later.
		if step.Logical() && step.New() != nil {
			if prior, has := w.pendingNews.Load(urn); has {
				return resource.StatusOK,
					errors.Errorf("resource '%s' registered twice (%s and %s)", urn, prior.(Step).Op(), step.Op())
			}

			w.pendingNews.Store(urn, step)
		}
	}

	// If there is a post-event, raise it, and in any case, return the results.
	if w.events != nil {
		if eventerr := w.events.OnResourceStepPost(eventctx, step, status, err); eventerr != nil {
			return status, errors.Wrapf(eventerr, "post-step event returned an error")
		}
	}

	return status, err
}

// canceled returns true if this plan has been canceled.
func (ex *planExecutor) canceled() bool {
	return ex.context.Err() != nil
}

// worker is the function executed by each of a plan executor's worker goroutines in order to process step chains off
// the given channel.
func (ex *planExecutor) worker(id int, chains <-chan []Step) error {
	for {
		select {
		case chain, ok := <-chains:
			// If the channel is closed, there is no more work to do. Return now.
			if !ok {
				return nil
			}
			// Otherwise, loop over each step in the chain and apply each in turn after checking for cancellation.
			for _, s := range chain {
				if ex.canceled() {
					return nil
				}

				logging.V(7).Infof("worker[%v].step(%v, %v)", id, s.Op(), s.URN())
				_, err := ex.apply(s)
				if err != nil {
					return err
				}
			}
		case <-ex.context.Done():
			// If the execution was cancelled, just return.
			return nil
		}
	}
}

// registerResourceOutputs processes a single RegisterResourceOutputsEvent
func (ex *planExecutor) registerResourceOutputs(e RegisterResourceOutputsEvent) error {
	// Look up the final state in the pending registration list.
	urn := e.URN()
	step, has := ex.pendingNews.Load(urn)
	contract.Assertf(has, "cannot complete a resource '%v' whose registration isn't pending", urn)
	contract.Assertf(step != nil, "expected a non-nil resource step ('%v')", urn)
	ex.pendingNews.Delete(urn)
	reg := step.(Step)

	// Unconditionally set the resource's outputs to what was provided.  This intentionally overwrites whatever
	// might already be there, since otherwise "deleting" outputs would have no affect.
	outs := e.Outputs()
	logging.V(7).Infof("Registered resource outputs %s: old=#%d, new=#%d", urn, len(reg.New().Outputs), len(outs))
	reg.New().Outputs = e.Outputs()

	// If there is an event subscription for finishing the resource, execute them.
	if ex.events != nil {
		if eventerr := ex.events.OnResourceOutputs(reg); eventerr != nil {
			return errors.Wrapf(eventerr, "resource complete event returned an error")
		}
	}

	// Finally, let the language provider know that we're done processing the event.
	e.Done()
	return nil
}

// drive iterates the given SourceIterator, calculates indepentently-executable step chains, and sends each chain down
// the provided channel for execution.
func (ex *planExecutor) drive(src SourceIterator, stepGen *stepGenerator, chains chan<- []Step) error {
	// Iterate the source in a separate goroutine so that we don't get stuck in src.Next during cancellation.
	type sourceEvent struct {
		event SourceEvent
		err   error
	}
	events := make(chan sourceEvent)
	go func() {
		for {
			// If the execution has been cancelled, just return.
			if ex.canceled() {
				return
			}

			// Otherwise, fetch the next step from the source iterator and send it down the event channel.
			event, err := src.Next()
			select {
			case events <- sourceEvent{event: event, err: err}:
				if event == nil || err != nil {
					close(events)
					return
				}
			case <-ex.context.Done():
				return
			}
		}
	}()

	// Iterate the event channel and generate step chains that can be executed in parallel.
	for {
		var sev sourceEvent
		select {
		case sev = <-events:
		case <-ex.context.Done():
			return nil
		}
		// If there was an error, return it; if there was no event, we're done.
		if sev.err != nil {
			return sev.err
		}
		if sev.event == nil {
			break
		}

		// If we have an event, drive the behavior based on which kind it is.
		switch e := sev.event.(type) {
		case RegisterResourceEvent:
			// If the intent is to register a resource, compute the plan steps necessary to do so.
			steps, steperr := stepGen.GenerateSteps(e)
			if steperr != nil {
				return steperr
			}
			select {
			case chains <- steps:
			case <-ex.context.Done():
				return nil
			}
		case RegisterResourceOutputsEvent:
			// If the intent is to complete a prior resource registration, do so. We do this by just
			// processing the request from the existing state, and do not expose our callers to it.
			if err := ex.registerResourceOutputs(e); err != nil {
				return err
			}
		default:
			contract.Failf("Unrecognized intent from source iterator: %T", e)
		}
	}

	// Create and drain a delete graph.
	deletes := newDeleteGraph(stepGen.GenerateDeletes())
	for deletes.Len() > 0 {
		// Peel off the latest set of leaves.
		steps := deletes.pruneLeaves()
		contract.Assert(len(steps) > 0)

		// Create a wait group and a completion notifier.
		wg, done := &sync.WaitGroup{}, make(chan bool)
		wg.Add(len(steps))
		go func() {
			wg.Wait()
			close(done)
		}()

		// Fire off each step in turn.
		for _, s := range steps {
			s.wg = wg
			select {
			case chains <- []Step{s}:
			case <-ex.context.Done():
				return nil
			}
		}

		// Wait for the deletes to finish before going around again.
		select {
		case <-done:
		case <-ex.context.Done():
			return nil
		}
	}

	return nil
}

// Execute executes the given plan according to the supplied options. If parallel execution is enabled, it uses a pool
// of workers to execute steps concurrently.
func (p *Plan) Execute(ctx context.Context, opts Options, preview bool) (PlanSummary, error) {
	// Create a step generator.
	stepGen := newStepGenerator(p, opts)

	// Obtain a source iterator.
	src, err := p.Source().Iterate(opts)
	if err != nil {
		return stepGen, err
	}
	defer contract.IgnoreClose(src)

	// Set up our executor.
	chains := make(chan []Step)
	pexContext, cancel := context.WithCancel(ctx)
	pex := &planExecutor{
		context: pexContext,
		preview: preview,
		events: opts.Events,
		pendingNews: &sync.Map{},
	}

	// Set up the wait group we'll use to track our workers.
	fanout := opts.DegreeOfParallelism()
	wg := &sync.WaitGroup{}
	wg.Add(fanout)

	// Start our workers.
	logging.V(7).Infof("plan.Execute is starting %v workers...", fanout)
	errors := make([]error, fanout)
	for i := 0; i < fanout; i++ {
		go func(idx int) {
			if err := pex.worker(idx, chains); err != nil {
				errors[idx] = err
				contract.IgnoreClose(src)
				cancel()
			}
			logging.V(7).Infof("worker[%v].finished(%v)", idx, errors[idx])
			wg.Done()
		}(i)
	}

	// Run the driver.
	if err = pex.drive(src, stepGen, chains); err != nil {
		cancel()
		return stepGen, err
	}
	close(chains)

	// Wait for the workers to finish, then accumulate any errors and return.
	wg.Wait()
	for _, e := range errors {
		if e != nil {
			err = multierror.Append(err, e)
		}
	}
	return stepGen, err
}
