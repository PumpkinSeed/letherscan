package communicator

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	// Address of the client
	Address string `json:"address"`

	eth *ethclient.Client
}

func NewClient(address string) *Client {
	client, err := ethclient.Dial(address)
	if err != nil {
		panic(err) // TODO: handle error properly
	}
	return &Client{
		Address: address,
		eth:     client,
	}
}
