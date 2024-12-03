package console

import (
	"context"
	"net/http"
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
	var (
		ctx      = context.TODO()
		tsClient = newTypesenseClient()
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
