package client

import (
	"context"
	"fmt"
	"time"
)

type DNSQuery struct {
	Time   time.Time
	Type   string
	Status string
	Domain string
}

type DNSQueries struct {
	Queries []struct {
		Id       int     `json:"id"`
		Time     float64 `json:"time"`
		Type     string  `json:"type"`
		Status   string  `json:"status"`
		Dnssec   string  `json:"dnssec"`
		Domain   string  `json:"domain"`
		Upstream *string `json:"upstream"`
		Reply    struct {
			Type string  `json:"type"`
			Time float64 `json:"time"`
		} `json:"reply"`
		Client struct {
			Ip   string      `json:"ip"`
			Name interface{} `json:"name"`
		} `json:"client"`
		ListId *int `json:"list_id"`
		Ede    struct {
			Code int         `json:"code"`
			Text interface{} `json:"text"`
		} `json:"ede"`
		Cname interface{} `json:"cname"`
	} `json:"queries"`
	Cursor          int     `json:"cursor"`
	RecordsTotal    int     `json:"recordsTotal"`
	RecordsFiltered int     `json:"recordsFiltered"`
	Draw            int     `json:"draw"`
	Took            float64 `json:"took"`
}

func (c *Client) GetDNSQueriesForClient(ctx context.Context, clientIP string, until time.Time) ([]DNSQuery, error) {
	var allQueries []DNSQuery

	untilUnix := float64(until.Unix())
	start := 0
	length := 100

	for {
		var res DNSQueries
		url := fmt.Sprintf("queries?client_ip=%s&start=%d&length=%d", clientIP, start, length)
		err := c.getJSON(ctx, url, &res)
		if err != nil {
			return nil, fmt.Errorf("error getting dns queries for the client: %w", err)
		}

		// If no queries returned, we're done
		if len(res.Queries) == 0 {
			break
		}

		// Process queries and stop if we reach the until timestamp
		for _, query := range res.Queries {
			if query.Time < untilUnix {
				// Reached queries older than until timestamp, stop pagination
				return allQueries, nil
			}
			allQueries = append(allQueries, DNSQuery{
				Time:   time.Unix(int64(query.Time), 0),
				Type:   query.Type,
				Status: query.Status,
				Domain: query.Domain,
			})
		}

		// Update pagination parameters
		start += length

		// If we've fetched all records, stop
		if start >= res.RecordsFiltered {
			break
		}
	}

	return allQueries, nil
}
