package tools

import (
	"context"
	"log/slog"

	"github.com/ajinux/pi-hole-mcp-server/pihole/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Registry holds all available tools and the Pi-hole client
type Registry struct {
	piholeClient *client.Client
	logger       *slog.Logger
}

// NewRegistry creates a new tool registry with the given Pi-hole client
func NewRegistry(piholeClient *client.Client, logger *slog.Logger) *Registry {
	return &Registry{
		piholeClient: piholeClient,
		logger:       logger,
	}
}

// RegisterAll registers all available tools with the MCP server
func (r *Registry) RegisterAll(server *mcp.Server) {
	r.registerTopActiveClients(server)
	r.registerTopDomainsForClient(server)
	r.registerTopDomains(server)
	r.registerDNSRecords(server)
	r.registerWhoisLookup(server)

	// Register prompts
	r.registerDomainOSINTPrompt(server)
}

// ToolHandler is a function type for handling tool requests
type ToolHandler func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error)

// withLogging wraps a tool handler with automatic logging
func (r *Registry) withLogging(toolName string, handler ToolHandler) mcp.ToolHandler {
	return func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		r.logger.Info("Tool invoked",
			"tool", toolName,
			"arguments", string(request.Params.Arguments),
		)

		result, err := handler(ctx, request)

		if err != nil {
			r.logger.Error("Tool handler error",
				"tool", toolName,
				"error", err,
			)
			return result, err
		}

		if result != nil && result.IsError {
			r.logger.Error("Tool execution failed",
				"tool", toolName,
			)
		} else {
			r.logger.Info("Tool executed successfully",
				"tool", toolName,
			)
		}

		return result, err
	}
}
