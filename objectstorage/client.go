package objectstorage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Error struct {
	APIError struct{
		Message string `json:"message,omitempty"`
	}
	StatusCode int
	Endpoint   string
}

func (e Error) Error() string {
	return fmt.Sprintf("API Error: %d %s %s", e.StatusCode, e.Endpoint, e.APIError.Message)
}

const (
	StorageAccountEndpoint string = "http://localhost:8084/"
)

type Client struct {
	Username string
	Password string
	AccessToken string
	HTTPClient *http.Client
}

func (c *Client) Do(method, endpoint string, payload *bytes.Buffer) (*http.Response, error) {
	absoluteendpoint := StorageAccountEndpoint + endpoint
	log.Printf("[DEBUG] Sending request to %s %s", method, absoluteendpoint)

	var bodyreader io.Reader

	if payload != nil {
		log.Printf("[DEBUG] with payload %s", payload.String())
		bodyreader = payload
	}

	req, err := http.NewRequest(method, absoluteendpoint, bodyreader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("token", c.AccessToken)
	if payload != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	req.Close = true

	resp, err := c.HTTPClient.Do(req)
	log.Printf("[DEBUG] Resp: %v Err: %v", resp, err)
	if resp.StatusCode >= 400 || resp.StatusCode < 200 {
		apiError := Error{
			StatusCode: resp.StatusCode,
			Endpoint: endpoint,
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Resp Body: %s", string(body))

		err = json.Unmarshal(body, &apiError)
		if err != nil {
			apiError.APIError.Message = string(body)
		}

		return resp, error(apiError)
	}
	return resp, err
}

func (c *Client) AuthenticateIfRequire() error {
	if c.AccessToken != "" {
		return nil
	}

	tokenRequest := c.newTokenRequest();

	var jsonbuffer []byte
	jsonpayload := bytes.NewBuffer(jsonbuffer)
	enc := json.NewEncoder(jsonpayload)
	enc.Encode(tokenRequest)

	resp, err := c.Do("POST", "token", jsonpayload)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		var tokenResponse TokenResponse

		body, readerr := ioutil.ReadAll(resp.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &tokenResponse)
		if decodeerr != nil {
			return decodeerr
		}

		c.AccessToken = tokenResponse.Token
	}

	return nil
}

func (c *Client) Get(endpoint string) (*http.Response, error) {
	err := c.AuthenticateIfRequire()
	if err != nil {
		log.Printf("[DEBUG] Err: %v", err)
		return nil, err
	}
	return c.Do("GET", endpoint, nil)
}

func (c *Client) Post(endpoint string, jsonpayload *bytes.Buffer) (*http.Response, error) {
	err := c.AuthenticateIfRequire()
	if err != nil {
		log.Printf("[DEBUG] Err: %v", err)
		return nil, err
	}
	return c.Do("POST", endpoint, jsonpayload)
}

func (c *Client) Put(endpoint string, jsonpayload *bytes.Buffer) (*http.Response, error) {
	err := c.AuthenticateIfRequire()
	if err != nil {
		log.Printf("[DEBUG] Err: %v", err)
		return nil, err
	}
	return c.Do("PUT", endpoint, jsonpayload)
}

func (c *Client) PutOnly(endpoint string) (*http.Response, error) {
	err := c.AuthenticateIfRequire()
	if err != nil {
		log.Printf("[DEBUG] Err: %v", err)
		return nil, err
	}
	return c.Do("PUT", endpoint, nil)
}

func (c *Client) Delete(endpoint string) (*http.Response, error) {
	err := c.AuthenticateIfRequire()
	if err != nil {
		log.Printf("[DEBUG] Err: %v", err)
		return nil, err
	}
	return c.Do("DELETE", endpoint, nil)
}

type TokenRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type TokenResponse struct {
	Token string `json:"token,omitempty"`
}

func (c *Client) newTokenRequest() *TokenRequest {
	return &TokenRequest{
		Username: c.Username,
		Password: c.Password,
	}
}