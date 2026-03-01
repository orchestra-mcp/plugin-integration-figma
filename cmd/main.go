// Command integration-figma is the entry point for the integration.figma
// plugin binary. It provides 6 MCP tools for interacting with the Figma API.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/orchestra-mcp/plugin-integration-figma/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

func main() {
	builder := plugin.New("integration.figma").
		Version("0.1.0").
		Description("Figma REST API integration for design-to-code workflows").
		Author("Orchestra").
		Binary("integration-figma")

	tp := &internal.FigmaPlugin{}
	tp.RegisterTools(builder)

	p := builder.BuildWithTools()
	p.ParseFlags()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	if err := p.Run(ctx); err != nil {
		log.Fatalf("integration.figma: %v", err)
	}
}
