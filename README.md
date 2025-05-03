# Ethereum Address Tracker

A simple Go-based service to track Ethereum transactions involving specific addresses. You can subscribe to addresses via HTTP, and the service will watch the blockchain for any incoming or outgoing transactions involving those addresses.

## Features

* Subscribe to Ethereum addresses
* Track transactions in new blocks
* Fetch transactions for a given address
* Get the latest processed block

## API Endpoints

### Subscribe to an address

**POST** `/subscribe`
**Body:**

```json
{
  "address": "0xYourEthereumAddress"
}
```

### Get transactions for an address

**GET** `/transactions?address=0xYourEthereumAddress`
**Response:**

```json
[
  {
    "from": "0x...",
    "to": "0x...",
    "value": "1000000000000000000",
    "block": 19000000
  }
]
```

### Get the latest processed block

**GET** `/current`
**Response:**

```json
{
  "latest_block": 19000021
}
```
