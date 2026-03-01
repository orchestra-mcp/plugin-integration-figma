package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-integration-figma/internal/figma"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// FigmaGetNodeSchema returns the JSON Schema for the figma_get_node tool.
func FigmaGetNodeSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file_key": map[string]any{
				"type":        "string",
				"description": "The Figma file key",
			},
			"node_id": map[string]any{
				"type":        "string",
				"description": "The Figma node ID",
			},
		},
		"required": []any{"file_key", "node_id"},
	})
	return s
}

// FigmaGetNode returns a handler that fetches a specific node from a Figma file.
func FigmaGetNode() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "file_key", "node_id"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}
		fileKey := helpers.GetString(req.Arguments, "file_key")
		nodeId := helpers.GetString(req.Arguments, "node_id")
		result, err := figma.NewClient().GetFormatted(ctx, "files/"+fileKey+"/nodes?ids="+nodeId)
		if err != nil {
			return helpers.ErrorResult("figma_error", err.Error()), nil
		}
		return helpers.TextResult(result), nil
	}
}
