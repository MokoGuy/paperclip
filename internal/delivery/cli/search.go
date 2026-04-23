package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/MokoGuy/paperclip/internal/domain"
)

func newSearchCmd() *cobra.Command {
	var (
		fromFlag    string
		typeFlag    string
		tagFlag     string
		yearFlag    int
		afterFlag   string
		beforeFlag  string
		recentFlag  int
		idsOnly     bool
		noCache     bool
	)

	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search documents with composable filters",
		Long: `Search documents by combining free-text query with structured filters.
All filters can be combined. Filter names (--from, --type, --tag) are matched
case-insensitively as substrings against cached taxonomy.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := loadAllDeps()
			if err != nil {
				return err
			}
			defer d.db.Close()

			ctx := context.Background()

			if recentFlag > 0 {
				docs, err := d.searchService.ListRecent(ctx, recentFlag, noCache)
				if err != nil {
					return err
				}
				if idsOnly {
					d.formatter.RenderIDs(docs)
					return nil
				}
				return d.formatter.RenderDocuments(docs)
			}

			query := ""
			if len(args) > 0 {
				query = args[0]
			}

			if query == "" && fromFlag == "" && typeFlag == "" && tagFlag == "" && yearFlag == 0 && afterFlag == "" && beforeFlag == "" {
				return fmt.Errorf("provide a query or at least one filter (--from, --type, --tag, --year, --after, --before, --recent)")
			}

			filters := domain.SearchFilters{
				Query:             query,
				CorrespondentName: fromFlag,
				DocumentTypeName:  typeFlag,
				TagName:           tagFlag,
				After:             afterFlag,
				Before:            beforeFlag,
			}
			if yearFlag > 0 {
				filters.Year = &yearFlag
			}

			docs, err := d.searchService.Search(ctx, filters, noCache)
			if err != nil {
				return err
			}

			if idsOnly {
				d.formatter.RenderIDs(docs)
				return nil
			}
			return d.formatter.RenderDocuments(docs)
		},
	}

	cmd.Flags().StringVar(&fromFlag, "from", "", "Filter by correspondent name")
	cmd.Flags().StringVar(&typeFlag, "type", "", "Filter by document type name")
	cmd.Flags().StringVar(&tagFlag, "tag", "", "Filter by tag name")
	cmd.Flags().IntVar(&yearFlag, "year", 0, "Filter by year (e.g. 2025)")
	cmd.Flags().StringVar(&afterFlag, "after", "", "Filter documents created after date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&beforeFlag, "before", "", "Filter documents created before date (YYYY-MM-DD)")
	cmd.Flags().IntVar(&recentFlag, "recent", 0, "Show N most recently added documents")
	cmd.Flags().BoolVar(&idsOnly, "ids-only", false, "Output only document IDs (one per line)")
	cmd.Flags().BoolVar(&noCache, "no-cache", false, "Bypass cache, query API directly")

	return cmd
}
