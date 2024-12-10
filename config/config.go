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

// MigrationSourceTypesenseHost :nodoc:
func MigrationSourceTypesenseHost() string {
	return viper.GetString("migration.source.typesense.host")
}

// MigrationSourceTypesenseAPIKey :nodoc:
func MigrationSourceTypesenseAPIKey() string {
	return viper.GetString("migration.source.typesense.api_key")
}

// MigrationDestinationTypesenseHost :nodoc:
func MigrationDestinationTypesenseHost() string {
	return viper.GetString("migration.destination.typesense.host")
}

// MigrationDestinationTypesenseAPIKey :nodoc:
func MigrationDestinationTypesenseAPIKey() string {
	return viper.GetString("migration.destination.typesense.api_key")
}

// MigrationBatchSize :nodoc:
func MigrationBatchSize() int {
	return utils.ValueOrDefault[int](viper.GetInt("migration.batch_size"), DefaultMigrationBatchSize)
}

// MigrationSleepInterval :nodoc:
func MigrationSleepInterval() time.Duration {
	return utils.ValueOrDefault[time.Duration](viper.GetDuration("migration.sleep_interval"), DefaultMigrationSleepInterval)
}

// MigrationSourceCollection :nodoc:
func MigrationSourceCollection() string {
	return viper.GetString("migration.source.collection")
}

// MigrationDestinationCollection :nodoc:
func MigrationDestinationCollection() string {
	return viper.GetString("migration.destination.collection")
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

// BackupTypesenseHost :nodoc:
func BackupTypesenseHost() string {
	return viper.GetString("backup.typesense.host")
}

// BackupTypesenseAPIKey :nodoc:
func BackupTypesenseAPIKey() string {
	return viper.GetString("backup.typesense.api_key")
}

// BackupCollection :nodoc:
func BackupCollection() string {
	return viper.GetString("backup.collection")
}

// BackupSleepInterval :nodoc:
func BackupSleepInterval() time.Duration {
	return utils.ValueOrDefault[time.Duration](viper.GetDuration("backup.sleep_interval"), DefaultBackupSleepInterval)
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

// BackupFolderPath :nodoc:
func BackupFolderPath() string {
	return viper.GetString("backup.folder_path")
}

// BackupMaxDocsPerFile :nodoc:
func BackupMaxDocsPerFile() int {
	return viper.GetInt("backup.max_docs_per_file")
}

// RestoreSleepInterval :nodoc:
func RestoreSleepInterval() time.Duration {
	return utils.ValueOrDefault[time.Duration](viper.GetDuration("restore.sleep_interval"), DefaultRestoreSleepInterval)
}

// BackupBatchSize :nodoc:
func BackupBatchSize() int {
	return utils.ValueOrDefault[int](viper.GetInt("backup.batch_size"), DefaultBackupBatchSize)
}

// BackupSorter :nodoc:
func BackupSorter() string {
	return viper.GetString("backup.sorter")
}

// RestoreTypesenseHost :nodoc:
func RestoreTypesenseHost() string {
	return viper.GetString("restore.typesense.host")
}

// RestoreTypesenseAPIKey :nodoc:
func RestoreTypesenseAPIKey() string {
	return viper.GetString("restore.typesense.api_key")
}

// RestoreCollection :nodoc:
func RestoreCollection() string {
	return viper.GetString("restore.collection")
}

// RestoreFolderPath :nodoc:
func RestoreFolderPath() string {
	return viper.GetString("restore.folder_path")
}

// RestoreBatchSize :nodoc:
func RestoreBatchSize() int {
	return utils.ValueOrDefault[int](viper.GetInt("restore.batch_size"), DefaultRestoreBatchSize)
}

// TypesenseHostForCollectionDeletion :nodoc:
func TypesenseHostForCollectionDeletion() string {
	return viper.GetString("delete_collection.typesense.host")
}

// TypesenseAPIKeyForCollectionDeletion :nodoc:
func TypesenseAPIKeyForCollectionDeletion() string {
	return viper.GetString("delete_collection.typesense.api_key")
}

// CollectionNameToDelete :nodoc:
func CollectionNameToDelete() string {
	return viper.GetString("delete_collection.collection")
}

// BatchSizeForCollectionDeletion :nodoc:
func BatchSizeForCollectionDeletion() int {
	return utils.ValueOrDefault[int](viper.GetInt("delete_collection.batch_size"), DefaultBatchSizeForCollectionDeletion)
}

// SorterForCollectionDeletion :nodoc:
func SorterForCollectionDeletion() string {
	return viper.GetString("delete_collection.sorter")
}

// SleepIntervalForCollectionDeletion :nodoc:
func SleepIntervalForCollectionDeletion() time.Duration {
	return utils.ValueOrDefault[time.Duration](viper.GetDuration("delete_collection.sleep_interval"), DefaultSleepIntervalForCollectionDeletion)
}

// ExcludeFieldsForCollectionDeletion :nodoc:
func ExcludeFieldsForCollectionDeletion() []string {
	return viper.GetStringSlice("delete_collection.exclude_fields")
}

// GetConf :nodoc:
func GetConf() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./..")
	viper.AddConfigPath("./../..")
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
