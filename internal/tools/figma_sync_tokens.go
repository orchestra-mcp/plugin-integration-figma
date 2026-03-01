package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-integration-figma/internal/figma"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// FigmaSyncTokensSchema returns the JSON Schema for the figma_sync_tokens tool.
func FigmaSyncTokensSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file_key": map[string]any{
				"type":        "string",
				"description": "The Figma file key",
			},
			"output_path": map[string]any{
				"type":        "string",
				"description": "Optional path to write the design tokens JSON file",
			},
		},
		"required": []any{"file_key"},
	})
	return s
}

// FigmaSyncTokens returns a handler that extracts design tokens from a Figma file.
func FigmaSyncTokens() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "file_key"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}
		fileKey := helpers.GetString(req.Arguments, "file_key")
		outputPath := helpers.GetString(req.Arguments, "output_path")

		// Fetch file styles.
		body, err := figma.NewClient().Get(ctx, "files/"+fileKey+"/styles")
		if err != nil {
			return helpers.ErrorResult("figma_error", err.Error()), nil
		}

		// Parse the styles response.
		var stylesResp map[string]any
		if err := json.Unmarshal(body, &stylesResp); err != nil {
			return helpers.ErrorResult("parse_error", fmt.Sprintf("failed to parse styles: %v", err)), nil
		}

		// Extract tokens by style type.
		tokens := map[string]any{
			"colors":      map[string]any{},
			"typography":  map[string]any{},
			"spacing":     map[string]any{},
		}

		if meta, ok := stylesResp["meta"].(map[string]any); ok {
			if styles, ok := meta["styles"].([]any); ok {
				for _, styleRaw := range styles {
					style, ok := styleRaw.(map[string]any)
					if !ok {
						continue
					}
					name, _ := style["name"].(string)
					styleType, _ := style["style_type"].(string)
					nodeID, _ := style["node_id"].(string)

					tokenKey := strings.ReplaceAll(strings.ToLower(name), "/", ".")
					tokenKey = strings.ReplaceAll(tokenKey, " ", "_")

					switch styleType {
					case "FILL":
						colors := tokens["colors"].(map[string]any)
						colors[tokenKey] = map[string]any{
							"name":    name,
							"node_id": nodeID,
							"type":    "color",
						}
					case "TEXT":
						typography := tokens["typography"].(map[string]any)
						typography[tokenKey] = map[string]any{
							"name":    name,
							"node_id": nodeID,
							"type":    "typography",
						}
					case "EFFECT", "GRID":
						spacing := tokens["spacing"].(map[string]any)
						spacing[tokenKey] = map[string]any{
							"name":    name,
							"node_id": nodeID,
							"type":    styleType,
						}
					}
				}
			}
		}

		formatted, _ := json.MarshalIndent(tokens, "", "  ")
		tokensJSON := string(formatted)

		if outputPath != "" {
			if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
				return helpers.ErrorResult("write_error", fmt.Sprintf("failed to write tokens to %s: %v", outputPath, err)), nil
			}
			return helpers.TextResult(fmt.Sprintf("Design tokens written to %s\n\n%s", outputPath, tokensJSON)), nil
		}

		return helpers.TextResult(tokensJSON), nil
	}
}
