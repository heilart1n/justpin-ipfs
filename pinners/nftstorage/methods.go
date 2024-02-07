package nftstorage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/heilart1n/justpin-ipfs/file"
	httpretry "github.com/heilart1n/justpin-ipfs/http"
	"io"
	"net/http"
	"os"
)

func (client *Client) Name() string {
	return ClientName
}

// PinFile pins content to NFTStorage by providing a file path, it returns an IPFS
// hash and an error.
func (client *Client) PinFile(fp string) (string, error) {
	fi, err := os.Stat(fp)
	if err != nil {
		return "", err
	}

	// For regular file
	if fi.Mode().IsRegular() {
		f, err := os.Open(fp)
		if err != nil {
			return "", err
		}
		defer f.Close()

		return client.pinFile(f, file.MediaType(f))
	}

	// For directory, or etc
	f, err := file.NewSerialFile(fp)
	if err != nil {
		return "", err
	}

	mfr, err := file.CreateMultiForm(f, true)
	if err != nil {
		return "", err
	}
	boundary := "multipart/form-data; boundary=" + mfr.Boundary()

	return client.pinFile(mfr, boundary)
}

// PinWithReader pins content to NFTStorage by given io.Reader, it returns an IPFS hash and an error.
func (client *Client) PinWithReader(rd io.Reader) (string, error) {
	return client.pinFile(rd, file.MediaType(rd))
}

// PinWithBytes pins content to NFTStorage by given byte slice, it returns an IPFS hash and an error.
func (client *Client) PinWithBytes(buf []byte) (string, error) {
	return client.pinFile(bytes.NewReader(buf), file.MediaType(buf))
}

func (client *Client) pinFile(r io.Reader, boundary string) (string, error) {
	endpoint := APIUrl + "/upload"

	req, err := http.NewRequest(http.MethodPost, endpoint, r)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", boundary)
	req.Header.Add("Authorization", "Bearer "+client.cfg.Apikey)
	httpClient := httpretry.NewClient(client.Client)
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var out addEvent
	if err := json.Unmarshal(data, &out); err != nil {
		var e *json.SyntaxError
		if errors.As(err, &e) {
			return "", fmt.Errorf("json syntax error at byte offset %d", e.Offset)
		}
		return "", err
	}

	return out.Value.Cid, nil
}

// PinHash pins content to NFTStorage by giving an IPFS hash, it returns the result and an error.
// Note: unsupported
func (client *Client) PinHash(hash string) (bool, error) {
	return false, fmt.Errorf("not yet supported")
}

// PinDir pins a directory to the NFT.Storage pinning service.
// It alias to PinFile.
func (client *Client) PinDir(name string) (string, error) {
	return client.PinFile(name)
}

func (client *Client) Pin(path interface{}) (cid string, err error) {
	err = fmt.Errorf("unsupported pinner")
	switch v := path.(type) {
	case string:
		_, err = os.Lstat(v)
		if err != nil {
			return
		}
		cid, err = client.PinFile(v)
	case io.Reader:
		cid, err = client.PinWithReader(v)
	case []byte:
		cid, err = client.PinWithBytes(v)
	}
	if err != nil {
		err = fmt.Errorf("%s: %w", client.Name(), err)
	}
	return cid, err
}
