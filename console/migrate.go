package console

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
	"typesense-migration-tools/config"

	"github.com/kumparan/go-utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	typesenseAPI "github.com/typesense/typesense-go/v2/typesense/api"
	typesensePtr "github.com/typesense/typesense-go/v2/typesense/api/pointer"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate typesense documents",
	Long:  `This subcommand migrate data between typesense collections`,
	Run:   runMigrate,
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}

func runMigrate(_ *cobra.Command, _ []string) {
	err := validateMigrationConfig()
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Printf("Source Typesense Host: %s\n", config.MigrationSourceTypesenseHost())
	fmt.Printf("Source Typesense API Key: %s\n", config.MigrationSourceTypesenseAPIKey())
	fmt.Printf("Source Collection Name: %s\n", config.MigrationSourceCollection())
	fmt.Printf("Destination Typesense Host: %s\n", config.MigrationDestinationTypesenseHost())
	fmt.Printf("Destination Typesense API Key: %s\n", config.MigrationDestinationTypesenseAPIKey())
	fmt.Printf("Destination Collection Name: %s\n", config.MigrationDestinationCollection())
	fmt.Printf("Filter: %s\n", config.MigrationFilter())
	fmt.Printf("Sorter: %s\n", config.MigrationSorter())
	fmt.Printf("Included Fields: %s\n", strings.Join(config.MigrationIncludedFields(), ","))
	fmt.Printf("Excluded Fields: %s\n", strings.Join(config.MigrationExcludedFields(), ","))
	fmt.Printf("Batch Size: %d\n", config.MigrationBatchSize())
	fmt.Print("Do you want to proceed with these config? (yes/no): ")

	var confirmation string
	fmt.Scanln(&confirmation)
	if confirmation != "yes" {
		log.Println("Export operation cancelled.")
		return
	}

	var (
		ctx                        = context.TODO()
		sourceTypesenseClient      = newTypesenseClient(config.MigrationSourceTypesenseHost(), config.MigrationSourceTypesenseAPIKey())
		destinationTypesenseClient = newTypesenseClient(config.MigrationDestinationTypesenseHost(), config.MigrationDestinationTypesenseAPIKey())
		page                       = 1
	)

	for {
		searchParams := buildMigrationSearchParams(page)

		logger := log.WithFields(log.Fields{
			"context":                    utils.DumpIncomingContext(ctx),
			"searchParams":               utils.Dump(searchParams),
			"sourceCollection":           config.MigrationSourceCollection(),
			"destinationCollection":      config.MigrationDestinationCollection(),
			"sourceTypesenseHost":        config.MigrationSourceTypesenseHost(),
			"destinationTypesenseHost":   config.MigrationDestinationTypesenseHost(),
			"sourceTypesenseAPIKey":      config.MigrationSourceTypesenseAPIKey(),
			"destinationTypesenseAPIKey": config.MigrationDestinationTypesenseAPIKey(),
		})

		searchResult, err := sourceTypesenseClient.SearchCollectionWithResponse(ctx, config.MigrationSourceCollection(), searchParams)
		switch {
		case err != nil:
			logger.Error(err)
			return
		case isTypesenseErrorResponse(searchResult):
			err = dumpTypesenseSearchResponseError(searchResult)
			logger.Error(err)
			return
		case len(*searchResult.JSON200.Hits) <= 0:
			goto LogSuccess
		}

		var docs []map[string]interface{}
		for _, item := range *searchResult.JSON200.Hits {
			docs = append(docs, *item.Document)
		}

		logger.Infof("start migrating page: %d/%d", page, int64(math.Ceil(float64(*searchResult.JSON200.Found)/float64(config.MigrationBatchSize()))))
		var buf bytes.Buffer
		jsonEncoder := json.NewEncoder(&buf)
		for _, doc := range docs {
			if doc == nil {
				continue
			}

			if err = jsonEncoder.Encode(doc); err != nil {
				return
			}
		}

		if buf.Len() <= 0 {
			return
		}

		var resp *typesenseAPI.ImportDocumentsResponse
		resp, err = destinationTypesenseClient.ImportDocumentsWithBodyWithResponse(ctx, config.MigrationDestinationCollection(), &typesenseAPI.ImportDocumentsParams{
			Action:    typesensePtr.String("upsert"),
			BatchSize: typesensePtr.Int(len(docs)),
		}, "application/octet-stream", &buf)
		switch {
		case err != nil:
			logger.Error(err)
			return
		case resp.StatusCode() != http.StatusOK:
			err = dumpTypesenseError(resp.JSON400, resp.JSON404)
			logger.Error(err)
			return
		}

		time.Sleep(config.MigrationSleepInterval())
		page++
	}

LogSuccess:
	log.Printf("Documents successfully migrated from %s to %s", config.MigrationSourceCollection(), config.MigrationDestinationCollection())
}

func buildMigrationSearchParams(page int) (searchParams *typesenseAPI.SearchCollectionParams) {
	searchParams = &typesenseAPI.SearchCollectionParams{
		Q:       typesensePtr.String("*"),
		PerPage: typesensePtr.Int(config.MigrationBatchSize()),
		Page:    typesensePtr.Int(page),
	}
	if len(config.MigrationSorter()) > 0 {
		searchParams.SortBy = typesensePtr.String(config.MigrationSorter())
	}

	if len(config.MigrationIncludedFields()) > 0 {
		searchParams.IncludeFields = typesensePtr.String(strings.Join(config.MigrationIncludedFields(), ","))
	}

	if len(config.MigrationExcludedFields()) > 0 {
		searchParams.ExcludeFields = typesensePtr.String(strings.Join(config.MigrationExcludedFields(), ","))
	}

	if len(config.MigrationFilter()) > 0 {
		searchParams.FilterBy = typesensePtr.String(config.MigrationFilter())
	}

	return
}

func validateMigrationConfig() error {
	parsedURL, err := url.Parse(config.MigrationSourceTypesenseHost())
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("invalid source typesense host URL: %s", config.MigrationSourceTypesenseHost())
	}

	parsedURL, err = url.Parse(config.MigrationDestinationTypesenseHost())
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("invalid destination typesense host URL: %s", config.MigrationDestinationTypesenseHost())
	}

	switch {
	case config.MigrationSourceTypesenseAPIKey() == "":
		return fmt.Errorf("migration.source.typesense.api_key cannot be empty")
	case config.MigrationDestinationTypesenseAPIKey() == "":
		return fmt.Errorf("migration.destination.typesense.api_key cannot be empty")
	case config.MigrationSourceCollection() == "":
		return fmt.Errorf("migration.source.collection cannot be empty")
	case config.MigrationDestinationCollection() == "":
		return fmt.Errorf("migration.destination.collection cannot be empty")
	case config.MigrationFilter() == "":
		return fmt.Errorf("migration.filter cannot be empty")
	case config.MigrationSorter() == "":
		return fmt.Errorf("migration.sorter cannot be empty")
	case config.MigrationBatchSize() <= 0:
		return fmt.Errorf("migration.batch_size must be a positive integer")
	}

	return nil
}
