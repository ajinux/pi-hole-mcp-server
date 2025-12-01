package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type domainInfo struct {
	Domain        string `json:"domain"`
	QueryCount    int    `json:"query_count"`
	RejectedCount int    `json:"rejected_count"`
}

type topDomainsResponse struct {
	ClientIP        string       `json:"client_ip"`
	HoursAnalyzed   int          `json:"hours_analyzed"`
	TotalQueries    int          `json:"total_queries"`
	RejectedQueries int          `json:"rejected_queries"`
	Domains         []domainInfo `json:"domains"`
}

// registerTopDomainsForClient registers the tool for getting top domains queried by a client
func (r *Registry) registerTopDomainsForClient(server *mcp.Server) {
	server.AddTool(&mcp.Tool{
		Name:        "get_top_domains_for_client",
		Description: "Get the top N most queried domains by a specific client IP address in the last X hours from Pi-hole",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"client_ip": map[string]interface{}{
					"type":        "string",
					"description": "The IP address of the client",
				},
				"hours": map[string]interface{}{
					"type":        "number",
					"description": "Number of hours to look back (default: 24, max: 168 for 1 week)",
					"minimum":     1,
					"maximum":     168,
				},
				"count": map[string]interface{}{
					"type":        "number",
					"description": "Number of top domains to return (default: 10)",
					"minimum":     1,
					"maximum":     100,
				},
			},
			"required": []string{"client_ip"},
		},
	}, r.withLogging("get_top_domains_for_client", r.handleTopDomainsForClient))
}

// handleTopDomainsForClient handles requests for the get_top_domains_for_client tool
func (r *Registry) handleTopDomainsForClient(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse arguments from JSON
	var args struct {
		ClientIP string  `json:"client_ip"`
		Hours    float64 `json:"hours"`
		Count    int     `json:"count"`
	}

	// Set defaults
	args.Hours = 24
	args.Count = 10

	// Parse arguments
	if len(request.Params.Arguments) > 0 {
		if err := json.Unmarshal(request.Params.Arguments, &args); err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("Failed to parse arguments: %v", err),
					},
				},
			}, nil
		}
	}

	// Validate client_ip is provided
	if args.ClientIP == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "client_ip is required",
				},
			},
		}, nil
	}

	// Validate hours range
	if args.Hours < 1 {
		args.Hours = 1
	} else if args.Hours > 168 {
		args.Hours = 168
	}

	// Validate count range
	if args.Count < 1 {
		args.Count = 1
	} else if args.Count > 100 {
		args.Count = 100
	}

	// Calculate time range
	until := time.Now().Add(-time.Duration(args.Hours) * time.Hour)

	// Get DNS queries for the client
	queries, err := r.piholeClient.GetDNSQueriesForClient(ctx, args.ClientIP, until)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to get DNS queries for client: %v", err),
				},
			},
		}, nil
	}

	// Aggregate domains
	domainCounts := make(map[string]int)
	domainRejectedCounts := make(map[string]int)
	rejectedCount := 0

	for _, query := range queries {
		domainCounts[query.Domain]++
		if query.Status == "GRAVITY" {
			rejectedCount++
			domainRejectedCounts[query.Domain]++
		}
	}

	// Convert to slice and sort by count
	var domains []domainInfo
	for domain, count := range domainCounts {
		domains = append(domains, domainInfo{
			Domain:        domain,
			QueryCount:    count,
			RejectedCount: domainRejectedCounts[domain],
		})
	}

	sort.Slice(domains, func(i, j int) bool {
		return domains[i].QueryCount > domains[j].QueryCount
	})

	// Limit to top N
	if len(domains) > args.Count {
		domains = domains[:args.Count]
	}

	// Build response
	response := topDomainsResponse{
		ClientIP:        args.ClientIP,
		HoursAnalyzed:   int(args.Hours),
		TotalQueries:    len(queries),
		RejectedQueries: rejectedCount,
		Domains:         domains,
	}

	// Format response as JSON
	resultJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to marshal response: %v", err),
				},
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(resultJSON),
			},
		},
	}, nil
}
