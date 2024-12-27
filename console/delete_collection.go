package console

import (
	"context"
	"fmt"
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

var deleteCollectionCmd = &cobra.Command{
	Use:   "delete-collection",
	Short: "delete typesense collection",
	Long:  `This subcommand delete typesense collection gracefully`,
	Run:   runDeleteCollection,
}

func init() {
	RootCmd.AddCommand(deleteCollectionCmd)
}

func runDeleteCollection(_ *cobra.Command, _ []string) {
	err := validateCollectionDeletionConfig()
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Printf("Typesense Host: %s\n", config.TypesenseHostForCollectionDeletion())
	fmt.Printf("Typesense API Key: %s\n", config.TypesenseAPIKeyForCollectionDeletion())
	fmt.Printf("Collection Name: %s\n", config.CollectionNameToDelete())
	fmt.Printf("Batch Size: %d\n", config.BatchSizeForCollectionDeletion())
	fmt.Printf("Excluded Fields: %s\n", strings.Join(config.ExcludedFieldsForCollectionDeletion(), ","))
	fmt.Print("Do you want to proceed with these config? (yes/no): ")

	var confirmation string
	_, err = fmt.Scanln(&confirmation)
	switch {
	case err != nil:
		log.Error(err)
		return
	case confirmation != "yes":
		log.Println("Delete collection operation cancelled.")
		return
	}

	var (
		ctx      = context.TODO()
		tsClient = newTypesenseClient(config.TypesenseHostForCollectionDeletion(), config.TypesenseAPIKeyForCollectionDeletion())
	)

	logger := log.WithFields(log.Fields{
		"context":          utils.DumpIncomingContext(ctx),
		"sourceCollection": config.CollectionNameToDelete(),
	})

	for {
		searchParams := buildSearchParamsForCollectionDeletion()
		logger = logger.WithField("searchParams", utils.Dump(searchParams))

		searchResult, err := tsClient.SearchCollectionWithResponse(ctx, config.CollectionNameToDelete(), searchParams)
		switch {
		case err != nil:
			logger.Error(err)
			return
		case isTypesenseErrorResponse(searchResult):
			err = dumpTypesenseSearchResponseError(searchResult)
			logger.Error(err)
			return
		case len(*searchResult.JSON200.Hits) <= 0:
			goto DeleteCollection
		}

		var ids []string
		for _, item := range *searchResult.JSON200.Hits {
			doc := *item.Document
			docID, ok := doc["id"].(string)
			if !ok {
				continue
			}

			ids = append(ids, docID)
		}

		logger.Infof("start deleting documents with ids: %s", fmt.Sprintf("id:=[%s]", strings.Join(ids, ",")))

		resp, err := tsClient.DeleteDocumentsWithResponse(ctx, config.CollectionNameToDelete(), &typesenseAPI.DeleteDocumentsParams{
			BatchSize: typesensePtr.Int(config.BatchSizeForCollectionDeletion()),
			FilterBy:  typesensePtr.String(fmt.Sprintf("id:=[%s]", strings.Join(ids, ","))),
		})
		switch {
		case err != nil:
			logger.Error(err)
			return
		case resp.StatusCode() != http.StatusOK:
			err = dumpTypesenseError(resp.JSON404)
			logger.Error(err)
			return
		}

		time.Sleep(config.SleepIntervalForCollectionDeletion())
	}

DeleteCollection:
	logger.Infof("start deleting collection: %s", config.CollectionNameToDelete())

	resp, err := tsClient.DeleteCollectionWithResponse(ctx, config.CollectionNameToDelete())
	switch {
	case err != nil:
		logger.Error(err)
		return
	case resp.StatusCode() != http.StatusOK:
		err = dumpTypesenseError(resp.JSON404)
		logger.Error(err)
		return
	}

	log.Printf("Collection %s successfully deleted", config.CollectionNameToDelete())
}

func buildSearchParamsForCollectionDeletion() (searchParams *typesenseAPI.SearchCollectionParams) {
	searchParams = &typesenseAPI.SearchCollectionParams{
		Q:             typesensePtr.String("*"),
		PerPage:       typesensePtr.Int(config.BatchSizeForCollectionDeletion()),
		Page:          typesensePtr.Int(1),
		IncludeFields: typesensePtr.String("id"),
	}

	if len(config.ExcludedFieldsForCollectionDeletion()) > 0 {
		searchParams.ExcludeFields = typesensePtr.String(strings.Join(config.ExcludedFieldsForCollectionDeletion(), ","))
	}

	if config.SorterForCollectionDeletion() != "" {
		searchParams.SortBy = typesensePtr.String(config.SorterForCollectionDeletion())
	}

	return
}

func validateCollectionDeletionConfig() error {
	parsedURL, err := url.Parse(config.TypesenseHostForCollectionDeletion())
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("invalid typesense host URL: %s", config.TypesenseHostForCollectionDeletion())
	}

	switch {
	case config.TypesenseAPIKeyForCollectionDeletion() == "":
		return fmt.Errorf("delete_collection.typesense.api_key cannot be empty")
	case config.CollectionNameToDelete() == "":
		return fmt.Errorf("delete_collection.collection cannot be empty")
	case config.RestoreBatchSize() <= 0:
		return fmt.Errorf("delete_collection.batch_size must be a positive integer")
	}

	return nil
}
