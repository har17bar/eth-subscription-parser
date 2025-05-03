package domain

type Transaction struct {
	Hash  string `json:"hash"`
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
	Block int    `json:"block"`
}

type Block struct {
	Txs []Transaction `json:"transactions"`
}
