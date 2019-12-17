package hashmodels

import (
	"database/sql"
	"testing"
)

func TestSqlIteHashStoreHappyPath(t *testing.T) {
	filename := "c:\\temp\\testdb.db"
	hashTable := "hashes"
	db, _ := sql.Open("sqlite3", filename)
	hashStore, initErr := NewHashStore(db, hashTable)

	if initErr != nil {
		t.Errorf("Failed to init db: %q", initErr)
	}

	countID, storeError := hashStore.StoreHash("testHash", 500)
	if storeError != nil {
		t.Errorf("Failed to store hash: %q", storeError)
	}

	result, hashErr := hashStore.GetHashStat(countID)
	if hashErr != nil {
		t.Errorf("Error getting hashStats: %q", hashErr)
	}

	expected := HashStat{"testHash", countID, 500}
	if expected.HashValue != result.HashValue {
		t.Errorf("Got incorrect hash value: %q and expected %q", expected.HashValue, result.HashValue)
	}
	if expected.CountID != result.CountID {
		t.Errorf("Got incorrect countId: %q and expected %q", expected.CountID, result.CountID)
	}
	if expected.HashTimeInMilliseconds != result.HashTimeInMilliseconds {
		t.Errorf("Got incorrect hashtime value: %q and expected %q", expected.HashTimeInMilliseconds, result.HashTimeInMilliseconds)
	}

	dropErr := hashStore.ClearStore()
	if dropErr != nil {
		t.Errorf("Failed to drop table %q", dropErr)
	}
	closeErr := hashStore.Close()
	if closeErr != nil {
		t.Errorf("Failed to close db: %q", closeErr)
	}
}
