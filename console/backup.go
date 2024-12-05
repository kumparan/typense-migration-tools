package console

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	fmt.Printf("JSONL File Path: %s\n", config.BackupJSONLFilePath())
	fmt.Printf("Filter: %s\n", config.BackupFilter())
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
		ctx      = context.TODO()
		tsClient = newTypesenseClient(config.BackupTypesenseHost(), config.BackupTypesenseAPIKey())
	)

	exportParams := &typesenseAPI.ExportDocumentsParams{}
	if len(config.BackupIncludeFields()) > 0 {
		exportParams.IncludeFields = typesensePtr.String(strings.Join(config.BackupIncludeFields(), ","))
	}

	if len(config.BackupExcludeFields()) > 0 {
		exportParams.ExcludeFields = typesensePtr.String(strings.Join(config.BackupExcludeFields(), ","))
	}

	if len(config.BackupFilter()) > 0 {
		exportParams.FilterBy = typesensePtr.String(config.BackupFilter())
	}

	logger := log.WithFields(log.Fields{
		"context":      utils.DumpIncomingContext(ctx),
		"exportParams": utils.Dump(exportParams),
		"collection":   config.BackupCollection(),
	})

	resp, err := tsClient.ExportDocumentsWithResponse(ctx, config.BackupCollection(), exportParams)
	if err != nil {
		logger.Error(err)
		return
	}

	if resp.StatusCode() != http.StatusOK {
		err = dumpTypesenseError(resp.JSON404)
		logger.Error(err)
		return
	}

	outputFile, err := os.Create(config.BackupJSONLFilePath())
	if err != nil {
		logger.Fatalf("Failed to create output file: %v", err)
	}
	defer outputFile.Close()

	_, err = outputFile.Write(resp.Body)
	if err != nil {
		log.Fatalf("Failed to write data to file: %v", err)
	}

	log.Printf("Documents successfully exported to %s", config.BackupJSONLFilePath())
}

func validateBackupConfig() error {
	parsedURL, err := url.Parse(config.BackupTypesenseHost())
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("invalid typesense host URL: %s", config.BackupTypesenseHost())
	}

	switch {
	case config.BackupTypesenseAPIKey() == "":
		return fmt.Errorf("restore.typesense.api_key cannot be empty")
	case config.BackupCollection() == "":
		return fmt.Errorf("restore.collection cannot be empty")
	case config.BackupJSONLFilePath() == "":
		return fmt.Errorf("restore.jsonl_file_path cannot be empty")
	}

	return nil
}
