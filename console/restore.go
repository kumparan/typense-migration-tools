package console

import (
	"bytes"
	"context"
	"net/http"
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
	var (
		ctx      = context.TODO()
		tsClient = newTypesenseClient()
	)

	logger := log.WithFields(log.Fields{
		"context":    utils.DumpIncomingContext(ctx),
		"collection": config.RestoreCollection(),
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
}
