package client

import (
	"context"
	"fmt"
	"sort"
)

type Domain struct {
	Name    string `json:"domain"`
	Count   int    `json:"count"`
	Blocked bool   `json:"blocked"`
}

type TopDomainStats struct {
	Domains        []Domain `json:"domains"`
	TotalQueries   int      `json:"total_queries"`
	BlockedQueries int      `json:"blocked_queries"`
}

func (c *Client) GetTopDomainsQueried(ctx context.Context) ([]*Domain, error) {
	var topAllowedDomains TopDomainStats
	err := c.getJSON(ctx, "stats/top_domains", &topAllowedDomains)
	if err != nil {
		return nil, fmt.Errorf("error getting top domains queried: %w", err)
	}
	fmt.Printf("top allowed domains: %v\n", topAllowedDomains.Domains)
	var topBlockedDomains TopDomainStats
	err = c.getJSON(ctx, "stats/top_domains?blocked=true", &topBlockedDomains)
	if err != nil {
		return nil, fmt.Errorf("error getting top blocked domains queried: %w", err)
	}
	allDomains := make([]*Domain, 0, len(topAllowedDomains.Domains)+len(topBlockedDomains.Domains))
	for _, d := range topAllowedDomains.Domains {
		allDomains = append(allDomains, &d)
	}
	for _, d := range topBlockedDomains.Domains {
		d.Blocked = true
		allDomains = append(allDomains, &d)
	}

	sort.Slice(allDomains, func(i, j int) bool {
		return allDomains[i].Count > allDomains[j].Count
	})

	return allDomains, nil
}
