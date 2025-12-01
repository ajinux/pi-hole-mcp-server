package client

import (
	"context"
	"fmt"
)

type TopDeviceStats struct {
	Clients []struct {
		Ip    string `json:"ip"`
		Name  string `json:"name"`
		Count int    `json:"count"`
	} `json:"clients"`
	TotalQueries   int     `json:"total_queries"`
	BlockedQueries int     `json:"blocked_queries"`
	Took           float64 `json:"took"`
}

func (c *Client) GetTopActiveClientsByUsage(ctx context.Context, topN int) (*TopDeviceStats, error) {
	var res TopDeviceStats
	err := c.getJSON(ctx, fmt.Sprintf("stats/top_clients?count=%d", topN), &res)
	if err != nil {
		return nil, fmt.Errorf("failed to get top active clients by usage: %w", err)
	}
	return &res, nil
}
