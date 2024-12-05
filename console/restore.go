package console

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"typesense-migration-tools/config"

	"github.com/kumparan/go-utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	typesenseAPI "github.com/typesense/typesense-go/v2/typesense/api"
	typesensePtr "github.com/typesense/typesense-go/v2/typesense/api/pointer"
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "restore typesense documents",
	Long:  `This subcommand restore documents from typesense collections`,
	Run:   runRestore,
}

func init() {
	RootCmd.AddCommand(restoreCmd)
}

func runRestore(_ *cobra.Command, _ []string) {
	err := validateRestoreConfig()
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Printf("Typesense Host: %s\n", config.RestoreTypesenseHost())
	fmt.Printf("Typesense API Key: %s\n", config.RestoreTypesenseAPIKey())
	fmt.Printf("Collection Name: %s\n", config.RestoreCollection())
	fmt.Printf("JSONL File Path: %s\n", config.RestoreJSONLFilePath())
	fmt.Printf("Batch Size: %d\n", config.RestoreBatchSize())
	fmt.Print("Do you want to proceed with these config? (yes/no): ")

	var confirmation string
	fmt.Scanln(&confirmation)
	if confirmation != "yes" {
		log.Println("Export operation cancelled.")
		return
	}

	var (
		ctx      = context.TODO()
		tsClient = newTypesenseClient(config.RestoreTypesenseHost(), config.RestoreTypesenseAPIKey())
	)

	logger := log.WithFields(log.Fields{
		"context":         utils.DumpIncomingContext(ctx),
		"collection":      config.RestoreCollection(),
		"jsonlFilePath":   config.RestoreJSONLFilePath(),
		"batchSize":       config.RestoreBatchSize(),
		"typesenseHost":   config.RestoreTypesenseHost(),
		"typesenseAPIKey": config.RestoreTypesenseAPIKey(),
	})

	jsonlData, err := os.ReadFile(config.RestoreJSONLFilePath())
	if err != nil {
		logger.Fatalf("Failed to read JSONL file: %v", err)
	}
	resp, err := tsClient.ImportDocumentsWithBodyWithResponse(ctx, config.RestoreCollection(), &typesenseAPI.ImportDocumentsParams{
		Action:    typesensePtr.String("upsert"),
		BatchSize: typesensePtr.Int(config.RestoreBatchSize()),
	}, "application/jsonl", bytes.NewReader(jsonlData))
	if err != nil {
		logger.Error(err)
		return
	}

	if resp.StatusCode() != http.StatusOK {
		err = dumpTypesenseError(resp.JSON404)
		logger.Error(err)
		return
	}

	logger.Printf("Import response: %s", string(resp.Body))
	log.Printf("Documents successfully imported to %s", config.RestoreCollection())
}

func validateRestoreConfig() error {
	parsedURL, err := url.Parse(config.RestoreTypesenseHost())
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("invalid typesense host URL: %s", config.RestoreTypesenseHost())
	}

	switch {
	case config.RestoreTypesenseAPIKey() == "":
		return fmt.Errorf("restore.typesense.api_key cannot be empty")
	case config.RestoreCollection() == "":
		return fmt.Errorf("restore.collection cannot be empty")
	case config.RestoreJSONLFilePath() == "":
		return fmt.Errorf("restore.jsonl_file_path cannot be empty")
	case config.RestoreBatchSize() <= 0:
		return fmt.Errorf("restore.batch_size must be a positive integer")
	}

	return nil
}
