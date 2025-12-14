package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/ajinux/pi-hole-mcp-server/config"
	"github.com/ajinux/pi-hole-mcp-server/pihole/client"
	"github.com/ajinux/pi-hole-mcp-server/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Build-time variables injected via ldflags during compilation
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	log.Printf(
		"starting pihole-mcp server version=%s commit=%s built=%s",
		version, commit, date,
	)

	// Load configuration
	cfg := config.Load()

	// Create Pi-hole client
	ctx := context.Background()
	piholeClient, err := client.NewClient(ctx, cfg.PiHoleURL, cfg.PiHolePassword)
	if err != nil {
		log.Fatalf("Failed to create Pi-hole client: %v", err)
	}

	// Create MCP server
	mserv := mcp.NewServer(&mcp.Implementation{
		Name:    "Pi hole mcp server",
		Title:   "Pi hole MCP Server",
		Version: "1.0",
	}, nil)

	// Create logger for MCP connections
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Register all tools
	toolRegistry := tools.NewRegistry(piholeClient, logger)
	toolRegistry.RegisterAll(mserv)

	// Create StreamableHTTP handler that returns our MCP server
	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		logger.Info("New MCP client connection",
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"session_id", r.Header.Get("Mcp-Session-Id"),
		)
		return mserv
	}, &mcp.StreamableHTTPOptions{
		Logger: logger,
	})

	log.Printf("Starting Pi-hole MCP server on http://localhost:%s", cfg.Port)
	log.Printf("Connected to Pi-hole at: %s", cfg.PiHoleURL)

	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
