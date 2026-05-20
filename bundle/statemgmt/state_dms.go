package statemgmt

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/bundle/direct/dstate"
	"github.com/databricks/cli/libs/tmpdms"
)

// LoadStateFromDMS populates the state DB's resource state from the deployment
// metadata service. Lineage / serial / CLI version come from the local state
// cache (resources.json) -- the server only stores per-resource state, not the
// envelope. If the local cache is absent (e.g. fresh checkout of an existing
// DMS bundle), Database defaults are used (lineage will be generated on the
// next Finalize and persisted locally).
func LoadStateFromDMS(ctx context.Context, b *bundle.Bundle) error {
	if b.DeploymentID == "" {
		return nil
	}

	// Open the local state cache first so the state DB path is set. The
	// envelope fields (lineage / serial / cli_version / state_version) come
	// from this file -- the DMS only stores per-resource state.
	db := &b.DeploymentBundle.StateDB
	_, localPath := b.StateFilenameDirect(ctx)
	if err := db.Open(localPath); err != nil {
		return fmt.Errorf("opening local state: %w", err)
	}

	svc, err := tmpdms.NewDeploymentMetadataAPI(b.WorkspaceClient())
	if err != nil {
		return fmt.Errorf("failed to create metadata service client: %w", err)
	}

	resources, err := svc.ListResources(ctx, tmpdms.ListResourcesRequest{
		DeploymentID: b.DeploymentID,
	})
	if err != nil {
		return fmt.Errorf("failed to list resources from deployment metadata service: %w", err)
	}

	// Populate resource state from the server.
	db.Data.State = make(map[string]dstate.ResourceEntry)

	for _, r := range resources {
		// The DMS stores keys without the "resources." prefix (e.g., "jobs.foo").
		// The state DB expects the full key (e.g., "resources.jobs.foo").
		resourceKey := "resources." + r.ResourceKey

		var stateBytes json.RawMessage
		if r.State != nil {
			stateBytes, err = json.Marshal(r.State)
			if err != nil {
				return fmt.Errorf("marshaling state for %s: %w", resourceKey, err)
			}
		}

		db.Data.State[resourceKey] = dstate.ResourceEntry{
			ID:    r.ResourceID,
			State: stateBytes,
		}
	}

	return nil
}
