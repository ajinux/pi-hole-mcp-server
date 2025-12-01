package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ajinux/pi-hole-mcp-server/dnsclient"
	"github.com/ajinux/pi-hole-mcp-server/domain"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type whoisResponse struct {
	Domain            string   `json:"domain"`
	Status            []string `json:"status,omitempty"`
	NameServers       []string `json:"name_servers,omitempty"`
	CreatedDate       string   `json:"created_date,omitempty"`
	UpdatedDate       string   `json:"updated_date,omitempty"`
	ExpirationDate    string   `json:"expiration_date,omitempty"`
	RegistrarName     string   `json:"registrar_name,omitempty"`
	RegistrarEmail    string   `json:"registrar_email,omitempty"`
	RegistrarPhone    string   `json:"registrar_phone,omitempty"`
	RegistrantOrg     string   `json:"registrant_organization,omitempty"`
	RegistrantCountry string   `json:"registrant_country,omitempty"`
	RegistrantEmail   string   `json:"registrant_email,omitempty"`
	AdminEmail        string   `json:"admin_email,omitempty"`
	TechnicalEmail    string   `json:"technical_email,omitempty"`
	BillingEmail      string   `json:"billing_email,omitempty"`
}

// registerWhoisLookup registers the tool for performing WHOIS lookups
func (r *Registry) registerWhoisLookup(server *mcp.Server) {
	server.AddTool(&mcp.Tool{
		Name:        "get_domain_whois",
		Description: "Perform a WHOIS lookup on a domain to get registration information (registrar, creation date, expiration date, registrant details, etc.). Automatically extracts top-level domain if a subdomain is provided.",
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
	}, r.withLogging("get_domain_whois", r.handleWhoisLookup))
}

// handleWhoisLookup handles requests for the get_domain_whois tool
func (r *Registry) handleWhoisLookup(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	// Strip to top-level domain
	tld := dnsclient.ExtractTopLevelDomain(args.Domain)

	// Perform WHOIS lookup
	whoisInfo, err := domain.Whois(tld)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to perform WHOIS lookup: %v", err),
				},
			},
		}, nil
	}

	// Extract important OSINT fields with nil checks
	response := whoisResponse{
		Domain: tld,
	}

	// Domain information
	if whoisInfo.Domain != nil {
		response.Status = whoisInfo.Domain.Status
		response.NameServers = whoisInfo.Domain.NameServers
		response.CreatedDate = whoisInfo.Domain.CreatedDate
		response.UpdatedDate = whoisInfo.Domain.UpdatedDate
		response.ExpirationDate = whoisInfo.Domain.ExpirationDate
	}

	// Registrar information
	if whoisInfo.Registrar != nil {
		response.RegistrarName = whoisInfo.Registrar.Name
		response.RegistrarEmail = whoisInfo.Registrar.Email
		response.RegistrarPhone = whoisInfo.Registrar.Phone
	}

	// Registrant information
	if whoisInfo.Registrant != nil {
		response.RegistrantOrg = whoisInfo.Registrant.Organization
		response.RegistrantCountry = whoisInfo.Registrant.Country
		response.RegistrantEmail = whoisInfo.Registrant.Email
	}

	// Contact information
	if whoisInfo.Administrative != nil {
		response.AdminEmail = whoisInfo.Administrative.Email
	}
	if whoisInfo.Technical != nil {
		response.TechnicalEmail = whoisInfo.Technical.Email
	}
	if whoisInfo.Billing != nil {
		response.BillingEmail = whoisInfo.Billing.Email
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
