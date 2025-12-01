package client

import (
	"context"
	"fmt"
	"log"
)

const (
	loginEndpoint = "auth"
)

type authRequest struct {
	Password string `json:"password"`
}

type authResponse struct {
	Session struct {
		Valid    bool   `json:"valid"`
		Totp     bool   `json:"totp"`
		Sid      string `json:"sid"`
		Csrf     string `json:"csrf"`
		Validity int    `json:"validity"`
		Message  string `json:"message"`
	} `json:"session"`
	Took float64 `json:"took"`
}

func (c *Client) obtainSessionID(ctx context.Context) (string, error) {
	res := &authResponse{}
	err := c.postJSON(ctx, loginEndpoint, authRequest{Password: c.password}, res)
	if err != nil {
		return "", err
	}
	if !res.Session.Valid {
		return "", fmt.Errorf("authentication failed: %s", res.Session.Message)
	}
	log.Printf("obtained a new pi-hole session token successfully")
	return res.Session.Sid, nil
}
