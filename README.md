# Ethereum Address Tracker

A simple Go-based service to track Ethereum transactions involving specific addresses. You can subscribe to addresses via HTTP, and the service will watch the blockchain for any incoming or outgoing transactions involving those addresses.

## Features

- Subscribe to Ethereum addresses
- Track transactions in new blocks
- Fetch transactions for a given address

## API Endpoints

### Subscribe to an address

**POST** `/subscribe`

Body:
```json
{
  "address": "0xYourEthereumAddress"
}
