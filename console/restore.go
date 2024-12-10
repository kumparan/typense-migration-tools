package console

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
	"typesense-migration-tools/config"

	"github.com/kumparan/go-utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/typesense/typesense-go/v2/typesense"
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
	fmt.Printf("Folder Path: %s\n", config.RestoreFolderPath())
	fmt.Printf("Batch Size: %d\n", config.RestoreBatchSize())
	fmt.Print("Do you want to proceed with these config? (yes/no): ")

	var confirmation string
	fmt.Scanln(&confirmation)
	if confirmation != "yes" {
		log.Println("Export operation cancelled.")
		return
	}

	files, err := filepath.Glob(filepath.Join(config.RestoreFolderPath(), "*.jsonl"))
	if err != nil {
		log.Error(err)
		return
	}

	if len(files) == 0 {
		log.Error(fmt.Errorf("no backup files found in folder: %s", config.RestoreFolderPath()))
		return
	}

	var (
		ctx      = context.TODO()
		tsClient = newTypesenseClient(config.RestoreTypesenseHost(), config.RestoreTypesenseAPIKey())
	)

	for _, file := range files {
		log.Printf("Restoring from file: %s\n", file)
		if err := restoreFromFile(ctx, tsClient, file); err != nil {
			log.Error(fmt.Errorf("error restoring file %s: %w", file, err))
			return
		}
	}

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
	case config.RestoreFolderPath() == "":
		return fmt.Errorf("restore.folder_path cannot be empty")
	case config.RestoreBatchSize() <= 0:
		return fmt.Errorf("restore.batch_size must be a positive integer")
	}

	return nil
}

func restoreFromFile(ctx context.Context, client typesense.APIClientInterface, filePath string) error {
	logger := log.WithFields(log.Fields{
		"context":         utils.DumpIncomingContext(ctx),
		"collection":      config.RestoreCollection(),
		"batchSize":       config.RestoreBatchSize(),
		"typesenseHost":   config.RestoreTypesenseHost(),
		"typesenseAPIKey": config.RestoreTypesenseAPIKey(),
		"folderPath":      config.RestoreFolderPath(),
		"filePath":        filePath,
	})

	file, err := os.Open(filePath)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer file.Close()

	var (
		scanner     = bufio.NewScanner(file)
		buffer      bytes.Buffer
		batches     [][]byte
		currentLine = 0
		startLine   = 1
	)

	for scanner.Scan() {
		currentLine++
		if currentLine < startLine {
			continue
		}

		buffer.Write(scanner.Bytes())
		buffer.WriteByte('\n')

		if (currentLine-startLine+1)%config.RestoreBatchSize() == 0 {
			batches = append(batches, append([]byte{}, buffer.Bytes()...))
			buffer.Reset()
		}
	}

	if buffer.Len() > 0 {
		batches = append(batches, append([]byte{}, buffer.Bytes()...))
	}

	if err := scanner.Err(); err != nil {
		logger.Error(err)
		return err
	}

	for i, batch := range batches {
		log.Printf("Sending batch %d from file %s\n", i+1, filePath)

		if err := sendBatch(ctx, client, batch); err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

func sendBatch(ctx context.Context, client typesense.APIClientInterface, batchData []byte) error {
	resp, err := client.ImportDocumentsWithBodyWithResponse(ctx, config.RestoreCollection(), &typesenseAPI.ImportDocumentsParams{
		Action:    typesensePtr.String("upsert"),
		BatchSize: typesensePtr.Int(config.RestoreBatchSize()),
	}, "application/jsonl", bytes.NewReader(batchData))
	if err != nil {
		log.Error(err)
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		err = dumpTypesenseError(resp.JSON404, resp.JSON400)
		log.Error(err)
		return err
	}

	time.Sleep(config.RestoreSleepInterval())

	return nil
}
