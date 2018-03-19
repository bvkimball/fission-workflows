package integration

import (
	"context"

	"github.com/fission/fission-workflows/cmd/fission-workflows-bundle/bundle"
)

// SetupBundle sets up and runs the workflows-bundle.
//
// By default the bundle runs with all components are enabled, setting up a NATS cluster as the
// backing event store, and internal fnenv and workflow runtime
func SetupBundle(ctx context.Context, opts ...bundle.Options) bundle.Options {
	var bundleOpts bundle.Options
	if len(opts) > 0 {
		bundleOpts = opts[0]
	} else {
		bundleOpts = bundle.Options{
			InternalRuntime:       true,
			InvocationController:  true,
			WorkflowController:    true,
			ApiHttp:               true,
			ApiWorkflowInvocation: true,
			ApiWorkflow:           true,
			ApiAdmin:              true,
		}
	}
	go bundle.Run(ctx, &bundleOpts)
	return bundleOpts
}
