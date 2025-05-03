package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"
	"txparser/internal/client"
	"txparser/internal/service/domain"
	"txparser/internal/storage"
)

type Eth interface {
	Process(ctx context.Context) error
	GetLastBlock() int64
	GetLatestBlockNumber(ctx context.Context) (int64, error)
	GetTransactions(address string) ([]domain.Transaction, error)
}

// Struct Note:
// `currentBlock` and `lastBlock` are atomic to support safe concurrent updates,
// and prepare the system for future parallel processing with worker pools.
//
// If we scale to a high-performance setup with multiple goroutines writing frequently,
// we might run into CPU cache issues like false sharing. In that case, we can consider
// padding or aligning these fields (e.g., using cache line padding) to reduce contention.

type eth struct {
	currentBlock atomic.Int64
	lastBlock    atomic.Int64
	storage      storage.Storage
	ethClient    *client.ETHClient
}

func NewEth(storage storage.Storage, ethClient *client.ETHClient) Eth {
	return &eth{
		storage:   storage,
		ethClient: ethClient,
	}
}

func (e *eth) GetLastBlock() int64 {
	return e.lastBlock.Load()
}

func (e *eth) GetLatestBlockNumber(ctx context.Context) (int64, error) {
	log.Println("Fetching latest block number...")
	return e.ethClient.GetLatestBlockNumber(ctx)
}

func (e *eth) updateLatestBlockNumber(val int64) {
	e.lastBlock.Store(val)
	log.Printf("Updated latest block number to %d", val)
}

func (e *eth) GetBlock(ctx context.Context, blockNumber int64) (*domain.Block, error) {
	log.Printf("Fetching block %d", blockNumber)
	res, err := e.ethClient.GetBlock(ctx, blockNumber)
	if err != nil {
		log.Printf("Error fetching block %d: %v", blockNumber, err)
		return nil, err
	}

	var result struct {
		Result domain.Block `json:"result"`
	}
	if err := json.Unmarshal(res, &result); err != nil {
		log.Printf("Failed to decode block %d JSON: %v", blockNumber, err)
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &result.Result, nil
}

// Design Note:
// I used atomic operations in the parsing logic to safely update data from multiple goroutines.
// This makes it easier to scale later—for example, we could process blocks in parallel
// using a worker pool without running into race conditions.
//
// The code is modular, so in the future we could run it as a cron job instead of a long-running service.
// Since Ethereum creates a new block about every 12 seconds (~5 per minute, ~300 per hour, ~7,200 per day),
// this gives us a good idea of how often we should poll or schedule the job.
//
// To avoid sending too many requests to the Ethereum node, I cache the latest global block number.
// This way, we don’t need to request the current block separately for every subscriber.
//
// Trade-offs:
// - We’re choosing better performance and availability over perfect consistency.
//   Because we use a shared global block number, it’s possible we include some transactions
//   that happened before the user subscribed.
// - If we want stricter consistency, we can store the block number when someone subscribes
//   and start looking for transactions from that point.
// - If we run the cron job every 12 seconds (in sync with new block creation),
//   the risk of showing outdated or irrelevant data becomes very small,
//   especially because we’re using the global block as a cache.

func (e *eth) Process(ctx context.Context) error {
	log.Println("Starting Ethereum block processing...")
	blockNumber, err := e.GetLatestBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest block: %w", err)
	}
	e.updateLatestBlockNumber(blockNumber)
	e.currentBlock.CompareAndSwap(0, blockNumber)
	for i := e.currentBlock.Load(); i < blockNumber; i++ {

		block, err := e.GetBlock(ctx, i)
		if err != nil {
			return fmt.Errorf("failed to fetch block %d: %w", i, err)
		}
		log.Printf("Processing block %d with %d transactions", i, len(block.Txs))
		e.currentBlock.Store(i)
		e.parseTrx(block.Txs, i)
	}
	log.Println("Finished processing blocks")
	return nil
}

func (e *eth) parseTrx(transactions []domain.Transaction, currentBlock int64) {
	subs := e.storage.GetSubscribers()
	addressTrxs := make(map[string][]domain.Transaction, len(subs))
	for _, trx := range transactions {
		from, to := trx.From, trx.To
		blockNumber, exists := subs[from]
		if exists && blockNumber <= currentBlock {
			addressTrxs[from] = append(addressTrxs[from], trx)
		}
		blockNumber, exists = subs[to]
		if exists && blockNumber <= currentBlock {
			addressTrxs[to] = append(addressTrxs[to], trx)
		}
	}
	for user, trxs := range addressTrxs {
		log.Printf("Saving %d transactions for user %s at block %d", len(trxs), user, currentBlock)
		e.storage.SaveTransactions(user, trxs)
	}
}

func (e *eth) GetTransactions(address string) ([]domain.Transaction, error) {
	log.Printf("Fetching transactions for address %s", address)
	return e.storage.GetTransactions(address)
}
