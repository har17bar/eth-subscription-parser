package storage

import (
	"fmt"
	"sync"
	"txparser/internal/service/domain"
)

type Storage interface {
	SetAddressStartBlock(address string, block int64)
	GetAddressStartBlock(address string) (int64, bool)
	SaveTransactions(address string, transactions []domain.Transaction)
	GetTransactions(address string) ([]domain.Transaction, error)
	GetSubscribers() map[string]int64
}

type Transaction struct {
	Hash  string `json:"hash"`
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
	Block int64  `json:"block"`
}

type storage struct {
	mu           sync.RWMutex
	subscribed   map[string]int64
	transactions map[string][]domain.Transaction
}

func NewStorage() Storage {
	return &storage{
		subscribed:   make(map[string]int64),
		transactions: make(map[string][]domain.Transaction),
	}
}

func (s *storage) SetAddressStartBlock(address string, block int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subscribed[address] = block
}

func (s *storage) GetAddressStartBlock(address string) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	block, ok := s.subscribed[address]
	return block, ok
}

func (s *storage) SaveTransactions(address string, transactions []domain.Transaction) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, tx := range transactions {
		s.transactions[address] = append(s.transactions[address], tx)
	}
}

func (s *storage) GetTransactions(address string) ([]domain.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	transaction, ok := s.transactions[address]
	if !ok {
		return nil, fmt.Errorf("transaction not found by address: %s", address)
	}
	return transaction, nil
}

func (s *storage) GetSubscribers() map[string]int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.subscribed
}
