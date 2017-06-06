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

/* RPC stubs for UsagePlanKey resource provider */

// UsagePlanKeyToken is the type token corresponding to the UsagePlanKey package type.
const UsagePlanKeyToken = tokens.Type("aws:apigateway/usagePlanKey:UsagePlanKey")

// UsagePlanKeyProviderOps is a pluggable interface for UsagePlanKey-related management functionality.
type UsagePlanKeyProviderOps interface {
    Check(ctx context.Context, obj *UsagePlanKey) ([]error, error)
    Create(ctx context.Context, obj *UsagePlanKey) (resource.ID, error)
    Get(ctx context.Context, id resource.ID) (*UsagePlanKey, error)
    InspectChange(ctx context.Context,
        id resource.ID, old *UsagePlanKey, new *UsagePlanKey, diff *resource.ObjectDiff) ([]string, error)
    Update(ctx context.Context,
        id resource.ID, old *UsagePlanKey, new *UsagePlanKey, diff *resource.ObjectDiff) error
    Delete(ctx context.Context, id resource.ID) error
}

// UsagePlanKeyProvider is a dynamic gRPC-based plugin for managing UsagePlanKey resources.
type UsagePlanKeyProvider struct {
    ops UsagePlanKeyProviderOps
}

// NewUsagePlanKeyProvider allocates a resource provider that delegates to a ops instance.
func NewUsagePlanKeyProvider(ops UsagePlanKeyProviderOps) lumirpc.ResourceProviderServer {
    contract.Assert(ops != nil)
    return &UsagePlanKeyProvider{ops: ops}
}

func (p *UsagePlanKeyProvider) Check(
    ctx context.Context, req *lumirpc.CheckRequest) (*lumirpc.CheckResponse, error) {
    contract.Assert(req.GetType() == string(UsagePlanKeyToken))
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

func (p *UsagePlanKeyProvider) Name(
    ctx context.Context, req *lumirpc.NameRequest) (*lumirpc.NameResponse, error) {
    contract.Assert(req.GetType() == string(UsagePlanKeyToken))
    obj, _, err := p.Unmarshal(req.GetProperties())
    if err != nil {
        return nil, err
    }
    if obj.Name == nil || *obj.Name == "" {
        if req.Unknowns[UsagePlanKey_Name] {
            return nil, errors.New("Name property cannot be computed from unknown outputs")
        }
        return nil, errors.New("Name property cannot be empty")
    }
    return &lumirpc.NameResponse{Name: *obj.Name}, nil
}

func (p *UsagePlanKeyProvider) Create(
    ctx context.Context, req *lumirpc.CreateRequest) (*lumirpc.CreateResponse, error) {
    contract.Assert(req.GetType() == string(UsagePlanKeyToken))
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

func (p *UsagePlanKeyProvider) Get(
    ctx context.Context, req *lumirpc.GetRequest) (*lumirpc.GetResponse, error) {
    contract.Assert(req.GetType() == string(UsagePlanKeyToken))
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

func (p *UsagePlanKeyProvider) InspectChange(
    ctx context.Context, req *lumirpc.InspectChangeRequest) (*lumirpc.InspectChangeResponse, error) {
    contract.Assert(req.GetType() == string(UsagePlanKeyToken))
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
        if diff.Changed("key") {
            replaces = append(replaces, "key")
        }
        if diff.Changed("usagePlan") {
            replaces = append(replaces, "usagePlan")
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

func (p *UsagePlanKeyProvider) Update(
    ctx context.Context, req *lumirpc.UpdateRequest) (*pbempty.Empty, error) {
    contract.Assert(req.GetType() == string(UsagePlanKeyToken))
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

func (p *UsagePlanKeyProvider) Delete(
    ctx context.Context, req *lumirpc.DeleteRequest) (*pbempty.Empty, error) {
    contract.Assert(req.GetType() == string(UsagePlanKeyToken))
    id := resource.ID(req.GetId())
    if err := p.ops.Delete(ctx, id); err != nil {
        return nil, err
    }
    return &pbempty.Empty{}, nil
}

func (p *UsagePlanKeyProvider) Unmarshal(
    v *pbstruct.Struct) (*UsagePlanKey, resource.PropertyMap, error) {
    var obj UsagePlanKey
    props := resource.UnmarshalProperties(nil, v, resource.MarshalOptions{RawResources: true})
    return &obj, props, mapper.MapIU(props.Mappable(), &obj)
}

/* Marshalable UsagePlanKey structure(s) */

// UsagePlanKey is a marshalable representation of its corresponding IDL type.
type UsagePlanKey struct {
    Name *string `lumi:"name,optional"`
    Key resource.ID `lumi:"key"`
    UsagePlan resource.ID `lumi:"usagePlan"`
}

// UsagePlanKey's properties have constants to make dealing with diffs and property bags easier.
const (
    UsagePlanKey_Name = "name"
    UsagePlanKey_Key = "key"
    UsagePlanKey_UsagePlan = "usagePlan"
)


