package console

import (
	"bytes"
	"context"
	"encoding/json"
	"math"
	"net/http"
	"strings"
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
	var (
		ctx      = context.TODO()
		tsClient = newTypesenseClient()
		page     = 1
	)

	for {
		searchParams := buildSearchParams(page)

		logger := log.WithFields(log.Fields{
			"context":               utils.DumpIncomingContext(ctx),
			"searchParams":          utils.Dump(searchParams),
			"sourceCollection":      config.MigrationSourceCollection(),
			"destinationCollection": config.MigrationDestinationCollection(),
		})

		setTypesenseCacheParams(searchParams)
		searchResult, err := tsClient.SearchCollectionWithResponse(ctx, config.MigrationSourceCollection(), searchParams)
		switch {
		case err != nil:
			logger.Error(err)
			return
		case isTypesenseErrorResponse(searchResult):
			err = dumpTypesenseSearchResponseError(searchResult)
			logger.Error(err)
			return
		case len(*searchResult.JSON200.Hits) <= 0:
			return
		}

		var docs []map[string]interface{}
		for _, item := range *searchResult.JSON200.Hits {
			docs = append(docs, *item.Document)
		}

		logger.Infof("start migrating page: %d/%d", page, int64(math.Ceil(float64(*searchResult.JSON200.Found)/float64(config.MigrationSizePerPage()))))
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
		resp, err = tsClient.ImportDocumentsWithBodyWithResponse(ctx, config.MigrationDestinationCollection(), &typesenseAPI.ImportDocumentsParams{
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

		page++
	}

}

func buildSearchParams(page int) (searchParams *typesenseAPI.SearchCollectionParams) {
	searchParams = &typesenseAPI.SearchCollectionParams{
		Q:       typesensePtr.String("*"),
		PerPage: typesensePtr.Int(config.MigrationSizePerPage()),
		Page:    typesensePtr.Int(page),
	}
	if len(config.MigrationSorter()) > 0 {
		searchParams.SortBy = typesensePtr.String(config.MigrationSorter())
	}

	if len(config.MigrationIncludeFields()) > 0 {
		searchParams.IncludeFields = typesensePtr.String(strings.Join(config.MigrationIncludeFields(), ","))
	}

	if len(config.MigrationExcludeFields()) > 0 {
		searchParams.ExcludeFields = typesensePtr.String(strings.Join(config.MigrationExcludeFields(), ","))
	}

	if len(config.MigrationFilter()) > 0 {
		searchParams.FilterBy = typesensePtr.String(config.MigrationFilter())
	}

	return
}
