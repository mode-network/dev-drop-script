package scripts

import (
	"testing"
)

func TestLoadTransactionBatch(t *testing.T) {
	config := GetConfig()
	path := "sample-photons.csv"
	jsonPath := "tx-batch.json"
	GenerateDevDropSafeFile(path, jsonPath, config)
}
