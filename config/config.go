package config

import (
	"strings"
	"time"

	"github.com/kumparan/go-utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// LogLevel :nodoc:
func LogLevel() string {
	return viper.GetString("log_level")
}

// HTTPTLSHandshakeTimeout :nodoc:
func HTTPTLSHandshakeTimeout() time.Duration {
	return utils.ValueOrDefault[time.Duration](viper.GetDuration("http_connection_settings.tls_handshake_timeout"), DefaultHTTPTLSHandshakeTimeout)
}

// HTTPTimeout :nodoc:
func HTTPTimeout() time.Duration {
	return utils.ValueOrDefault[time.Duration](viper.GetDuration("http_connection_settings.timeout"), DefaultHTTPTimeout)
}

// TypesenseHost :nodoc:
func TypesenseHost() string {
	return viper.GetString("typesense.host")
}

// TypesenseAPIKey :nodoc:
func TypesenseAPIKey() string {
	return viper.GetString("typesense.api_key")
}

// TypesenseEnableRead :nodoc:
func TypesenseEnableRead() bool {
	return viper.GetBool("typesense.enable_read")
}

// TypesenseEnableCache :nodoc:
func TypesenseEnableCache() bool {
	return viper.GetBool("typesense.enable_cache")
}

// TypesenseCacheTTL :nodoc:
func TypesenseCacheTTL() int {
	return utils.ValueOrDefault(viper.GetInt("typesense.cache_ttl"), DefaultTypesenseCacheTTL)
}

// MigrationSizePerPage :nodoc:
func MigrationSizePerPage() int {
	return utils.ValueOrDefault[int](viper.GetInt("migration.size_per_page"), DefaultMigrationSizePerPage)
}

// MigrationSourceCollection :nodoc:
func MigrationSourceCollection() string {
	return viper.GetString("migration.source_collection")
}

// MigrationDestinationCollection :nodoc:
func MigrationDestinationCollection() string {
	return viper.GetString("migration.destination_collection")
}

// MigrationIncludeFields :nodoc:
func MigrationIncludeFields() []string {
	return viper.GetStringSlice("migration.include_fields")
}

// MigrationExcludeFields :nodoc:
func MigrationExcludeFields() []string {
	return viper.GetStringSlice("migration.exclude_fields")
}

// MigrationSorter :nodoc:
func MigrationSorter() string {
	return viper.GetString("migration.sorter")
}

// MigrationFilter :nodoc:
func MigrationFilter() string {
	return viper.GetString("migration.filter")
}

// BackupCollection :nodoc:
func BackupCollection() string {
	return viper.GetString("backup.collection")
}

// BackupIncludeFields :nodoc:
func BackupIncludeFields() []string {
	return viper.GetStringSlice("backup.include_fields")
}

// BackupExcludeFields :nodoc:
func BackupExcludeFields() []string {
	return viper.GetStringSlice("backup.exclude_fields")
}

// BackupFilter :nodoc:
func BackupFilter() string {
	return viper.GetString("backup.filter")
}

// BackupJSONLFilePath :nodoc:
func BackupJSONLFilePath() string {
	return viper.GetString("backup.jsonl_file_path")
}

// RestoreCollection :nodoc:
func RestoreCollection() string {
	return viper.GetString("restore.collection")
}

// RestoreJSONLFilePath :nodoc:
func RestoreJSONLFilePath() string {
	return viper.GetString("restore.jsonl_file_path")
}

// RestoreBatchSize :nodoc:
func RestoreBatchSize() int {
	return utils.ValueOrDefault[int](viper.GetInt("restore.bathc_size"), DefaultRestoreBatchSize)
}

// GetConf :nodoc:
func GetConf() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./..")
	viper.AddConfigPath("./../..")
	viper.AddConfigPath("./../../..")
	viper.SetConfigName("config")
	viper.SetEnvPrefix("svc")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Warnf("%v", err)
	}
}
