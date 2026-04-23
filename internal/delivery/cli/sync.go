package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Force a full cache refresh from the Paperless API",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := loadAllDeps()
			if err != nil {
				return err
			}
			defer deps.db.Close()

			fmt.Fprintln(cmd.ErrOrStderr(), "Syncing...")
			if err := deps.syncService.Sync(context.Background()); err != nil {
				return err
			}
			fmt.Fprintln(cmd.ErrOrStderr(), "Sync complete.")
			return nil
		},
	}
}
