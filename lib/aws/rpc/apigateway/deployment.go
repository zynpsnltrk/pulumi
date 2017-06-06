// *** WARNING: this file was generated by the Lumi IDL Compiler (LUMIDL). ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package apigateway

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

/* RPC stubs for Deployment resource provider */

// DeploymentToken is the type token corresponding to the Deployment package type.
const DeploymentToken = tokens.Type("aws:apigateway/deployment:Deployment")

// DeploymentProviderOps is a pluggable interface for Deployment-related management functionality.
type DeploymentProviderOps interface {
    Check(ctx context.Context, obj *Deployment) ([]error, error)
    Create(ctx context.Context, obj *Deployment) (resource.ID, error)
    Get(ctx context.Context, id resource.ID) (*Deployment, error)
    InspectChange(ctx context.Context,
        id resource.ID, old *Deployment, new *Deployment, diff *resource.ObjectDiff) ([]string, error)
    Update(ctx context.Context,
        id resource.ID, old *Deployment, new *Deployment, diff *resource.ObjectDiff) error
    Delete(ctx context.Context, id resource.ID) error
}

// DeploymentProvider is a dynamic gRPC-based plugin for managing Deployment resources.
type DeploymentProvider struct {
    ops DeploymentProviderOps
}

// NewDeploymentProvider allocates a resource provider that delegates to a ops instance.
func NewDeploymentProvider(ops DeploymentProviderOps) lumirpc.ResourceProviderServer {
    contract.Assert(ops != nil)
    return &DeploymentProvider{ops: ops}
}

func (p *DeploymentProvider) Check(
    ctx context.Context, req *lumirpc.CheckRequest) (*lumirpc.CheckResponse, error) {
    contract.Assert(req.GetType() == string(DeploymentToken))
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

func (p *DeploymentProvider) Name(
    ctx context.Context, req *lumirpc.NameRequest) (*lumirpc.NameResponse, error) {
    contract.Assert(req.GetType() == string(DeploymentToken))
    obj, _, err := p.Unmarshal(req.GetProperties())
    if err != nil {
        return nil, err
    }
    if obj.Name == nil || *obj.Name == "" {
        if req.Unknowns[Deployment_Name] {
            return nil, errors.New("Name property cannot be computed from unknown outputs")
        }
        return nil, errors.New("Name property cannot be empty")
    }
    return &lumirpc.NameResponse{Name: *obj.Name}, nil
}

func (p *DeploymentProvider) Create(
    ctx context.Context, req *lumirpc.CreateRequest) (*lumirpc.CreateResponse, error) {
    contract.Assert(req.GetType() == string(DeploymentToken))
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

func (p *DeploymentProvider) Get(
    ctx context.Context, req *lumirpc.GetRequest) (*lumirpc.GetResponse, error) {
    contract.Assert(req.GetType() == string(DeploymentToken))
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

func (p *DeploymentProvider) InspectChange(
    ctx context.Context, req *lumirpc.InspectChangeRequest) (*lumirpc.InspectChangeResponse, error) {
    contract.Assert(req.GetType() == string(DeploymentToken))
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
        if diff.Changed("restAPI") {
            replaces = append(replaces, "restAPI")
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

func (p *DeploymentProvider) Update(
    ctx context.Context, req *lumirpc.UpdateRequest) (*pbempty.Empty, error) {
    contract.Assert(req.GetType() == string(DeploymentToken))
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

func (p *DeploymentProvider) Delete(
    ctx context.Context, req *lumirpc.DeleteRequest) (*pbempty.Empty, error) {
    contract.Assert(req.GetType() == string(DeploymentToken))
    id := resource.ID(req.GetId())
    if err := p.ops.Delete(ctx, id); err != nil {
        return nil, err
    }
    return &pbempty.Empty{}, nil
}

func (p *DeploymentProvider) Unmarshal(
    v *pbstruct.Struct) (*Deployment, resource.PropertyMap, error) {
    var obj Deployment
    props := resource.UnmarshalProperties(nil, v, resource.MarshalOptions{RawResources: true})
    return &obj, props, mapper.MapIU(props.Mappable(), &obj)
}

/* Marshalable Deployment structure(s) */

// Deployment is a marshalable representation of its corresponding IDL type.
type Deployment struct {
    Name *string `lumi:"name,optional"`
    RestAPI resource.ID `lumi:"restAPI"`
    Description *string `lumi:"description,optional"`
    ID string `lumi:"id,optional"`
    CreatedDate string `lumi:"createdDate,optional"`
}

// Deployment's properties have constants to make dealing with diffs and property bags easier.
const (
    Deployment_Name = "name"
    Deployment_RestAPI = "restAPI"
    Deployment_Description = "description"
    Deployment_ID = "id"
    Deployment_CreatedDate = "createdDate"
)


