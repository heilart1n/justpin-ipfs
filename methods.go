package justpin_ipfs

import (
	"fmt"
	"github.com/heilart1n/justpin-ipfs/config"
	"github.com/heilart1n/justpin-ipfs/pinners"
	"github.com/heilart1n/justpin-ipfs/pinners/infura"
	"github.com/heilart1n/justpin-ipfs/pinners/nftstorage"
	"github.com/heilart1n/justpin-ipfs/pinners/pinata"
	"github.com/heilart1n/justpin-ipfs/pinners/web3storage"
	"net/http"
)

func NewPinnerWithHTTPClient(cfg config.Config, clientName ClientName, httpClient *http.Client) (pinners.Pinner, error) {
	switch clientName {
	case ClientNameInfura:
		return infura.NewClient(cfg, httpClient), nil
	case ClientNameNFTStorage:
		return nftstorage.NewClient(cfg, httpClient), nil
	case ClientNamePinata:
		return pinata.NewClient(cfg, httpClient), nil
	case ClientNameWeb3Storage:
		return web3storage.NewClient(cfg, httpClient), nil
	default:
		return nil, fmt.Errorf("client %s not implemented", clientName)
	}
}

func MustNewPinnerWithHTTPClient(cfg config.Config, clientName ClientName, httpClient *http.Client) pinners.Pinner {
	pinner, err := NewPinnerWithHTTPClient(cfg, clientName, httpClient)
	if err != nil {
		return nftstorage.NewClient(cfg, httpClient)
	}
	return pinner
}

func NewPinner(cfg config.Config, clientName ClientName) (pinners.Pinner, error) {
	return NewPinnerWithHTTPClient(cfg, clientName, http.DefaultClient)
}

func MustNewPinner(cfg config.Config, clientName ClientName) pinners.Pinner {
	return MustNewPinnerWithHTTPClient(cfg, clientName, http.DefaultClient)
}

func (pinners *Pinners) GetPinner(client ClientName) (pinners.Pinner, error) {
	switch client {
	case ClientNameInfura:
		return pinners.Infura, nil
	case ClientNameNFTStorage:
		return pinners.NFTStorage, nil
	case ClientNamePinata:
		return pinners.Pinata, nil
	case ClientNameWeb3Storage:
		return pinners.Web3Storage, nil
	default:
		return nil, fmt.Errorf("client %s not implemented", client)
	}
}

func (pinners *Pinners) MustGetPinner(client ClientName) pinners.Pinner {
	switch client {
	case ClientNameInfura:
		return pinners.Infura
	case ClientNameNFTStorage:
		return pinners.NFTStorage
	case ClientNamePinata:
		return pinners.Pinata
	case ClientNameWeb3Storage:
		return pinners.Web3Storage
	default:
		return pinners.NFTStorage
	}
}
