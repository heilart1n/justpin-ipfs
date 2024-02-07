package pinata

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/heilart1n/justpin-ipfs/file"
	httpretry "github.com/heilart1n/justpin-ipfs/http"
	"github.com/heilart1n/justpin-ipfs/pinners"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func (client *Client) Name() string {
	return ClientName
}

// PinFile pins content to Pinata by providing a file path, it returns an IPFS
// hash and an error.
func (client *Client) PinFile(fp string) (pinners.Result, error) {
	f, err := file.NewSerialFile(fp)
	if err != nil {
		return nil, err
	}
	f.MapDirectory(filepath.Base(fp))

	mfr, err := file.CreateMultiForm(f, true)
	if err != nil {
		return nil, err
	}
	boundary := "multipart/form-data; boundary=" + mfr.Boundary()

	return client.pinFile(mfr, boundary)
}

// PinWithReader pins content to Pinata by given io.Reader, it returns an IPFS hash and an error.
func (client *Client) PinWithReader(rd io.Reader) (pinners.Result, error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	fn := file.RandString(6, "lower")

	go func() {
		defer w.Close()
		defer m.Close()

		part, err := m.CreateFormFile("file", fn)
		if err != nil {
			return
		}

		if _, err = io.Copy(part, rd); err != nil {
			return
		}
	}()

	return client.pinFile(r, m.FormDataContentType())
}

// PinWithBytes pins content to Infura by given byte slice, it returns an IPFS hash and an error.
func (client *Client) PinWithBytes(buf []byte) (pinners.Result, error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	fn := file.RandString(6, "lower")

	go func() {
		defer w.Close()
		defer m.Close()

		// m.WriteField("pinataOptions", `{cidVersion: 1}`)
		part, err := m.CreateFormFile("file", fn)
		if err != nil {
			return
		}

		if _, err = part.Write(buf); err != nil {
			return
		}
	}()

	return client.pinFile(r, m.FormDataContentType())
}

func (client *Client) pinFile(r io.Reader, boundary string) (pinners.Result, error) {
	// if fr, ok := r.(*file.MultiFileReader); ok {
	// 	// Metadata part.
	// 	metadataHeader := textproto.MIMEHeader{}
	// 	metadataHeader.Set("Content-Disposition", `form-data; name="pinataMetadata"`)
	// 	// Metadata content.
	// 	metadata := fmt.Sprintf(`{"name":"%s"}`, "adsasdfa")
	// 	fr.Write(metadataHeader, []byte(metadata))

	// 	// options part.
	// 	optsHeader := textproto.MIMEHeader{}
	// 	optsHeader.Set("Content-Disposition", `form-data; name="pinataOptions"`)
	// 	// options content.
	// 	opts := `{"cidVersion":"1","wrapWithDirectory":false}`
	// 	fr.Write(optsHeader, []byte(opts))
	// }
	req, err := http.NewRequest(http.MethodPost, PinFileUrl, r)
	if err != nil {
		return nil, err
	}
	client.setAuth(req)
	req.Header.Add("Content-Type", boundary)

	httpClient := httpretry.NewClient(client.Client)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
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

	return client.NewResult(out.IpfsHash), nil
}

// PinHash pins content to Pinata by giving an IPFS hash, it returns the result and an error.
func (client *Client) PinHash(hash string) (bool, error) {
	if hash == "" {
		return false, fmt.Errorf("invalid hash: %s", hash)
	}

	jsonValue, _ := json.Marshal(map[string]string{"hashToPin": hash})

	req, err := http.NewRequest(http.MethodPost, PinHashUrl, bytes.NewBuffer(jsonValue))
	if err != nil {
		return false, err
	}
	client.setAuth(req)

	req.Header.Set("Content-Type", "application/json")

	httpClient := httpretry.NewClient(client.Client)
	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf(resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(data, &dat); err != nil {
		var e *json.SyntaxError
		if errors.As(err, &e) {
			return false, fmt.Errorf("json syntax error at byte offset %d", e.Offset)
		}
		return false, err
	}

	if h, ok := dat["hashToPin"].(string); ok {
		return h == hash, nil
	}

	return false, fmt.Errorf("pin hash to Pinata failed")
}

// PinDir pins a directory to the Pinata pinning service.
// It alias to PinFile.
func (client *Client) PinDir(name string) (pinners.Result, error) {
	return client.PinFile(name)
}

func (client *Client) setAuth(req *http.Request) {
	if client.cfg.Secret != "" && client.cfg.Apikey != "" {
		req.Header.Add("pinata_secret_api_key", client.cfg.Secret)
		req.Header.Add("pinata_api_key", client.cfg.Apikey)
	} else {
		req.Header.Add("Authorization", "Bearer "+client.cfg.Apikey)
	}
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
