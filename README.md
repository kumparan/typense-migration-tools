# typesense-migration-tools
==================

Tools for backup, restore, migrating documents, and deleting collection gracefully on Typesense. The configurations is provided via a `config.yml` file. The application includes validation and credential confirmation before proceeding with the export.

## Requirements
- Go 1.22
- Typesense server
- YAML configuration file (`config.yml`)

## Installation
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Backup
Backup console application allows you to export documents from a Typesense collection to a JSONL file.
- Ensure that your Typesense server is running and accessible.
- The specified collection must exist and contain data for export.

### Usage
1. Run the application:
   ```bash
   go run main.go backup
   ```

2. Confirm credentials:
   The application will display the config for confirmation:
   ```
   Typesense Host: http://localhost:8108
   Typesense API Key: YOUR_API_KEY
   Collection Name: collection_name
   JSONL File Path: path/to/example.jsonl
   Filter: created_at:<1488325530496000000
   Include Fields: field1,field2,field3
   Exclude Fields: out_of
   Do you want to proceed with these credentials? (yes/no):
   ```
   Type `yes` to proceed or `no` to cancel the operation.

3. The application will export the documents from the specified collection and save them to a JSONL file:
   ```
   Documents successfully exported to path/to/example.jsonl
   ```

### Example Output
```jsonl
{"id": "1", "name": "Document 1", "description": "This is the first document."}
{"id": "2", "name": "Document 2", "description": "This is the second document."}
```

## Restore
Restore console application allows you to import documents to a Typesense collection from a JSONL file.
- Ensure that your Typesense server is running and accessible.
- The specified collection must exist and contain data for restore.

### Usage
1. Run the application:
   ```bash
   go run main.go restore
   ```

2. Confirm credentials:
   The application will display the config for confirmation:
   ```
   Typesense Host: http://localhost:8108
   Typesense API Key: YOUR_API_KEY
   Collection Name: collection_name
   JSONL File Path: path/to/example.jsonl
   Batch Size: 100
   Do you want to proceed with these credentials? (yes/no):
   ```
   Type `yes` to proceed or `no` to cancel the operation.

3. The application will import the documents to the specified collection from a JSONL file:
   ```
   Documents successfully imported to collection_name
   ```

## Migrate
Migrate console application allows you to import documents to a Typesense collection from another Typesense collection.
- Ensure that your Typesense server is running and accessible.
- The source collection must exist and contain data for migration.
- The destination collection must exist.

### Usage
1. Run the application:
   ```bash
   go run main.go migrate
   ```

2. Confirm credentials:
   The application will display the config for confirmation:
   ```
   Source Typesense Host: http://localhost:8108
   Source Typesense API Key: YOUR_API_KEY
   Source Collection Name: source_collection_name
   Destination Typesense Host: http://localhost:8108
   Destination Typesense API Key: YOUR_API_KEY
   Destination Collection Name: destination_collection_name
   Filter: created_at:<1488325530496000000
   Sorter: created_at:desc
   Include Fields: field1,field2,field3
   Exclude Fields: out_of
   Batch Size: 100
   Do you want to proceed with these credentials? (yes/no):
   ```
   Type `yes` to proceed or `no` to cancel the operation.

3. The application will import the documents to the specified collection from a JSONL file:
   ```
   Documents successfully migrated from source_collection_name to destination_collection_name
   ```

## Delete Collection

Delete Collection console application helps manage large-scale document deletions in a Typesense collection. It enables batched deletions to optimize resource usage and avoid overloading the system.
- Ensure that your Typesense server is running and accessible.
- The specified collection must exist.

### Usage
1. Run the application:
   ```bash
   go run main.go delete-collection
   ```

2. Confirm credentials:
   The application will display the config for confirmation:
   ```
   Typesense Host: http://localhost:8108
   Typesense API Key: YOUR_API_KEY
   Collection Name: collection_name
   Batch Size: 100
   Exclude Fields: out_of
   Do you want to proceed with these credentials? (yes/no):
   ```
   Type `yes` to proceed or `no` to cancel the operation.

3. The application will delete documents in batches:
   ```
   Collection collection_name successfully deleted
   ```

## License
This project is licensed under the MIT License. See the `LICENSE` file for details.

## Contribution
Feel free to open issues or submit pull requests for improvements and bug fixes.
