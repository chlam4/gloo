package plugins

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/gogo/protobuf/types"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/internal/control-plane/filewatcher"
	"github.com/solo-io/gloo/pkg/secretwatcher"
)

type Stage int

const (
	PreInAuth Stage = iota
	InAuth
	PostInAuth
	PreOutAuth
	OutAuth
)

type EnvoyNameForUpstream func(upstreamName string) string

type Dependencies struct {
	SecretRefs []string
	FileRefs   []string
}

type TranslatorPlugin interface {
	GetDependencies(cfg *v1.Config) *Dependencies
}

// Parameters for ProcessUpstream()
type UpstreamPluginParams struct {
	EnvoyNameForUpstream EnvoyNameForUpstream
	Secrets              secretwatcher.SecretMap
	Files                filewatcher.Files
}

type UpstreamPlugin interface {
	TranslatorPlugin
	ProcessUpstream(params *UpstreamPluginParams, in *v1.Upstream, out *envoyapi.Cluster) error
}

// Params for ParseFunctionSpec()
type FunctionPluginParams struct {
	UpstreamType string
	ServiceType  string
}

type FunctionPlugin interface {
	UpstreamPlugin
	// if the FunctionSpec does not belong to this plugin, return nil, nil
	// if the FunctionSpec belongs to this plugin but is not valid, return nil, err
	// if the FunctionSpec belogns to this plugin and is valid, return *Struct, nil
	ParseFunctionSpec(params *FunctionPluginParams, in v1.FunctionSpec) (*types.Struct, error)
}

// Params for ProcessRoute()
type RoutePluginParams struct {
	// some route plugins need to know about the upstream(s) they route to
	Upstreams []*v1.Upstream
}

type RoutePlugin interface {
	TranslatorPlugin
	ProcessRoute(params *RoutePluginParams, in *v1.Route, out *envoyroute.Route) error
}

// Params for HttpFilters()
type FilterPluginParams struct{}

type StagedFilter struct {
	HttpFilter *envoyhttp.HttpFilter
	Stage      Stage
}

type FilterPlugin interface {
	TranslatorPlugin
	HttpFilters(params *FilterPluginParams) []StagedFilter
}
