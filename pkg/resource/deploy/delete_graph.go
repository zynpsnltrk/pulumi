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
	"sort"

	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/util/logging"
)

// deleteNode is a single node in a delete graph.
type deleteNode struct {
	incomingEdges int // the number of resources that are dependent on this node.
	step *DeleteStep // the delete step for this node.
}

// deleteGraph is a simple edge-counting graph structure that is used to determine groups of nodes that can be deleted
// in parallel with one another. Once created, it is used by repeatedly calling pruneLeaves to remove any leaf
// nodes--which by nature must not be depended on by other resources--from the graph until all nodes have been removed.
// 
// Note that we may over-estimate the number of incoming edges for a node due to resources that were pending deletion
// from an earlier execution: all of our dependency information is recorded by URN only, so pending deletes that share
// a URN with other resources may have false edges from resources that actually depend on newer resources with the
// same URN. Thankfully, this is benign, as we will remove any extra edges naturally during pruning.
type deleteGraph struct {
	nodes []deleteNode
}

func (dg *deleteGraph) Len() int {
	return len(dg.nodes)
}
func (dg *deleteGraph) Less(i, j int) bool {
	return dg.nodes[i].incomingEdges < dg.nodes[j].incomingEdges
}
func (dg *deleteGraph) Swap(i, j int) {
	dg.nodes[i], dg.nodes[j] = dg.nodes[j], dg.nodes[i]
}

// getDeleteDeps returns the URNs considered dependencies of the given resource for purposes of deletion.
func getDeleteDeps(res *resource.State) map[resource.URN]struct{} {
	deps := make(map[resource.URN]struct{})
	for _, d := range res.Dependencies {
		deps[d] = struct{}{}
	}
	if res.Parent != "" {
		deps[res.Parent] = struct{}{}
	}
	return deps
}

// pruneLeaves removes any leaf nodes from the graph and returns removed steps. These steps are safe to execute in
// parallel: because they are leaf nodes, no other nodes may be dependent upon them, including other leaf nodes. Once
// these nodes are removed, the graph is walked incoming edge counts for the remaining nodes are decremented as needed.
func (dg *deleteGraph) pruneLeaves() []*DeleteStep {
	// sort the graph in ascending order by incoming edge count
	sort.Sort(dg)

	// find all nodes with no incoming edges
	var leaves []*DeleteStep
	for _, n := range dg.nodes {
		if n.incomingEdges != 0 {
			break
		}
		leaves = append(leaves, n.step)
	}

	// remove the leaf nodes from the graph
	dg.nodes = dg.nodes[len(leaves):]

	// decrement incoming edges
	for _, l := range leaves {
		deps := getDeleteDeps(l.old)
		for i := range dg.nodes {
			n := &dg.nodes[i]
			if _, ok := deps[n.step.old.URN]; ok {
				n.incomingEdges--
			}
		}
	}

	return leaves
}

// newDeleteGraph creates a new delete graph from the given list of steps. All of these steps must be delete steps, and
// they must be presented in a valid topological order.
func newDeleteGraph(steps []Step) *deleteGraph {
	// add all steps to the node list
	nodes := make([]deleteNode, len(steps))
	for i, s := range steps {
		nodes[i] = deleteNode{step: s.(*DeleteStep)}
	}

	// count incoming edges
	for i := range nodes {
		deps := getDeleteDeps(nodes[i].step.old)
		for j := i+1; j < len(nodes); j++ {
			n := &nodes[j]
			if _, ok := deps[n.step.old.URN]; ok {
				n.incomingEdges++
			}
		}
	}

	// dump the graph
	for i := range nodes {
		logging.V(7).Infof("%v: %v %v", i, nodes[i].step.URN(), nodes[i].incomingEdges)
	}

	return &deleteGraph{nodes: nodes}
}
