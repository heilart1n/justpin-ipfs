package justpin_ipfs

import (
	"fmt"
	"github.com/heilart1n/justpin-ipfs/config"
	"github.com/heilart1n/justpin-ipfs/pinners/infura"
	"github.com/heilart1n/justpin-ipfs/pinners/nftstorage"
	"github.com/heilart1n/justpin-ipfs/pinners/pinata"
	"github.com/heilart1n/justpin-ipfs/pinners/web3storage"
	"io"
	"net/http"
)

type ClientName string

const (
	ClientNameInfura      ClientName = "Infura"
	ClientNameNFTStorage  ClientName = "NFTStorage"
	ClientNamePinata      ClientName = "Pinata"
	ClientNameWeb3Storage ClientName = "Web3Storage"
)

type Pinner interface {
	Name() string
	PinFile(fp string) (string, error)
	PinWithReader(rd io.Reader) (string, error)
	PinWithBytes(buf []byte) (string, error)
	PinHash(hash string) (bool, error)
	PinDir(name string) (string, error)
	Pin(path interface{}) (cid string, err error)
}

type Pinners struct {
	Infura      *infura.Client
	NFTStorage  *nftstorage.Client
	Pinata      *pinata.Client
	Web3Storage *web3storage.Client
}

func NewPinners(infuraCfg, nftStorageCfg, pinataCfg, web3StorageCfg config.Config) *Pinners {
	return &Pinners{
		Infura:      infura.NewClient(infuraCfg, http.DefaultClient),
		NFTStorage:  nftstorage.NewClient(nftStorageCfg, http.DefaultClient),
		Pinata:      pinata.NewClient(pinataCfg, http.DefaultClient),
		Web3Storage: web3storage.NewClient(web3StorageCfg, http.DefaultClient),
	}
}

func (pinners *Pinners) GetPinner(client ClientName) (Pinner, error) {
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

func (pinners *Pinners) MustGetPinner(client ClientName) Pinner {
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

type Handler struct {
	clientName ClientName
	cfg        config.Config
	httpClient *http.Client
}

func NewHandler(cfg config.Config, clientName ClientName) *Handler {
	return &Handler{cfg: cfg, clientName: clientName}
}

func (handler *Handler) CreateAndGetPinner() (Pinner, error) {
	switch handler.clientName {
	case ClientNameInfura:
		return infura.NewClient(handler.cfg, handler.httpClient), nil
	case ClientNameNFTStorage:
		return nftstorage.NewClient(handler.cfg, handler.httpClient), nil
	case ClientNamePinata:
		return pinata.NewClient(handler.cfg, handler.httpClient), nil
	case ClientNameWeb3Storage:
		return web3storage.NewClient(handler.cfg, handler.httpClient), nil
	default:
		return nil, fmt.Errorf("client %s not implemented", handler.clientName)
	}
}

func (handler *Handler) MustCreateAndGetPinner() Pinner {
	switch handler.clientName {
	case ClientNameInfura:
		return infura.NewClient(handler.cfg, handler.httpClient)
	case ClientNameNFTStorage:
		return nftstorage.NewClient(handler.cfg, handler.httpClient)
	case ClientNamePinata:
		return pinata.NewClient(handler.cfg, handler.httpClient)
	case ClientNameWeb3Storage:
		return web3storage.NewClient(handler.cfg, handler.httpClient)
	default:
		return nftstorage.NewClient(handler.cfg, handler.httpClient)
	}
}

// WithHttpClient attach http.Client
func (handler *Handler) WithHttpClient(c *http.Client) *Handler {
	handler.httpClient = c
	return handler
}