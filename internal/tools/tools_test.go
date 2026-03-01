package tools

import (
	"context"
	"os"
	"testing"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// ---------- helpers ----------

func callTool(t *testing.T, handler func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error), args map[string]any) *pluginv1.ToolResponse {
	t.Helper()
	var s *structpb.Struct
	if args != nil {
		var err error
		s, err = structpb.NewStruct(args)
		if err != nil {
			t.Fatalf("NewStruct: %v", err)
		}
	}
	resp, err := handler(context.Background(), &pluginv1.ToolRequest{Arguments: s})
	if err != nil {
		t.Fatalf("handler returned Go error: %v", err)
	}
	return resp
}

func isError(resp *pluginv1.ToolResponse) bool {
	return resp != nil && !resp.Success
}

func errorCode(resp *pluginv1.ToolResponse) string {
	if resp == nil {
		return ""
	}
	return resp.GetErrorCode()
}

func figmaTokenSet() bool {
	return os.Getenv("FIGMA_ACCESS_TOKEN") != ""
}

// ---------- figma_get_file ----------

func TestFigmaGetFile_MissingFileKey(t *testing.T) {
	resp := callTool(t, FigmaGetFile(), map[string]any{})
	if !isError(resp) {
		t.Fatal("expected error response for missing file_key")
	}
	if code := errorCode(resp); code != "validation_error" {
		t.Errorf("expected code=validation_error, got %q", code)
	}
}

func TestFigmaGetFile_NoToken(t *testing.T) {
	if figmaTokenSet() {
		t.Skip("FIGMA_ACCESS_TOKEN is set — skipping no-token test")
	}
	// Ensure the token is definitely unset for this test.
	t.Setenv("FIGMA_ACCESS_TOKEN", "")

	resp := callTool(t, FigmaGetFile(), map[string]any{"file_key": "test123"})
	if !isError(resp) {
		t.Fatal("expected error response when FIGMA_ACCESS_TOKEN is not set")
	}
	if code := errorCode(resp); code != "figma_error" {
		t.Errorf("expected code=figma_error, got %q", code)
	}
}

// ---------- figma_get_components ----------

func TestFigmaGetComponents_MissingFileKey(t *testing.T) {
	resp := callTool(t, FigmaGetComponents(), map[string]any{})
	if !isError(resp) {
		t.Fatal("expected error response for missing file_key")
	}
	if code := errorCode(resp); code != "validation_error" {
		t.Errorf("expected code=validation_error, got %q", code)
	}
}

func TestFigmaGetComponents_NoToken(t *testing.T) {
	if figmaTokenSet() {
		t.Skip("FIGMA_ACCESS_TOKEN is set — skipping no-token test")
	}
	t.Setenv("FIGMA_ACCESS_TOKEN", "")

	resp := callTool(t, FigmaGetComponents(), map[string]any{"file_key": "test123"})
	if !isError(resp) {
		t.Fatal("expected error response when FIGMA_ACCESS_TOKEN is not set")
	}
	if code := errorCode(resp); code != "figma_error" {
		t.Errorf("expected code=figma_error, got %q", code)
	}
}

// ---------- figma_get_styles ----------

func TestFigmaGetStyles_MissingFileKey(t *testing.T) {
	resp := callTool(t, FigmaGetStyles(), map[string]any{})
	if !isError(resp) {
		t.Fatal("expected error response for missing file_key")
	}
	if code := errorCode(resp); code != "validation_error" {
		t.Errorf("expected code=validation_error, got %q", code)
	}
}

func TestFigmaGetStyles_NoToken(t *testing.T) {
	if figmaTokenSet() {
		t.Skip("FIGMA_ACCESS_TOKEN is set — skipping no-token test")
	}
	t.Setenv("FIGMA_ACCESS_TOKEN", "")

	resp := callTool(t, FigmaGetStyles(), map[string]any{"file_key": "test123"})
	if !isError(resp) {
		t.Fatal("expected error response when FIGMA_ACCESS_TOKEN is not set")
	}
	if code := errorCode(resp); code != "figma_error" {
		t.Errorf("expected code=figma_error, got %q", code)
	}
}

// ---------- figma_get_node ----------

func TestFigmaGetNode_MissingArgs(t *testing.T) {
	// Missing both file_key and node_id.
	resp := callTool(t, FigmaGetNode(), map[string]any{})
	if !isError(resp) {
		t.Fatal("expected error response for missing file_key and node_id")
	}
	if code := errorCode(resp); code != "validation_error" {
		t.Errorf("expected code=validation_error, got %q", code)
	}
}

func TestFigmaGetNode_MissingNodeID(t *testing.T) {
	// file_key present but node_id missing.
	resp := callTool(t, FigmaGetNode(), map[string]any{"file_key": "test123"})
	if !isError(resp) {
		t.Fatal("expected error response for missing node_id")
	}
	if code := errorCode(resp); code != "validation_error" {
		t.Errorf("expected code=validation_error, got %q", code)
	}
}

func TestFigmaGetNode_NoToken(t *testing.T) {
	if figmaTokenSet() {
		t.Skip("FIGMA_ACCESS_TOKEN is set — skipping no-token test")
	}
	t.Setenv("FIGMA_ACCESS_TOKEN", "")

	resp := callTool(t, FigmaGetNode(), map[string]any{
		"file_key": "test123",
		"node_id":  "1:2",
	})
	if !isError(resp) {
		t.Fatal("expected error response when FIGMA_ACCESS_TOKEN is not set")
	}
	if code := errorCode(resp); code != "figma_error" {
		t.Errorf("expected code=figma_error, got %q", code)
	}
}

// ---------- figma_export_node ----------

func TestFigmaExportNode_MissingArgs(t *testing.T) {
	// Missing both file_key and node_id.
	resp := callTool(t, FigmaExportNode(), map[string]any{})
	if !isError(resp) {
		t.Fatal("expected error response for missing file_key and node_id")
	}
	if code := errorCode(resp); code != "validation_error" {
		t.Errorf("expected code=validation_error, got %q", code)
	}
}

func TestFigmaExportNode_MissingNodeID(t *testing.T) {
	// file_key present but node_id missing.
	resp := callTool(t, FigmaExportNode(), map[string]any{"file_key": "test123"})
	if !isError(resp) {
		t.Fatal("expected error response for missing node_id")
	}
	if code := errorCode(resp); code != "validation_error" {
		t.Errorf("expected code=validation_error, got %q", code)
	}
}

func TestFigmaExportNode_NoToken(t *testing.T) {
	if figmaTokenSet() {
		t.Skip("FIGMA_ACCESS_TOKEN is set — skipping no-token test")
	}
	t.Setenv("FIGMA_ACCESS_TOKEN", "")

	resp := callTool(t, FigmaExportNode(), map[string]any{
		"file_key": "test123",
		"node_id":  "1:2",
	})
	if !isError(resp) {
		t.Fatal("expected error response when FIGMA_ACCESS_TOKEN is not set")
	}
	if code := errorCode(resp); code != "figma_error" {
		t.Errorf("expected code=figma_error, got %q", code)
	}
}

func TestFigmaExportNode_NoToken_WithFormat(t *testing.T) {
	if figmaTokenSet() {
		t.Skip("FIGMA_ACCESS_TOKEN is set — skipping no-token test")
	}
	t.Setenv("FIGMA_ACCESS_TOKEN", "")

	// Provide all args including optional format and scale.
	resp := callTool(t, FigmaExportNode(), map[string]any{
		"file_key": "test123",
		"node_id":  "1:2",
		"format":   "svg",
		"scale":    float64(2),
	})
	if !isError(resp) {
		t.Fatal("expected error response when FIGMA_ACCESS_TOKEN is not set")
	}
	if code := errorCode(resp); code != "figma_error" {
		t.Errorf("expected code=figma_error, got %q", code)
	}
}

// ---------- figma_sync_tokens ----------

func TestFigmaSyncTokens_MissingFileKey(t *testing.T) {
	resp := callTool(t, FigmaSyncTokens(), map[string]any{})
	if !isError(resp) {
		t.Fatal("expected error response for missing file_key")
	}
	if code := errorCode(resp); code != "validation_error" {
		t.Errorf("expected code=validation_error, got %q", code)
	}
}

func TestFigmaSyncTokens_NoToken(t *testing.T) {
	if figmaTokenSet() {
		t.Skip("FIGMA_ACCESS_TOKEN is set — skipping no-token test")
	}
	t.Setenv("FIGMA_ACCESS_TOKEN", "")

	resp := callTool(t, FigmaSyncTokens(), map[string]any{"file_key": "test123"})
	if !isError(resp) {
		t.Fatal("expected error response when FIGMA_ACCESS_TOKEN is not set")
	}
	if code := errorCode(resp); code != "figma_error" {
		t.Errorf("expected code=figma_error, got %q", code)
	}
}
