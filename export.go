package integrationfigma

import (
	"github.com/orchestra-mcp/plugin-integration-figma/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// Register adds all Figma tools to the builder.
func Register(builder *plugin.PluginBuilder) {
	fp := &internal.FigmaPlugin{}
	fp.RegisterTools(builder)
}
