package config

import "time"

const (
	DefaultHTTPTimeout             = 10 * time.Second
	DefaultHTTPTLSHandshakeTimeout = 5 * time.Second
	DefaultTLSInsecureSkipVerify   = true

	DefaultTypesenseCacheTTL = 5

	DefaultMigrationSizePerPage = 100
	DefaultBackupSizePerPage    = 100
	DefaultRestoreBatchSize     = 100
)
