package env

import (
	"context"

	"github.com/databricks/cli/libs/env"
)

// managedStateVariable names the environment variable that opts a bundle
// into server-managed state (the deployment metadata service) on its next
// deploy. Truthy values accepted by env.GetBool ("1", "true", "yes", "on",
// case-insensitive) all opt in; everything else is treated as false.
const managedStateVariable = "DATABRICKS_BUNDLE_MANAGED_STATE"

// ManagedStateOptIn reports whether the DATABRICKS_BUNDLE_MANAGED_STATE env
// var is set to a truthy value, opting a not-yet-deployed bundle into DMS on
// its first deploy. Most callers should instead use
// statemgmt.IsDmsActive, which also picks up bundles whose deployment_id is
// already pinned in managed_service.json -- the env var is only relevant
// before that pin exists.
func ManagedStateOptIn(ctx context.Context) bool {
	v, _ := env.GetBool(ctx, managedStateVariable)
	return v
}
