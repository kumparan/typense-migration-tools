package console

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
	"typesense-migration-tools/config"

	"github.com/kumparan/go-utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	typesenseAPI "github.com/typesense/typesense-go/v2/typesense/api"
	typesensePtr "github.com/typesense/typesense-go/v2/typesense/api/pointer"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "backup typesense documents",
	Long:  `This subcommand backup documents from typesense collections`,
	Run:   runBackup,
}

func init() {
	RootCmd.AddCommand(backupCmd)
}

func runBackup(_ *cobra.Command, _ []string) {
	err := validateBackupConfig()
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Printf("Typesense Host: %s\n", config.BackupTypesenseHost())
	fmt.Printf("Typesense API Key: %s\n", config.BackupTypesenseAPIKey())
	fmt.Printf("Collection Name: %s\n", config.BackupCollection())
	fmt.Printf("Folder Path: %s\n", config.BackupFolderPath())
	fmt.Printf("Batch Size: %d\n", config.BackupBatchSize())
	fmt.Printf("Max Docs Per File: %d\n", config.BackupMaxDocsPerFile())
	fmt.Printf("Filter: %s\n", config.BackupFilter())
	fmt.Printf("Sorter: %s\n", config.BackupSorter())
	fmt.Printf("Include Fields: %s\n", strings.Join(config.BackupIncludeFields(), ","))
	fmt.Printf("Exclude Fields: %s\n", strings.Join(config.BackupExcludeFields(), ","))
	fmt.Print("Do you want to proceed with these config? (yes/no): ")

	var confirmation string
	fmt.Scanln(&confirmation)
	if confirmation != "yes" {
		log.Println("Export operation cancelled.")
		return
	}

	var (
		ctx        = context.TODO()
		tsClient   = newTypesenseClient(config.BackupTypesenseHost(), config.BackupTypesenseAPIKey())
		page       = 1
		chunk      = make([]string, 0, config.BackupMaxDocsPerFile())
		chunkCount = 0
	)

	for {
		searchParams := buildBackupSearchParams(page)

		logger := log.WithFields(log.Fields{
			"context":         utils.DumpIncomingContext(ctx),
			"searchParams":    utils.Dump(searchParams),
			"collection":      config.BackupCollection(),
			"typesenseHost":   config.BackupTypesenseHost(),
			"typesenseAPIKey": config.BackupTypesenseAPIKey(),
		})

		searchResult, err := tsClient.SearchCollectionWithResponse(ctx, config.BackupCollection(), searchParams)
		switch {
		case err != nil:
			logger.Error(err)
			return
		case isTypesenseErrorResponse(searchResult):
			err = dumpTypesenseSearchResponseError(searchResult)
			logger.Error(err)
			return
		case len(*searchResult.JSON200.Hits) <= 0:
			goto WriteRemainingData
		}

		for _, item := range *searchResult.JSON200.Hits {
			doc, err := json.Marshal(*item.Document)
			if err != nil {
				logger.Error(err)
				return
			}
			chunk = append(chunk, string(doc))
		}

		chunkTotalLines := len(chunk)
		log.Printf("Chunk progress: %d/%d", chunkTotalLines, config.BackupMaxDocsPerFile())
		if chunkTotalLines >= config.BackupMaxDocsPerFile() {
			if err := writeChunkToFile(chunk, chunkCount); err != nil {
				logger.Error(err)
				return
			}

			chunk = chunk[:0]
			chunkCount++
		}

		time.Sleep(config.BackupSleepInterval())
		page++
	}

WriteRemainingData:
	if len(chunk) > 0 {
		if err := writeChunkToFile(chunk, chunkCount); err != nil {
			log.Error(err)
			return
		}
	}

	log.Printf("Documents successfully exported to folder %s", config.BackupFolderPath())
}

func validateBackupConfig() error {
	parsedURL, err := url.Parse(config.BackupTypesenseHost())
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("invalid typesense host URL: %s", config.BackupTypesenseHost())
	}

	switch {
	case config.BackupTypesenseAPIKey() == "":
		return fmt.Errorf("backup.typesense.api_key cannot be empty")
	case config.BackupCollection() == "":
		return fmt.Errorf("backup.collection cannot be empty")
	case config.BackupFolderPath() == "":
		return fmt.Errorf("backup.folder_path cannot be empty")
	case config.BackupBatchSize() <= 0:
		return fmt.Errorf("backup.batch_size must be a positive integer")
	case config.BackupMaxDocsPerFile() <= 0:
		return fmt.Errorf("backup.max_docs_per_file must be a positive integer")
	}

	return nil
}

func buildBackupSearchParams(page int) (searchParams *typesenseAPI.SearchCollectionParams) {
	searchParams = &typesenseAPI.SearchCollectionParams{
		Q:       typesensePtr.String("*"),
		PerPage: typesensePtr.Int(config.BackupBatchSize()),
		Page:    typesensePtr.Int(page),
	}
	if len(config.BackupSorter()) > 0 {
		searchParams.SortBy = typesensePtr.String(config.BackupSorter())
	}

	if len(config.BackupIncludeFields()) > 0 {
		searchParams.IncludeFields = typesensePtr.String(strings.Join(config.BackupIncludeFields(), ","))
	}

	if len(config.BackupExcludeFields()) > 0 {
		searchParams.ExcludeFields = typesensePtr.String(strings.Join(config.BackupExcludeFields(), ","))
	}

	if len(config.BackupFilter()) > 0 {
		searchParams.FilterBy = typesensePtr.String(config.BackupFilter())
	}

	return
}

func writeChunkToFile(chunk []string, chunkCount int) error {
	filename := fmt.Sprintf("%s/backup_chunk_%d.jsonl", config.BackupFolderPath(), chunkCount)
	logger := log.WithField("filename", filename)

	file, err := os.Create(filename)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	chunkTotalLines := len(chunk)
	for i, line := range chunk {
		if i+1 == chunkTotalLines {
			if _, err := writer.WriteString(line); err != nil {
				logger.Error(err)
				return err
			}
			break
		}
		if _, err := writer.WriteString(line + "\n"); err != nil {
			logger.Error(err)
			return err
		}
	}

	writer.Flush()

	log.Printf("Documents successfully exported to file %s", filename)

	return nil
}
