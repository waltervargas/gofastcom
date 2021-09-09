package fastclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

type Client struct {
	baseURL string
	apiURL string
	maxRequests int
	maxTime int
	oca *OCA
	userAgent string
}

type OCA struct {
	Client struct {
		IP       string `json:"ip"`
		Asn      string `json:"asn"`
		Location struct {
			City    string `json:"city"`
			Country string `json:"country"`
		} `json:"location"`
	} `json:"client"`
	Targets []struct {
		Name     string `json:"name"`
		URL      string `json:"url"`
		Location struct {
			City    string `json:"city"`
			Country string `json:"country"`
		} `json:"location"`
	} `json:"targets"`
}

var (
	errScriptNotFound = errors.New("unable to get script path")
	errTokenNotFound = errors.New("unable to get token")

	oca OCA
)

const (
	baseURL    = "https://fast.com"
	userAgent  = "github.com/waltervargas/gofastcom"
	apiURL     = "https://api.fast.com/netflix/speedtest/v2"
)

func New(maxReq, maxTime int) (*Client, error) {
	c := Client{
		baseURL: baseURL,
		userAgent: userAgent,
		maxRequests: maxReq,
		maxTime: maxTime,
	}

	t, err := getToken(baseURL)
	if err != nil {
		return nil, err
	}
	c.apiURL = fmt.Sprintf("%s?https=true&token=%s", apiURL, t)

	oca, err := c.getOCAs()
	if err != nil {
		return nil, err
	}
	c.oca = oca

	return &c, nil
}

func (c Client) getOCAs() (*OCA, error) {
	var oca OCA

	resp, err := http.Get(c.apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &oca)
	if err != nil {
		return nil, err
	}

	return &oca, nil
}

func getToken(baseURL string) ([]byte, error) {
	baseResp, err := http.Get(baseURL)
	if err != nil {
		return nil, err
	}
	defer baseResp.Body.Close()

	path, err := getScriptPath(baseResp.Body)
	if err != nil {
		return nil, err
	}

	scriptResp, err := http.Get(fmt.Sprintf("%s%s", baseURL, path))
	if err != nil {
		return nil, err
	}
	defer scriptResp.Body.Close()

	scriptBody, err := io.ReadAll(scriptResp.Body)
	if err != nil {
		return nil, err
	}

	t := getTokenFromScriptBody(scriptBody)
	if t == nil {
		return nil, errTokenNotFound
	}

	return t, nil
}

func getTokenFromScriptBody(b []byte) []byte {
	match := []byte("https:!0,token:\"")
	if bytes.Contains(b, match) {
		index := bytes.Index(b, match)
		if len(b[index+len(match):]) > 32 {
			return b[index+len(match):index+len(match)+32]
		}
	}
	return nil
}

func getScriptPath(r io.Reader) ([]byte, error) {
	z := html.NewTokenizer(r)

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if z.Err().Error() != "EOF" {
				return nil, z.Err()
			}
			return nil, errScriptNotFound
		}

		if tt == html.EndTagToken {
			continue
		}

		n, hasAttr := z.TagName()
		if hasAttr && len(n) == 6 && n[0] == 's' && n[5] == 't' {
			k, v, _ := z.TagAttr()
			if len(k) == 3 && k[0] == 's' {
				return v, nil
			}
		}
	}
}
