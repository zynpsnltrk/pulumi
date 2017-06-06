// *** WARNING: this file was generated by the Lumi IDL Compiler (LUMIDL). ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package elasticbeanstalk

import (
    "errors"

    pbempty "github.com/golang/protobuf/ptypes/empty"
    pbstruct "github.com/golang/protobuf/ptypes/struct"
    "golang.org/x/net/context"

    "github.com/pulumi/lumi/pkg/resource"
    "github.com/pulumi/lumi/pkg/tokens"
    "github.com/pulumi/lumi/pkg/util/contract"
    "github.com/pulumi/lumi/pkg/util/mapper"
    "github.com/pulumi/lumi/sdk/go/pkg/lumirpc"
)

/* RPC stubs for Application resource provider */

// ApplicationToken is the type token corresponding to the Application package type.
const ApplicationToken = tokens.Type("aws:elasticbeanstalk/application:Application")

// ApplicationProviderOps is a pluggable interface for Application-related management functionality.
type ApplicationProviderOps interface {
    Check(ctx context.Context, obj *Application) ([]error, error)
    Create(ctx context.Context, obj *Application) (resource.ID, error)
    Get(ctx context.Context, id resource.ID) (*Application, error)
    InspectChange(ctx context.Context,
        id resource.ID, old *Application, new *Application, diff *resource.ObjectDiff) ([]string, error)
    Update(ctx context.Context,
        id resource.ID, old *Application, new *Application, diff *resource.ObjectDiff) error
    Delete(ctx context.Context, id resource.ID) error
}

// ApplicationProvider is a dynamic gRPC-based plugin for managing Application resources.
type ApplicationProvider struct {
    ops ApplicationProviderOps
}

// NewApplicationProvider allocates a resource provider that delegates to a ops instance.
func NewApplicationProvider(ops ApplicationProviderOps) lumirpc.ResourceProviderServer {
    contract.Assert(ops != nil)
    return &ApplicationProvider{ops: ops}
}

func (p *ApplicationProvider) Check(
    ctx context.Context, req *lumirpc.CheckRequest) (*lumirpc.CheckResponse, error) {
    contract.Assert(req.GetType() == string(ApplicationToken))
    obj, _, err := p.Unmarshal(req.GetProperties())
    if err == nil {
        if failures, err := p.ops.Check(ctx, obj); err != nil {
            return nil, err
        } else if len(failures) > 0 {
            err = resource.NewCheckError(failures)
        }
    }
    return resource.NewCheckResponse(err), nil
}

func (p *ApplicationProvider) Name(
    ctx context.Context, req *lumirpc.NameRequest) (*lumirpc.NameResponse, error) {
    contract.Assert(req.GetType() == string(ApplicationToken))
    obj, _, err := p.Unmarshal(req.GetProperties())
    if err != nil {
        return nil, err
    }
    if obj.Name == nil || *obj.Name == "" {
        if req.Unknowns[Application_Name] {
            return nil, errors.New("Name property cannot be computed from unknown outputs")
        }
        return nil, errors.New("Name property cannot be empty")
    }
    return &lumirpc.NameResponse{Name: *obj.Name}, nil
}

func (p *ApplicationProvider) Create(
    ctx context.Context, req *lumirpc.CreateRequest) (*lumirpc.CreateResponse, error) {
    contract.Assert(req.GetType() == string(ApplicationToken))
    obj, _, err := p.Unmarshal(req.GetProperties())
    if err != nil {
        return nil, err
    }
    id, err := p.ops.Create(ctx, obj)
    if err != nil {
        return nil, err
    }
    return &lumirpc.CreateResponse{Id: string(id)}, nil
}

func (p *ApplicationProvider) Get(
    ctx context.Context, req *lumirpc.GetRequest) (*lumirpc.GetResponse, error) {
    contract.Assert(req.GetType() == string(ApplicationToken))
    id := resource.ID(req.GetId())
    obj, err := p.ops.Get(ctx, id)
    if err != nil {
        return nil, err
    }
    return &lumirpc.GetResponse{
        Properties: resource.MarshalProperties(
            nil, resource.NewPropertyMap(obj), resource.MarshalOptions{}),
    }, nil
}

func (p *ApplicationProvider) InspectChange(
    ctx context.Context, req *lumirpc.InspectChangeRequest) (*lumirpc.InspectChangeResponse, error) {
    contract.Assert(req.GetType() == string(ApplicationToken))
    id := resource.ID(req.GetId())
    old, oldprops, err := p.Unmarshal(req.GetOlds())
    if err != nil {
        return nil, err
    }
    new, newprops, err := p.Unmarshal(req.GetNews())
    if err != nil {
        return nil, err
    }
    var replaces []string
    diff := oldprops.Diff(newprops)
    if diff != nil {
        if diff.Changed("name") {
            replaces = append(replaces, "name")
        }
        if diff.Changed("applicationName") {
            replaces = append(replaces, "applicationName")
        }
    }
    more, err := p.ops.InspectChange(ctx, id, old, new, diff)
    if err != nil {
        return nil, err
    }
    return &lumirpc.InspectChangeResponse{
        Replaces: append(replaces, more...),
    }, err
}

func (p *ApplicationProvider) Update(
    ctx context.Context, req *lumirpc.UpdateRequest) (*pbempty.Empty, error) {
    contract.Assert(req.GetType() == string(ApplicationToken))
    id := resource.ID(req.GetId())
    old, oldprops, err := p.Unmarshal(req.GetOlds())
    if err != nil {
        return nil, err
    }
    new, newprops, err := p.Unmarshal(req.GetNews())
    if err != nil {
        return nil, err
    }
    diff := oldprops.Diff(newprops)
    if err := p.ops.Update(ctx, id, old, new, diff); err != nil {
        return nil, err
    }
    return &pbempty.Empty{}, nil
}

func (p *ApplicationProvider) Delete(
    ctx context.Context, req *lumirpc.DeleteRequest) (*pbempty.Empty, error) {
    contract.Assert(req.GetType() == string(ApplicationToken))
    id := resource.ID(req.GetId())
    if err := p.ops.Delete(ctx, id); err != nil {
        return nil, err
    }
    return &pbempty.Empty{}, nil
}

func (p *ApplicationProvider) Unmarshal(
    v *pbstruct.Struct) (*Application, resource.PropertyMap, error) {
    var obj Application
    props := resource.UnmarshalProperties(nil, v, resource.MarshalOptions{RawResources: true})
    return &obj, props, mapper.MapIU(props.Mappable(), &obj)
}

/* Marshalable Application structure(s) */

// Application is a marshalable representation of its corresponding IDL type.
type Application struct {
    Name *string `lumi:"name,optional"`
    ApplicationName *string `lumi:"applicationName,optional"`
    Description *string `lumi:"description,optional"`
}

// Application's properties have constants to make dealing with diffs and property bags easier.
const (
    Application_Name = "name"
    Application_ApplicationName = "applicationName"
    Application_Description = "description"
)


