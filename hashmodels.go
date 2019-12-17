package hashmodels

import (
	"database/sql"
	"fmt"

	// Getting the entire sqlite3 driver set.
	_ "github.com/mattn/go-sqlite3"
)

// HashRepository is a handler for wrapping the db pointer.
type HashRepository struct {
	db            *sql.DB
	hashTableName string
}

// HashStat represents the schema for the values stored in the Hashsets.
type HashStat struct {
	hashValue              string
	countID                int
	hashTimeInMilliseconds int
}

// TotalStats presents the schema for the total stats queried from the DB for all queries.
type TotalStats struct {
	count       int
	averageTime float32
}

// NewHashStore creates an initialized HashStore
func NewHashStore(db *sql.DB, hashTableName string) (*HashRepository, error) {
	hashStore := HashRepository{
		db:            db,
		hashTableName: hashTableName,
	}
	err := hashStore.InitTables()
	if err != nil {
		return &HashRepository{}, err
	}
	return &hashStore, nil
}

// InitTables creates the schemas for the stored hashes returning any errors or nil on success.
func (dataconnection *HashRepository) InitTables() error {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY NOT NULL, countID INTEGER, hashValue TEXT, hashTimeInMilliseconds INTEGER)", dataconnection.hashTableName)
	_, err := dataconnection.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// StoreHash stores a hashed string and the milliseconds to hash it returning an error on error or nil error on success.
func (dataconnection *HashRepository) StoreHash(hash string, hashTimeInMilliseconds int) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (hashValue, hashTimeInMilliseconds, countID) VALUES (?, ?, (SELECT COUNT(countID) + 1 FROM %s))", dataconnection.hashTableName, dataconnection.hashTableName)
	insert, err := dataconnection.db.Prepare(query)
	if err != nil {
		return -1, err
	}
	defer insert.Close()
	rows, execErr := insert.Exec(hash, hashTimeInMilliseconds)
	if execErr != nil {
		return -1, execErr
	}

	countID, _ := rows.LastInsertId()
	return int(countID), nil
}

// GetHashStat returns the HashStat for a stored countID or an error on error.
func (dataconnection *HashRepository) GetHashStat(countID int) (HashStat, error) {
	query := fmt.Sprintf("SELECT hashValue, hashTimeInMilliseconds from %s where countID=?", dataconnection.hashTableName)
	rows, queryError := dataconnection.db.Query(query, countID)
	if queryError != nil {
		return HashStat{}, queryError
	}
	defer rows.Close()

	var hashValue string
	var hashTimeInMilliseconds int
	rows.Next()
	rows.Scan(&hashValue, &hashTimeInMilliseconds)
	return HashStat{hashValue, countID, hashTimeInMilliseconds}, nil
}

// GetTotalStats gets the total stats for all hashes.
func (dataconnection *HashRepository) GetTotalStats() (TotalStats, error) {
	query := fmt.Sprintf("SELECT COUNT(countID), AVG(hashTimeInMilliseconds from %s", dataconnection.hashTableName)
	rows, queryError := dataconnection.db.Query(query)
	if queryError != nil {
		return TotalStats{}, queryError
	}
	rows.Close()

	var count int
	var avgTime float32
	rows.Next()
	rows.Scan(&count, &avgTime)
	return TotalStats{count, avgTime}, nil
}

// ClearStore clears all data for the hash tables.
func (dataconnection *HashRepository) ClearStore() error {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", dataconnection.hashTableName)
	_, err := dataconnection.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// Close should be called to clear all data.
func (dataconnection *HashRepository) Close() error {
	err := dataconnection.db.Close()
	if err != nil {
		return err
	}
	return nil
}
