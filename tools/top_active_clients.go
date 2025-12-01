package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ajinux/pi-hole-mcp-server/pihole/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type clientInfo struct {
	Ip                 string   `json:"ip"`
	Name               string   `json:"name"`
	DnsRequestsCount   int      `json:"dns_requests_count"`
	MacAddress         []string `json:"mac_address"`
	MacVendor          string   `json:"mac_vendor"`
	LastRequestAgoMins int      `json:"last_request_ago_mins"`
}

type activeClientsResponse struct {
	Clients []clientInfo `json:"clients"`
}

// registerTopActiveClients registers the tool for getting top active clients
func (r *Registry) registerTopActiveClients(server *mcp.Server) {
	server.AddTool(&mcp.Tool{
		Name:        "get_top_active_clients",
		Description: "Get the top N most active clients by DNS query usage from Pi-hole",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"count": map[string]interface{}{
					"type":        "number",
					"description": "Number of top clients to return (default: 10)",
					"minimum":     1,
					"maximum":     100,
				},
			},
		},
	}, r.withLogging("get_top_active_clients", r.handleTopActiveClients))
}

// handleTopActiveClients handles requests for the get_top_active_clients tool
func (r *Registry) handleTopActiveClients(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse arguments from JSON
	var args struct {
		Count int `json:"count"`
	}

	// Set default count
	args.Count = 10

	// Parse arguments if provided
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

	// Validate count range
	if args.Count < 1 {
		args.Count = 1
	} else if args.Count > 100 {
		args.Count = 100
	}

	// Call Pi-hole API
	stats, err := r.piholeClient.GetTopActiveClientsByUsage(ctx, args.Count)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to get top active clients: %v", err),
				},
			},
		}, nil
	}

	allClients, err := r.piholeClient.GetAllClients(ctx)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to get all clients: %v", err),
				},
			},
		}, nil
	}

	// Map IPs to Names
	ipToDeviceInfo := make(map[string]client.ConnectedDeviceInfo, len(allClients.Clients))
	for _, device := range allClients.Clients {
		ipToDeviceInfo[device.Addresses] = device
	}

	// Enrich stats with MAC and Names
	var response activeClientsResponse
	for _, clientStat := range stats.Clients {
		clientInfo := clientInfo{
			Ip:               clientStat.Ip,
			Name:             clientStat.Name,
			DnsRequestsCount: clientStat.Count,
		}

		if deviceInfo, exists := ipToDeviceInfo[clientStat.Ip]; exists {
			clientInfo.MacAddress = strings.Split(deviceInfo.Hwaddr, ",")
			clientInfo.MacVendor = deviceInfo.MacVendor
			if deviceInfo.LastQuery.Time.Unix() > 0 {
				clientInfo.LastRequestAgoMins = int(time.Since(deviceInfo.LastQuery.Time).Minutes())
			}
		}

		response.Clients = append(response.Clients, clientInfo)
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
