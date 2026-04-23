package cli

import (
	"context"
	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/spf13/cobra"

	"github.com/MokoGuy/paperclip/internal/domain"
)

func filterByName[T any](filter string, items []T, getName func(T) string) []T {
	if filter == "" {
		return items
	}

	// 1. Substring match first
	lower := strings.ToLower(filter)
	var matched []T
	for _, item := range items {
		if strings.Contains(strings.ToLower(getName(item)), lower) {
			matched = append(matched, item)
		}
	}
	if len(matched) > 0 {
		return matched
	}

	// 2. Fuzzy match fallback
	names := make([]string, len(items))
	for i, item := range items {
		names[i] = getName(item)
	}
	matches := fuzzy.Find(filter, names)
	result := make([]T, len(matches))
	for i, m := range matches {
		result[i] = items[m.Index]
	}
	return result
}

func newTagsCmd() *cobra.Command {
	var filter string

	cmd := &cobra.Command{
		Use:   "tags",
		Short: "List all tags with document counts",
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := loadAllDeps()
			if err != nil {
				return err
			}
			defer d.db.Close()

			tags, err := d.taxonomyService.ListTags(context.Background())
			if err != nil {
				return err
			}

			tags = filterByName(filter, tags, func(t domain.Tag) string { return t.Name })
			return d.formatter.RenderTags(tags)
		},
	}

	cmd.Flags().StringVar(&filter, "filter", "", "Filter tags by name (substring, then fuzzy)")
	return cmd
}

func newTypesCmd() *cobra.Command {
	var filter string

	cmd := &cobra.Command{
		Use:   "types",
		Short: "List all document types with document counts",
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := loadAllDeps()
			if err != nil {
				return err
			}
			defer d.db.Close()

			types, err := d.taxonomyService.ListDocumentTypes(context.Background())
			if err != nil {
				return err
			}

			types = filterByName(filter, types, func(t domain.DocumentType) string { return t.Name })
			return d.formatter.RenderDocumentTypes(types)
		},
	}

	cmd.Flags().StringVar(&filter, "filter", "", "Filter types by name (substring, then fuzzy)")
	return cmd
}

func newCorrespondentsCmd() *cobra.Command {
	var filter string

	cmd := &cobra.Command{
		Use:   "correspondents",
		Short: "List all correspondents with document counts",
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := loadAllDeps()
			if err != nil {
				return err
			}
			defer d.db.Close()

			correspondents, err := d.taxonomyService.ListCorrespondents(context.Background())
			if err != nil {
				return err
			}

			correspondents = filterByName(filter, correspondents, func(c domain.Correspondent) string { return c.Name })
			return d.formatter.RenderCorrespondents(correspondents)
		},
	}

	cmd.Flags().StringVar(&filter, "filter", "", "Filter correspondents by name (substring, then fuzzy)")
	return cmd
}
