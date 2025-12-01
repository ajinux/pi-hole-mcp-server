package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ajinux/pi-hole-mcp-server/dnsclient"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type dnsRecordResponse struct {
	Domain string   `json:"domain"`
	A      []string `json:"a_records"`
	AAAA   []string `json:"aaaa_records"`
	NS     []string `json:"ns_records"`
	MX     []string `json:"mx_records"`
	TXT    []string `json:"txt_records"`
}

// registerDNSRecords registers the tool for getting DNS records for a domain
func (r *Registry) registerDNSRecords(server *mcp.Server) {
	server.AddTool(&mcp.Tool{
		Name:        "get_domain_dns_records",
		Description: "Get DNS records (A, AAAA, NS, MX, TXT) for a domain. Automatically extracts top-level domain if a subdomain is provided.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"domain": map[string]interface{}{
					"type":        "string",
					"description": "The domain or subdomain to query (e.g., example.com or api.example.com)",
				},
			},
			"required": []string{"domain"},
		},
	}, r.withLogging("get_domain_dns_records", r.handleDNSRecords))
}

// handleDNSRecords handles requests for the get_domain_dns_records tool
func (r *Registry) handleDNSRecords(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse arguments from JSON
	var args struct {
		Domain string `json:"domain"`
	}

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

	// Validate domain is provided
	if args.Domain == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "domain is required",
				},
			},
		}, nil
	}

	// Get DNS records (this will automatically strip to TLD)
	records, err := dnsclient.GetAllRecords(args.Domain)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to get DNS records: %v", err),
				},
			},
		}, nil
	}

	// Build response
	response := dnsRecordResponse{
		Domain: records.Domain,
		A:      records.A,
		AAAA:   records.AAAA,
		NS:     records.NS,
		MX:     records.MX,
		TXT:    records.TXT,
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
