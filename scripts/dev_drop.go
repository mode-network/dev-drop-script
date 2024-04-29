package scripts

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var contractMethod = ContractMethodType{
	Inputs: []Input{
		{InternalType: "address", Name: "from", Type: "address"},
		{InternalType: "address", Name: "to", Type: "address"},
		{InternalType: "uint256", Name: "id", Type: "uint256"},
		{InternalType: "uint256", Name: "amount", Type: "uint256"},
		{InternalType: "bytes", Name: "data", Type: "bytes"},
	},
	Name:    "safeTransferFrom",
	Payable: false,
}

func LoadCSVFile(path string, columnCount int) (*[]GenericStringList, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Error reading csv file %s: %v", path, err)
		panic(err)
	}
	defer f.Close()

	var list []GenericStringList
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for i, line := range data {
		if i == 0 {
			continue
		}
		if len(line) != columnCount {
			return &list, errors.New(fmt.Sprintf(
				"Invalid file content, line %d has wrong number of columns %d, expected %d columns",
				i, len(line), columnCount))
		}
		l := GenericStringList{}
		for _, value := range line {
			l.list = append(l.list, strings.Trim(value, " "))
		}
		list = append(list, l)
	}
	return &list, nil
}

func loadTransactionBatch(filePath string) (TransactionsBatch, error) {
	var tb TransactionsBatch
	if filePath == "" {
		panic("filePath is required but is empty string")
	}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return tb, err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &tb)
	if err != nil {
		fmt.Println("error reading TransactionsBatch from json file: ", err)
		return tb, err
	}
	return tb, nil
}

func saveTransactionsBatch(filePath string, transactionsBatch TransactionsBatch) {
	data, _ := json.MarshalIndent(transactionsBatch, "", "  ")
	ioutil.WriteFile(filePath, data, 0644)
}

func newContractInputsRecord(from string, to string, id string, amount string, data string) ContractInputsRecord {
	return ContractInputsRecord{
		From: from, To: to, ID: id, Amount: amount, Data: data,
	}
}

func newTransactionRecord(to string, value string, data *[]byte, inputs ContractInputsRecord) TransactionRecord {
	return TransactionRecord{
		To: to, Value: value, Data: data, ContractMethod: contractMethod, ContractInputsValues: inputs,
	}
}

func newTransactionBatch(version string, chainId string, safeAddress string) TransactionsBatch {
	if version == "" {
		version = "1.0"
	}
	if chainId == "" {
		chainId = "34443"
	}
	var transactions []TransactionRecord
	txBatch := TransactionsBatch{
		Version:   version,
		ChainId:   chainId,
		CreatedAt: time.Now().Unix(),
		Meta: MetaType{
			Name:                    "Transactions Batch",
			Description:             "",
			TxBuilderVersion:        "1.16.5",
			CreatedFromSafeAddress:  safeAddress,
			CreatedFromOwnerAddress: "",
			Checksum:                "0x91049de6874cc7eb4591c73bcb84442deb9933cd487c0b75006b496fbfa492c7",
		},
		Transactions: transactions,
	}
	return txBatch
}

func validateInputAmount(value string, name string) {
	_, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("bad input value %s for %s (expected an integer value): %v", value, name, err))
	}
}

func GetConfig() Config {
	return Config{
		SafeAddress:    os.Getenv("SafeAddress"),
		DevDropFactory: os.Getenv("DevDropFactory"),
		PhotonId:       os.Getenv("PhotonId"),
		OrbId:          os.Getenv("OrbId"),
	}
}

func GenerateDevDropSafeFile(inFile string, outFile string, config Config) {
	path := inFile
	list, err := LoadCSVFile(path, 3)
	if err != nil {
		panic(fmt.Sprintf("error reading csv file %s: %v", path, err))
	}
	log.Printf("loaded csv data from file %s: %v", path, list)

	txBatch := newTransactionBatch("", "", config.SafeAddress)
	var tx TransactionRecord
	for _, l := range *list {
		wallet := l.list[0]
		photonCount := l.list[1]
		orbCount := l.list[2]
		validateInputAmount(photonCount, "photon count")
		validateInputAmount(orbCount, "orb count")

		if photonCount != "0" {
			inputs := newContractInputsRecord(config.SafeAddress, wallet, config.PhotonId, photonCount, "0x00")
			tx = newTransactionRecord(config.DevDropFactory, "0", nil, inputs)
		}
		if orbCount != "0" {
			inputs := newContractInputsRecord(config.SafeAddress, wallet, config.OrbId, orbCount, "0x00")
			tx = newTransactionRecord(config.DevDropFactory, "0", nil, inputs)
		}
		txBatch.Transactions = append(txBatch.Transactions, tx)
	}

	jsonPath := outFile
	saveTransactionsBatch(jsonPath, txBatch)

}
