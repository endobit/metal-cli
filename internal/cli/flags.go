package cli

import "github.com/spf13/cobra"

func AddRenameFlag(cmd *cobra.Command, object string) *string {
	return cmd.Flags().String("model", "", "rename the "+object)
}

func AddZoneFlag(cmd *cobra.Command, object string) *string {
	return cmd.Flags().String("zone", "", "zone for the "+object)
}
