package cli

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/MokoGuy/paperclip/internal/delivery/cli/output"
	"github.com/MokoGuy/paperclip/internal/domain"
	"github.com/MokoGuy/paperclip/internal/repository/api"
	sqliterepo "github.com/MokoGuy/paperclip/internal/repository/sqlite"
	"github.com/MokoGuy/paperclip/internal/usecase"
)

var jsonFlag bool

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "paperclip",
		Short: "CLI d'exploration Paperless-NGX",
		Long:  "paperCLIp — explore your Paperless-NGX instance from the terminal or as an LLM agent.",
		SilenceUsage: true,
	}

	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "Force JSON output")

	rootCmd.AddCommand(
		newTagsCmd(),
		newTypesCmd(),
		newCorrespondentsCmd(),
		newSyncCmd(),
		newSearchCmd(),
		newContentCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

type deps struct {
	cfg             *domain.Config
	db              *sql.DB
	formatter       *output.Formatter
	taxonomyService *usecase.TaxonomyService
	searchService   *usecase.SearchService
	contentService  *usecase.ContentService
	syncService     *usecase.SyncService
}

func loadAllDeps() (*deps, error) {
	cfg, err := domain.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	_ = domain.EnsureConfigPermissions()

	db, err := sqliterepo.NewConnection()
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	apiClient := api.NewClient(cfg.URL, cfg.Token)
	apiTaxRepo := api.NewTaxonomyRepository(apiClient)
	apiDocRepo := api.NewDocumentRepository(apiClient)
	cacheTaxRepo := sqliterepo.NewTaxonomyRepository(db)
	cacheDocRepo := sqliterepo.NewDocumentRepository(db)
	syncRepo := sqliterepo.NewSyncRepository(db)

	syncSvc := usecase.NewSyncService(apiTaxRepo, apiDocRepo, cacheTaxRepo, cacheDocRepo, syncRepo)
	taxSvc := usecase.NewTaxonomyService(cacheTaxRepo, syncSvc)
	searchSvc := usecase.NewSearchService(cacheDocRepo, cacheTaxRepo, apiDocRepo, syncSvc)
	contentSvc := usecase.NewContentService(apiDocRepo)

	formatter := output.NewFormatter(jsonFlag, cfg.URL)

	return &deps{
		cfg:             cfg,
		db:              db,
		formatter:       formatter,
		taxonomyService: taxSvc,
		searchService:   searchSvc,
		contentService:  contentSvc,
		syncService:     syncSvc,
	}, nil
}
