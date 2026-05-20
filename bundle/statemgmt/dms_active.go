package statemgmt

import (
	"context"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/bundle/env"
)

// IsDmsActive returns true if this bundle run should use the deployment
// metadata service for locking, version tracking, and resource state
// management.
//
// A bundle is "DMS-active" when EITHER:
//
//  1. The bundle already has a server-side deployment record (its deployment
//     ID was loaded from managed_service.json during state pull), OR
//  2. DATABRICKS_BUNDLE_MANAGED_STATE is set to a truthy value (case- and
//     spelling-tolerant via env.ManagedStateOptIn -- accepts "1", "true",
//     "yes", "on"), opting a not-yet-deployed bundle into DMS on its first
//     deploy. After the first successful deploy, managed_service.json
//     carries the bundle along and the env var is no longer required.
//
// All gating sites that decide "do we talk to the DMS this run?" should call
// this function. Sites that need to know "do we have a server-side deployment
// ID right now?" (e.g. before calling LoadStateFromDMS) should keep checking
// b.DeploymentID directly -- the two questions are semantically distinct:
// IsDmsActive is true on a first deploy before the deployment record exists,
// while b.DeploymentID only becomes non-empty once the record is created.
func IsDmsActive(ctx context.Context, b *bundle.Bundle) bool {
	if b.DeploymentID != "" {
		return true
	}
	return env.ManagedStateOptIn(ctx)
}
