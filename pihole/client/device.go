package client

import (
	"context"
	"fmt"
)

type ConnectedDeviceInfo struct {
	Hwaddr    string   `json:"hwaddr"`
	MacVendor string   `json:"macVendor"`
	LastQuery UnixTime `json:"lastQuery"`
	Addresses string   `json:"addresses"`
	Names     *string  `json:"names"`
}

type AllConnectedDeviceInfo struct {
	Clients []ConnectedDeviceInfo `json:"clients"`
	Took    float64               `json:"took"`
}

func (c *Client) GetAllClients(ctx context.Context) (*AllConnectedDeviceInfo, error) {
	var res AllConnectedDeviceInfo
	err := c.getJSON(ctx, "clients/_suggestions", &res)
	if err != nil {
		return nil, fmt.Errorf("failed to get active clients: %w", err)
	}
	return &res, nil
}
