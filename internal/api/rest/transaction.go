package rest

import (
	"encoding/json"
	"net/http"
	"txparser/internal/service"
)

type Transactions struct {
	ethService service.Eth
}

func NewTransactions(ethService service.Eth) *Transactions {
	return &Transactions{
		ethService: ethService,
	}
}

func (e *Transactions) GetCurrentProcessedBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	number, err := e.ethService.GetLatestBlockNumber(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch block number", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]int64{"current_block": number}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (e *Transactions) GetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing 'address' parameter", http.StatusBadRequest)
		return
	}
	txs, err := e.ethService.GetTransactions(address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(txs); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
