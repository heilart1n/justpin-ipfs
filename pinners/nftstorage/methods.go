package nftstorage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/heilart1n/justpin-ipfs/file"
	httpretry "github.com/heilart1n/justpin-ipfs/http"
	"github.com/heilart1n/justpin-ipfs/pinners"
	"io"
	"net/http"
	"os"
)

func (client *Client) Name() string {
	return ClientName
}

// PinFile pins content to NFTStorage by providing a file path, it returns an IPFS
// hash and an error.
func (client *Client) PinFile(fp string) (pinners.Result, error) {
	fi, err := os.Stat(fp)
	if err != nil {
		return nil, err
	}

	// For regular file
	if fi.Mode().IsRegular() {
		f, err := os.Open(fp)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		return client.pinFile(f, file.MediaType(f))
	}

	// For directory, or etc
	f, err := file.NewSerialFile(fp)
	if err != nil {
		return nil, err
	}

	mfr, err := file.CreateMultiForm(f, true)
	if err != nil {
		return nil, err
	}
	boundary := "multipart/form-data; boundary=" + mfr.Boundary()

	return client.pinFile(mfr, boundary)
}

// PinWithReader pins content to NFTStorage by given io.Reader, it returns an IPFS hash and an error.
func (client *Client) PinWithReader(rd io.Reader) (pinners.Result, error) {
	return client.pinFile(rd, file.MediaType(rd))
}

// PinWithBytes pins content to NFTStorage by given byte slice, it returns an IPFS hash and an error.
func (client *Client) PinWithBytes(buf []byte) (pinners.Result, error) {
	return client.pinFile(bytes.NewReader(buf), file.MediaType(buf))
}

func (client *Client) pinFile(r io.Reader, boundary string) (pinners.Result, error) {
	endpoint := APIUrl + "/upload"

	req, err := http.NewRequest(http.MethodPost, endpoint, r)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", boundary)
	req.Header.Add("Authorization", "Bearer "+client.cfg.Apikey)
	httpClient := httpretry.NewClient(client.Client)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out addEvent
	if err := json.Unmarshal(data, &out); err != nil {
		var e *json.SyntaxError
		if errors.As(err, &e) {
			return nil, fmt.Errorf("json syntax error at byte offset %d", e.Offset)
		}
		return nil, err
	}

	return newResult(out.Value.Cid), nil
}

// PinHash pins content to NFTStorage by giving an IPFS hash, it returns the result and an error.
// Note: unsupported
func (client *Client) PinHash(hash string) (bool, error) {
	return false, fmt.Errorf("not yet supported")
}

// PinDir pins a directory to the NFT.Storage pinning service.
// It alias to PinFile.
func (client *Client) PinDir(name string) (pinners.Result, error) {
	return client.PinFile(name)
}

func (client *Client) Pin(path interface{}) (result pinners.Result, err error) {
	err = fmt.Errorf("unsupported pinner")
	switch v := path.(type) {
	case string:
		_, err = os.Lstat(v)
		if err != nil {
			return
		}
		result, err = client.PinFile(v)
	case io.Reader:
		result, err = client.PinWithReader(v)
	case []byte:
		result, err = client.PinWithBytes(v)
	}
	if err != nil {
		err = fmt.Errorf("%s: %w", client.Name(), err)
	}
	return result, err
}
