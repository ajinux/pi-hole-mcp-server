package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type topDomainsGlobalResponse struct {
	TotalDomains int          `json:"total_domains"`
	Domains      []domainStat `json:"domains"`
}

type domainStat struct {
	Domain  string `json:"domain"`
	Count   int    `json:"count"`
	Blocked bool   `json:"blocked"`
}

// registerTopDomains registers the tool for getting top queried domains globally
func (r *Registry) registerTopDomains(server *mcp.Server) {
	server.AddTool(&mcp.Tool{
		Name:        "get_top_domains",
		Description: "Get the top queried domains (both allowed and blocked) from Pi-hole",
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	}, r.withLogging("get_top_domains", r.handleTopDomains))
}

// handleTopDomains handles requests for the get_top_domains tool
func (r *Registry) handleTopDomains(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get top domains from Pi-hole
	domains, err := r.piholeClient.GetTopDomainsQueried(ctx)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to get top domains: %v", err),
				},
			},
		}, nil
	}

	// Build response
	var domainStats []domainStat
	for _, d := range domains {
		domainStats = append(domainStats, domainStat{
			Domain:  d.Name,
			Count:   d.Count,
			Blocked: d.Blocked,
		})
	}

	response := topDomainsGlobalResponse{
		TotalDomains: len(domainStats),
		Domains:      domainStats,
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
