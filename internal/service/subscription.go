package service

import (
	"fmt"
	"txparser/internal/storage"
)

type subscription struct {
	storage storage.Storage
	eth     Eth
}

type Subscription interface {
	Subscribe(address string) error
}

func NewUser(s storage.Storage, eth Eth) Subscription {
	return &subscription{
		storage: s,
		eth:     eth,
	}
}

func (p *subscription) Subscribe(address string) error {
	if _, ok := p.storage.GetAddressStartBlock(address); ok {
		return fmt.Errorf("address %s is already subscribed", address)
	}
	p.storage.SetAddressStartBlock(address, p.eth.GetLastBlock())
	return nil
}
