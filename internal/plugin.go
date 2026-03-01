package internal

import (
	"github.com/orchestra-mcp/plugin-integration-figma/internal/tools"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// FigmaPlugin registers all Figma API tools with the plugin builder.
type FigmaPlugin struct{}

// RegisterTools registers all 6 Figma tools on the given plugin builder.
func (fp *FigmaPlugin) RegisterTools(builder *plugin.PluginBuilder) {
	builder.RegisterTool("figma_get_file",
		"Get a Figma file by key",
		tools.FigmaGetFileSchema(), tools.FigmaGetFile())

	builder.RegisterTool("figma_get_components",
		"Get all components in a Figma file",
		tools.FigmaGetComponentsSchema(), tools.FigmaGetComponents())

	builder.RegisterTool("figma_get_styles",
		"Get all styles in a Figma file",
		tools.FigmaGetStylesSchema(), tools.FigmaGetStyles())

	builder.RegisterTool("figma_get_node",
		"Get a specific node from a Figma file",
		tools.FigmaGetNodeSchema(), tools.FigmaGetNode())

	builder.RegisterTool("figma_export_node",
		"Export a Figma node as an image (png, jpg, svg, or pdf)",
		tools.FigmaExportNodeSchema(), tools.FigmaExportNode())

	builder.RegisterTool("figma_sync_tokens",
		"Extract design tokens (colors, typography, spacing) from a Figma file",
		tools.FigmaSyncTokensSchema(), tools.FigmaSyncTokens())
}
