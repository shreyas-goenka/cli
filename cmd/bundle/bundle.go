package bundle

import (
	"github.com/databricks/cli/cmd/bundle/deployment"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bundle",
		Short: "Declarative Automation Bundles let you express data/AI/analytics projects as code.",
		Long: `Declarative Automation Bundles let you express data/AI/analytics projects as code.

Common workflows:
  databricks bundle init default-python      # Initialize new project
  databricks bundle deploy --target dev      # Deploy to development
  databricks bundle run my_job               # Run jobs/pipelines
  databricks bundle deploy --target prod     # Deploy to production

Import existing resources:
  databricks bundle generate job --existing-job-id 123 --key my_job # Generate job configuration
  databricks bundle deployment bind my_job 123                      # Link to an existing job

Online documentation: https://docs.databricks.com/en/dev-tools/bundles/index.html`,
		GroupID: "development",
	}

	initVariableFlag(cmd)
	cmd.AddCommand(newDeployCommand())
	cmd.AddCommand(newDestroyCommand())
	cmd.AddCommand(newRunCommand())
	cmd.AddCommand(newSchemaCommand())
	cmd.AddCommand(newSyncCommand())
	cmd.AddCommand(newValidateCommand())
	cmd.AddCommand(newInitCommand())
	cmd.AddCommand(newSummaryCommand())
	cmd.AddCommand(newGenerateCommand())
	cmd.AddCommand(newDebugCommand())
	cmd.AddCommand(newOpenCommand())
	cmd.AddCommand(newPlanCommand())
	cmd.AddCommand(newConfigRemoteSyncCommand())

	// Bundle Metadata Service (DMS) command groups. The CRUD verbs come from
	// the auto-generated cmd/workspace/bundle tree; we group them under the
	// singular-noun parents and alias to the plural form for discoverability.
	dms := metadataServiceCommands()

	deploymentCmd := deployment.NewDeploymentCommand()
	deploymentCmd.Aliases = append(deploymentCmd.Aliases, "deployments")
	deploymentCmd.AddCommand(renameTo(dms["create-deployment"], "create"))
	deploymentCmd.AddCommand(renameTo(dms["get-deployment"], "get"))
	deploymentCmd.AddCommand(renameTo(dms["list-deployments"], "list"))
	deploymentCmd.AddCommand(renameTo(dms["delete-deployment"], "delete"))
	cmd.AddCommand(deploymentCmd)

	versionCmd := &cobra.Command{
		Use:     "version",
		Aliases: []string{"versions"},
		Short:   "Manage version records in the bundle metadata service.",
	}
	versionCmd.AddCommand(renameTo(dms["create-version"], "create"))
	versionCmd.AddCommand(renameTo(dms["get-version"], "get"))
	versionCmd.AddCommand(renameTo(dms["list-versions"], "list"))
	versionCmd.AddCommand(renameTo(dms["complete-version"], "complete"))
	versionCmd.AddCommand(renameTo(dms["heartbeat"], "heartbeat"))
	cmd.AddCommand(versionCmd)

	resourceCmd := &cobra.Command{
		Use:     "resource",
		Aliases: []string{"resources"},
		Short:   "Read resource records from the bundle metadata service.",
	}
	resourceCmd.AddCommand(renameTo(dms["get-resource"], "get"))
	resourceCmd.AddCommand(renameTo(dms["list-resources"], "list"))
	cmd.AddCommand(resourceCmd)

	operationCmd := &cobra.Command{
		Use:     "operation",
		Aliases: []string{"operations"},
		Short:   "Manage operation records in the bundle metadata service.",
	}
	operationCmd.AddCommand(renameTo(dms["create-operation"], "create"))
	operationCmd.AddCommand(renameTo(dms["get-operation"], "get"))
	operationCmd.AddCommand(renameTo(dms["list-operations"], "list"))
	cmd.AddCommand(operationCmd)

	return cmd
}
