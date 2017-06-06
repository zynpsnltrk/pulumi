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

/* RPC stubs for ApplicationVersion resource provider */

// ApplicationVersionToken is the type token corresponding to the ApplicationVersion package type.
const ApplicationVersionToken = tokens.Type("aws:elasticbeanstalk/applicationVersion:ApplicationVersion")

// ApplicationVersionProviderOps is a pluggable interface for ApplicationVersion-related management functionality.
type ApplicationVersionProviderOps interface {
    Check(ctx context.Context, obj *ApplicationVersion) ([]error, error)
    Create(ctx context.Context, obj *ApplicationVersion) (resource.ID, error)
    Get(ctx context.Context, id resource.ID) (*ApplicationVersion, error)
    InspectChange(ctx context.Context,
        id resource.ID, old *ApplicationVersion, new *ApplicationVersion, diff *resource.ObjectDiff) ([]string, error)
    Update(ctx context.Context,
        id resource.ID, old *ApplicationVersion, new *ApplicationVersion, diff *resource.ObjectDiff) error
    Delete(ctx context.Context, id resource.ID) error
}

// ApplicationVersionProvider is a dynamic gRPC-based plugin for managing ApplicationVersion resources.
type ApplicationVersionProvider struct {
    ops ApplicationVersionProviderOps
}

// NewApplicationVersionProvider allocates a resource provider that delegates to a ops instance.
func NewApplicationVersionProvider(ops ApplicationVersionProviderOps) lumirpc.ResourceProviderServer {
    contract.Assert(ops != nil)
    return &ApplicationVersionProvider{ops: ops}
}

func (p *ApplicationVersionProvider) Check(
    ctx context.Context, req *lumirpc.CheckRequest) (*lumirpc.CheckResponse, error) {
    contract.Assert(req.GetType() == string(ApplicationVersionToken))
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

func (p *ApplicationVersionProvider) Name(
    ctx context.Context, req *lumirpc.NameRequest) (*lumirpc.NameResponse, error) {
    contract.Assert(req.GetType() == string(ApplicationVersionToken))
    obj, _, err := p.Unmarshal(req.GetProperties())
    if err != nil {
        return nil, err
    }
    if obj.Name == nil || *obj.Name == "" {
        if req.Unknowns[ApplicationVersion_Name] {
            return nil, errors.New("Name property cannot be computed from unknown outputs")
        }
        return nil, errors.New("Name property cannot be empty")
    }
    return &lumirpc.NameResponse{Name: *obj.Name}, nil
}

func (p *ApplicationVersionProvider) Create(
    ctx context.Context, req *lumirpc.CreateRequest) (*lumirpc.CreateResponse, error) {
    contract.Assert(req.GetType() == string(ApplicationVersionToken))
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

func (p *ApplicationVersionProvider) Get(
    ctx context.Context, req *lumirpc.GetRequest) (*lumirpc.GetResponse, error) {
    contract.Assert(req.GetType() == string(ApplicationVersionToken))
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

func (p *ApplicationVersionProvider) InspectChange(
    ctx context.Context, req *lumirpc.InspectChangeRequest) (*lumirpc.InspectChangeResponse, error) {
    contract.Assert(req.GetType() == string(ApplicationVersionToken))
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
        if diff.Changed("application") {
            replaces = append(replaces, "application")
        }
        if diff.Changed("versionLabel") {
            replaces = append(replaces, "versionLabel")
        }
        if diff.Changed("sourceBundle") {
            replaces = append(replaces, "sourceBundle")
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

func (p *ApplicationVersionProvider) Update(
    ctx context.Context, req *lumirpc.UpdateRequest) (*pbempty.Empty, error) {
    contract.Assert(req.GetType() == string(ApplicationVersionToken))
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

func (p *ApplicationVersionProvider) Delete(
    ctx context.Context, req *lumirpc.DeleteRequest) (*pbempty.Empty, error) {
    contract.Assert(req.GetType() == string(ApplicationVersionToken))
    id := resource.ID(req.GetId())
    if err := p.ops.Delete(ctx, id); err != nil {
        return nil, err
    }
    return &pbempty.Empty{}, nil
}

func (p *ApplicationVersionProvider) Unmarshal(
    v *pbstruct.Struct) (*ApplicationVersion, resource.PropertyMap, error) {
    var obj ApplicationVersion
    props := resource.UnmarshalProperties(nil, v, resource.MarshalOptions{RawResources: true})
    return &obj, props, mapper.MapIU(props.Mappable(), &obj)
}

/* Marshalable ApplicationVersion structure(s) */

// ApplicationVersion is a marshalable representation of its corresponding IDL type.
type ApplicationVersion struct {
    Name *string `lumi:"name,optional"`
    Application resource.ID `lumi:"application"`
    VersionLabel *string `lumi:"versionLabel,optional"`
    Description *string `lumi:"description,optional"`
    SourceBundle resource.ID `lumi:"sourceBundle"`
}

// ApplicationVersion's properties have constants to make dealing with diffs and property bags easier.
const (
    ApplicationVersion_Name = "name"
    ApplicationVersion_Application = "application"
    ApplicationVersion_VersionLabel = "versionLabel"
    ApplicationVersion_Description = "description"
    ApplicationVersion_SourceBundle = "sourceBundle"
)


