package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerDomainOSINTPrompt registers a prompt that guides the LLM to perform comprehensive OSINT on a domain
func (r *Registry) registerDomainOSINTPrompt(server *mcp.Server) {
	server.AddPrompt(&mcp.Prompt{
		Name:        "domain-osint",
		Title:       "Domain OSINT Analysis",
		Description: "Perform comprehensive Open Source Intelligence (OSINT) gathering on a domain using DNS and WHOIS lookups",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "domain",
				Description: "The domain or subdomain to investigate (e.g., example.com or api.example.com)",
				Required:    true,
			},
		},
	}, r.handleDomainOSINTPrompt)
}

// handleDomainOSINTPrompt handles the domain-osint prompt request
func (r *Registry) handleDomainOSINTPrompt(ctx context.Context, request *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	// Extract the domain argument
	domain := request.Params.Arguments["domain"]
	if domain == "" {
		return nil, fmt.Errorf("domain argument is required")
	}

	// Create a comprehensive OSINT prompt for the LLM
	promptText := fmt.Sprintf(`Perform comprehensive Open Source Intelligence (OSINT) analysis on the domain: %s

Please follow these steps:

1. **DNS Records Analysis**
   - Use the 'get_domain_dns_records' tool to retrieve DNS information
   - Analyze the following records:
     * A records (IPv4 addresses) - Identify hosting infrastructure
     * AAAA records (IPv6 addresses) - Check for IPv6 support
     * NS records (Name servers) - Identify DNS providers
     * MX records (Mail servers) - Identify email infrastructure
     * TXT records - Look for SPF, DKIM, DMARC, domain verification, and other configurations

2. **WHOIS Information**
   - Use the 'get_domain_whois' tool to gather registration details
   - Analyze:
     * Domain registration and expiration dates (identify domain age and renewal status)
     * Registrar information (who manages the domain)
     * Registrant details (organization, country, contact information)
     * Administrative, technical, and billing contacts
     * Domain status flags (look for locks, holds, or other restrictions)
     * Name servers listed in WHOIS

3. **Security & Infrastructure Assessment**
   - Identify the hosting provider based on IP addresses
   - Check if domain uses CDN or DDoS protection services
   - Analyze email security configurations (SPF, DMARC, DKIM in TXT records)
   - Note any privacy/proxy protection services being used
   - Identify potential sister domains or related infrastructure

4. **Risk Indicators**
   - Recently registered domains (potential phishing indicator)
   - Privacy-protected WHOIS (could indicate legitimate privacy or malicious hiding)
   - Mismatched registrar and hosting locations
   - Lack of email security records
   - Suspicious TXT records or DNS configurations

5. **Summary Report**
   Provide a structured summary including:
   - Domain ownership and registration timeline
   - Technical infrastructure overview
   - Security posture assessment
   - Any notable findings or red flags
   - Recommendations if applicable

Please be thorough and present your findings in a clear, organized manner.`, domain)

	return &mcp.GetPromptResult{
		Description: fmt.Sprintf("OSINT analysis prompt for domain: %s", domain),
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: promptText},
			},
		},
	}, nil
}
