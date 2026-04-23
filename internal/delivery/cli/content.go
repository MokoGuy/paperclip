package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/MokoGuy/paperclip/internal/delivery/cli/output"
)

func newContentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "content <id> [id...]",
		Short: "Extract text content from one or more documents",
		Long:  "Fetch document text content from the Paperless API (always live, never cached).",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := loadAllDeps()
			if err != nil {
				return err
			}
			defer d.db.Close()

			ids := make([]int, 0, len(args))
			for _, arg := range args {
				id, err := strconv.Atoi(arg)
				if err != nil {
					return fmt.Errorf("invalid document ID %q: %w", arg, err)
				}
				ids = append(ids, id)
			}

			results := d.contentService.GetContents(context.Background(), ids)

			var items []output.ContentItemResponse
			for _, r := range results {
				if r.Err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "Warning: failed to fetch document %d: %v\n", r.ID, r.Err)
					continue
				}
				items = append(items, output.ContentItemResponse{
					ID:      r.ID,
					Title:   r.Title,
					Content: r.Content,
				})
			}

			return d.formatter.RenderContent(items)
		},
	}

	return cmd
}
