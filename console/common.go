package console

import (
	"fmt"
	"log"
	"net/http"
	"typesense-migration-tools/config"

	"github.com/kumparan/go-connect"
	"github.com/kumparan/go-utils"
	"github.com/typesense/typesense-go/v2/typesense"
	typesenseAPI "github.com/typesense/typesense-go/v2/typesense/api"
)

func newHTTPClient() *http.Client {
	return connect.NewHTTPConnection(&connect.HTTPConnectionOptions{
		TLSHandshakeTimeout:   config.HTTPTLSHandshakeTimeout(),
		Timeout:               config.HTTPTimeout(),
		TLSInsecureSkipVerify: config.DefaultTLSInsecureSkipVerify, //nolint:gosec
	})
}

func newTypesenseClient(host, apiKey string) typesense.APIClientInterface {
	cli, err := typesenseAPI.NewClientWithResponses(
		host,
		typesenseAPI.WithAPIKey(apiKey),
		typesenseAPI.WithHTTPClient(newHTTPClient()),
	)
	if err != nil {
		log.Fatal(err)
	}

	return cli
}

func isTypesenseErrorResponse(response *typesenseAPI.SearchCollectionResponse) bool {
	return response.StatusCode() != http.StatusOK || response.JSON200 == nil
}

func dumpTypesenseSearchResponseError(response *typesenseAPI.SearchCollectionResponse) error {
	return fmt.Errorf("unexpected response from typesense, code: %d, response: %s", response.StatusCode(), string(response.Body))
}

func dumpTypesenseError(messages ...any) error {
	errorMsg := ""
	for _, v := range messages {
		errorMsg = fmt.Sprintf("%s, %s", errorMsg, utils.Dump(v))
	}
	return fmt.Errorf("unexpected response from typesense, %s", errorMsg)
}
