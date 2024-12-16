package config

import (
	"strings"
	"time"

	"github.com/kumparan/go-utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// LogLevel determines the verbosity of the application's logging output (Trace, Debug, Info, Warning, Error, Fatal and Panic)
func LogLevel() string {
	return viper.GetString("log_level")
}

// HTTPTLSHandshakeTimeout set a limit on how long the application waits for a TLS handshake to complete when establishing a connection
func HTTPTLSHandshakeTimeout() time.Duration {
	return utils.ValueOrDefault[time.Duration](viper.GetDuration("http_connection_settings.tls_handshake_timeout"), DefaultHTTPTLSHandshakeTimeout)
}

// HTTPTimeout sets the maximum duration the application waits for an HTTP request to complete before timing out
func HTTPTimeout() time.Duration {
	return utils.ValueOrDefault[time.Duration](viper.GetDuration("http_connection_settings.timeout"), DefaultHTTPTimeout)
}

// MigrationSourceTypesenseHost used to specify the hostname or IP address of the Typesense instance from which migration data is sourced
func MigrationSourceTypesenseHost() string {
	return viper.GetString("migration.source.typesense.host")
}

// MigrationSourceTypesenseAPIKey used to authenticate requests to the Typesense instance from which migration data is sourced
func MigrationSourceTypesenseAPIKey() string {
	return viper.GetString("migration.source.typesense.api_key")
}

// MigrationDestinationTypesenseHost used to specify the hostname or IP address of the Typesense instance to which migration data will be written
func MigrationDestinationTypesenseHost() string {
	return viper.GetString("migration.destination.typesense.host")
}

// MigrationDestinationTypesenseAPIKey used to authenticate requests to the Typesense instance to which migration data will be written
func MigrationDestinationTypesenseAPIKey() string {
	return viper.GetString("migration.destination.typesense.api_key")
}

// MigrationBatchSize defines the maximum number of documents to process in a Typesense search batch
func MigrationBatchSize() int {
	return utils.ValueOrDefault[int](viper.GetInt("migration.batch_size"), DefaultMigrationBatchSize)
}

// MigrationSleepInterval defines the duration the application waits between processing consecutive migration batches
func MigrationSleepInterval() time.Duration {
	return utils.ValueOrDefault[time.Duration](viper.GetDuration("migration.sleep_interval"), DefaultMigrationSleepInterval)
}

// MigrationSourceCollection specifies the collection in the source Typesense instance from which data will be migrated
func MigrationSourceCollection() string {
	return viper.GetString("migration.source.collection")
}

// MigrationDestinationCollection specifies the collection in the destination Typesense instance to which data will be migrated
func MigrationDestinationCollection() string {
	return viper.GetString("migration.destination.collection")
}

// MigrationIncludedFields specifies which fields from the source collection should be migrated to the destination collection (optional)
func MigrationIncludedFields() []string {
	return viper.GetStringSlice("migration.included_fields")
}

// MigrationExcludedFields specifies which fields from the source collection should not be migrated to the destination collection (optional)
func MigrationExcludedFields() []string {
	return viper.GetStringSlice("migration.excluded_fields")
}

// MigrationSorter specifies the field or criteria by which the data should be sorted while being migrated from the source to the destination collection
func MigrationSorter() string {
	return viper.GetString("migration.sorter")
}

// MigrationFilter specifies a condition or query to filter which documents from the source collection should be migrated to the destination collection
func MigrationFilter() string {
	return viper.GetString("migration.filter")
}

// BackupTypesenseHost specifies the hostname or IP address of the Typesense server where backup operations are performed
func BackupTypesenseHost() string {
	return viper.GetString("backup.typesense.host")
}

// BackupTypesenseAPIKey used to authenticate requests to the Typesense instance during backup operations
func BackupTypesenseAPIKey() string {
	return viper.GetString("backup.typesense.api_key")
}

// BackupCollection specifies the collection in the Typesense instance from which data will be backed up
func BackupCollection() string {
	return viper.GetString("backup.collection")
}

// BackupSleepInterval defines the duration the application waits between consecutive backup tasks
func BackupSleepInterval() time.Duration {
	return utils.ValueOrDefault[time.Duration](viper.GetDuration("backup.sleep_interval"), DefaultBackupSleepInterval)
}

// BackupIncludedFields specifies which fields from the collection should be included in the backup (optional)
func BackupIncludedFields() []string {
	return viper.GetStringSlice("backup.included_fields")
}

// BackupExcludedFields specifies which fields from the collection should not be included in the backup (optional)
func BackupExcludedFields() []string {
	return viper.GetStringSlice("backup.excluded_fields")
}

// BackupFilter specifies a condition or query to filter which records from the collection should be backed up
func BackupFilter() string {
	return viper.GetString("backup.filter")
}

// BackupFolderPath specifies the directory location on the filesystem where the backup files will be saved
func BackupFolderPath() string {
	return viper.GetString("backup.folder_path")
}

// BackupMaxDocsPerFile defines how many documents will be stored in each backup file before creating a new one
func BackupMaxDocsPerFile() int {
	return viper.GetInt("backup.max_docs_per_file")
}

// BackupBatchSize defines the number of documents to be processed in each backup batch
func BackupBatchSize() int {
	return utils.ValueOrDefault[int](viper.GetInt("backup.batch_size"), DefaultBackupBatchSize)
}

// BackupSorter specifies the field or criteria by which the data should be sorted while performing backup operations
func BackupSorter() string {
	return viper.GetString("backup.sorter")
}

// RestoreTypesenseHost specifies the hostname or IP address of the Typesense server where restore operations will be performed
func RestoreTypesenseHost() string {
	return viper.GetString("restore.typesense.host")
}

// RestoreTypesenseAPIKey used to authenticate requests to the Typesense instance during restore operations
func RestoreTypesenseAPIKey() string {
	return viper.GetString("restore.typesense.api_key")
}

// RestoreCollection specifies the collection in the Typesense instance where data will be restored
func RestoreCollection() string {
	return viper.GetString("restore.collection")
}

// RestoreFolderPath specifies the directory location on the filesystem where the restore files can be found
func RestoreFolderPath() string {
	return viper.GetString("restore.folder_path")
}

// RestoreBatchSize defines the number of documents to be processed in each restore batch
func RestoreBatchSize() int {
	return utils.ValueOrDefault[int](viper.GetInt("restore.batch_size"), DefaultRestoreBatchSize)
}

// RestoreSleepInterval defines the duration the application waits between consecutive restore tasks
func RestoreSleepInterval() time.Duration {
	return utils.ValueOrDefault[time.Duration](viper.GetDuration("restore.sleep_interval"), DefaultRestoreSleepInterval)
}

// TypesenseHostForCollectionDeletion specifies the hostname or IP address of the Typesense server where the collection deletion operation will be performed
func TypesenseHostForCollectionDeletion() string {
	return viper.GetString("delete_collection.typesense.host")
}

// TypesenseAPIKeyForCollectionDeletion used to authenticate requests to the Typesense instance during the collection deletion process
func TypesenseAPIKeyForCollectionDeletion() string {
	return viper.GetString("delete_collection.typesense.api_key")
}

// CollectionNameToDelete specifies the name of the collection that will be deleted
func CollectionNameToDelete() string {
	return viper.GetString("delete_collection.collection")
}

// BatchSizeForCollectionDeletion defines the number of documents to be processed in each deletion batch
func BatchSizeForCollectionDeletion() int {
	return utils.ValueOrDefault[int](viper.GetInt("delete_collection.batch_size"), DefaultBatchSizeForCollectionDeletion)
}

// SorterForCollectionDeletion specifies the field or criteria by which the documents should be sorted during the deletion operation
func SorterForCollectionDeletion() string {
	return viper.GetString("delete_collection.sorter")
}

// SleepIntervalForCollectionDeletion defines the duration the application waits between consecutive deletion tasks
func SleepIntervalForCollectionDeletion() time.Duration {
	return utils.ValueOrDefault[time.Duration](viper.GetDuration("delete_collection.sleep_interval"), DefaultSleepIntervalForCollectionDeletion)
}

// ExcludedFieldsForCollectionDeletion specifies which fields from the collection should not be included in the deletion process (optional)
func ExcludedFieldsForCollectionDeletion() []string {
	return viper.GetStringSlice("delete_collection.excluded_fields")
}

// GetConf read the configuration file
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
