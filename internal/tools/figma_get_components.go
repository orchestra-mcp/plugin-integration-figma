package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-integration-figma/internal/figma"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// FigmaGetComponentsSchema returns the JSON Schema for the figma_get_components tool.
func FigmaGetComponentsSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file_key": map[string]any{
				"type":        "string",
				"description": "The Figma file key",
			},
		},
		"required": []any{"file_key"},
	})
	return s
}

// FigmaGetComponents returns a handler that fetches all components in a Figma file.
func FigmaGetComponents() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "file_key"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}
		fileKey := helpers.GetString(req.Arguments, "file_key")
		result, err := figma.NewClient().GetFormatted(ctx, "files/"+fileKey+"/components")
		if err != nil {
			return helpers.ErrorResult("figma_error", err.Error()), nil
		}
		return helpers.TextResult(result), nil
	}
}
