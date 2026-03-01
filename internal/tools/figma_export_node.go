package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-integration-figma/internal/figma"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// FigmaExportNodeSchema returns the JSON Schema for the figma_export_node tool.
func FigmaExportNodeSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file_key": map[string]any{
				"type":        "string",
				"description": "The Figma file key",
			},
			"node_id": map[string]any{
				"type":        "string",
				"description": "The Figma node ID to export",
			},
			"format": map[string]any{
				"type":        "string",
				"description": "Export format: png, jpg, svg, or pdf (default: png)",
				"enum":        []any{"png", "jpg", "svg", "pdf"},
			},
			"scale": map[string]any{
				"type":        "number",
				"description": "Export scale multiplier (default: 1)",
			},
		},
		"required": []any{"file_key", "node_id"},
	})
	return s
}

// FigmaExportNode returns a handler that exports a Figma node as an image.
func FigmaExportNode() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "file_key", "node_id"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}
		fileKey := helpers.GetString(req.Arguments, "file_key")
		nodeId := helpers.GetString(req.Arguments, "node_id")
		format := helpers.GetStringOr(req.Arguments, "format", "png")
		scale := helpers.GetFloat64(req.Arguments, "scale")
		if scale <= 0 {
			scale = 1
		}
		path := fmt.Sprintf("images/%s?ids=%s&format=%s&scale=%.2g", fileKey, nodeId, format, scale)
		result, err := figma.NewClient().GetFormatted(ctx, path)
		if err != nil {
			return helpers.ErrorResult("figma_error", err.Error()), nil
		}
		return helpers.TextResult(result), nil
	}
}
