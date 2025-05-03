package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

const apiUrl = "https://mainnet.infura.io/v3/"

type ETHClient struct {
	apiKey     string
	httpClient *http.Client
}

func MustETHClient(apiKey string) *ETHClient {
	if apiKey == "" {
		log.Fatal("ETH API key is missing. Please set your API key in the code or as an environment variable.")
	}
	return &ETHClient{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (c *ETHClient) GetLatestBlockNumber(ctx context.Context) (int64, error) {
	payload := []byte(`{
		"jsonrpc":"2.0",
		"method":"eth_blockNumber",
		"params":[],
		"id":1
	}`)
	res, err := c.callAPI(ctx, payload)
	if err != nil {
		return 0, err
	}

	var result struct {
		Result string `json:"result"`
	}
	if err := json.Unmarshal(res, &result); err != nil {
		return 0, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return strconv.ParseInt(result.Result, 0, 64)
}

func (c *ETHClient) GetBlock(ctx context.Context, blockNumber int64) ([]byte, error) {
	hexBlock := fmt.Sprintf("0x%x", blockNumber)
	payload := []byte(fmt.Sprintf(`{
		"jsonrpc":"2.0",
		"method":"eth_getBlockByNumber",
		"params":["%s", true],
		"id":1
	}`, hexBlock))

	return c.callAPI(ctx, payload)
}

func (c *ETHClient) callAPI(ctx context.Context, payload []byte) ([]byte, error) {
	url := apiUrl + c.apiKey
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("failed to close response body")
		}
	}()

	return io.ReadAll(resp.Body)
}
