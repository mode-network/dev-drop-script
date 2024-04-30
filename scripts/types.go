package scripts

type Config struct {
	SafeAddress    string
	DevDropFactory string
	PhotonId       string
	OrbId          string
}

type GenericStringList struct {
	list []string
}

type TransactionsBatch struct {
	Version      string              `json:"version"`
	ChainId      string              `json:"chainId"`
	CreatedAt    int64               `json:"createdAt"`
	Meta         MetaType            `json:"meta"`
	Transactions []TransactionRecord `json:"transactions"`
}

type MetaType struct {
	Name                    string `json:"name"`
	Description             string `json:"description"`
	TxBuilderVersion        string `json:"txBuilderVersion"`
	CreatedFromSafeAddress  string `json:"createdFromSafeAddress"`
	CreatedFromOwnerAddress string `json:"createdFromOwnerAddress"`
	Checksum                string `json:"checksum"`
}

type TransactionRecord struct {
	To                   string               `json:"to"`
	Value                string               `json:"value"`
	Data                 *[]byte              `json:"data"`
	ContractMethod       ContractMethodType   `json:"contractMethod"`
	ContractInputsValues ContractInputsRecord `json:"contractInputsValues"`
}

type ContractInputsRecord struct {
	From   string `json:"from"`
	To     string `json:"to"`
	ID     string `json:"id"`
	Amount string `json:"amount"`
	Data   string `json:"data"`
}

type ContractMethodType struct {
	Inputs  []Input `json:"inputs"`
	Name    string  `json:"name"`
	Payable bool    `json:"payable"`
}

type Input struct {
	InternalType string `json:"internalType"`
	Name         string `json:"name"`
	Type         string `json:"type"`
}
