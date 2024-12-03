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
	typesensePtr "github.com/typesense/typesense-go/v2/typesense/api/pointer"
)

func newHTTPClient() *http.Client {
	return connect.NewHTTPConnection(&connect.HTTPConnectionOptions{
		TLSHandshakeTimeout:   config.HTTPTLSHandshakeTimeout(),
		Timeout:               config.HTTPTimeout(),
		TLSInsecureSkipVerify: config.DefaultTLSInsecureSkipVerify, //nolint:gosec
	})
}

func newTypesenseClient() typesense.APIClientInterface {
	cli, err := typesenseAPI.NewClientWithResponses(
		config.TypesenseHost(),
		typesenseAPI.WithAPIKey(config.TypesenseAPIKey()),
		typesenseAPI.WithHTTPClient(newHTTPClient()),
	)
	if err != nil {
		log.Fatal(err)
	}

	return cli
}

func setTypesenseCacheParams(searchParams *typesenseAPI.SearchCollectionParams) {
	if config.TypesenseEnableCache() {
		searchParams.UseCache = typesensePtr.True()
		searchParams.CacheTtl = typesensePtr.Int(config.TypesenseCacheTTL())
	}
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
