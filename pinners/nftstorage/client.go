package nftstorage

import (
	"github.com/heilart1n/justpin-ipfs/config"
	"net/http"
)

const (
	APIUrl     = "https://api.nft.storage"
	ClientName = "NFTStorage"
)

// Client NFTStorage represents an NFTStorage configuration.
type Client struct {
	*http.Client
	cfg        config.Config
	clientName string
}

func NewClient(cfg config.Config, httpClient *http.Client) *Client {
	return &Client{cfg: cfg, clientName: ClientName, Client: httpClient}
}
