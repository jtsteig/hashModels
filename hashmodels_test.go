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

	countID, storeError := hashStore.storeHash("testHash", 500)
	if storeError != nil {
		t.Errorf("Failed to store hash: %q", storeError)
	}

	result, hashErr := hashStore.getHashStat(countID)
	if hashErr != nil {
		t.Errorf("Error getting hashStats: %q", hashErr)
	}

	expected := HashStat{"testHash", countID, 500}
	if expected.hashValue != result.hashValue {
		t.Errorf("Got incorrect hash value: %q and expected %q", expected.hashValue, result.hashValue)
	}
	if expected.countID != result.countID {
		t.Errorf("Got incorrect countId: %q and expected %q", expected.countID, result.countID)
	}
	if expected.hashTimeInMilliseconds != result.hashTimeInMilliseconds {
		t.Errorf("Got incorrect hashtime value: %q and expected %q", expected.hashTimeInMilliseconds, result.hashTimeInMilliseconds)
	}

	dropErr := hashStore.clearStore()
	if dropErr != nil {
		t.Errorf("Failed to drop table %q", dropErr)
	}
	closeErr := hashStore.close()
	if closeErr != nil {
		t.Errorf("Failed to close db: %q", closeErr)
	}
}
