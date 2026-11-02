package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine"
)

func main() {
	nodeID := flag.String("node", "", "node id to run, e.g. deepl_translate; empty lists every node")
	input := flag.String("input", "{}", "input object as JSON, or @path to a JSON file")
	credentials := flag.String("credentials", "{}", "credentials keyed by credential id as JSON, or @path to a JSON file")
	fileGatewayURL := flag.String("file-gateway-url", os.Getenv("FILE_GATEWAY_BASE_URL"), "file gateway base url, only for nodes that upload or download files")
	fileGatewayKey := flag.String("file-gateway-service-key", os.Getenv("FILE_GATEWAY_AUTH_SERVICE_KEY"), "file gateway service key")
	timeout := flag.Duration("timeout", 5*time.Minute, "execution timeout")
	flag.Parse()

	module := nodeengine.NewModule(nodeengine.Dependencies{
		FileGateway: filegateway.NewFileGateway(*fileGatewayURL, *fileGatewayKey),
	})

	if *nodeID == "" {
		for _, node := range module.NodeService.GetAllNodes() {
			_, _ = os.Stdout.WriteString(node.GetID() + "\t" + node.GetName() + "\n")
		}
		return
	}

	executor, ok := module.ExecutorService.GetByID(*nodeID)
	if !ok {
		slog.Error("unknown node", "node", *nodeID)
		os.Exit(1)
	}

	inputObject, err := readObject(*input)
	if err != nil {
		slog.Error("invalid input", "error", err)
		os.Exit(1)
	}

	credentialObjects, err := readObject(*credentials)
	if err != nil {
		slog.Error("invalid credentials", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	startedAt := time.Now()
	outputs, err := executor.ExecuteWithContext(ctx, credentialObjects, []map[string]any{inputObject})
	duration := time.Since(startedAt)
	if err != nil {
		slog.Error("node failed", "node", *nodeID, "duration", duration, "error", err)
		os.Exit(1)
	}

	raw, err := json.Marshal(outputs)
	if err != nil {
		slog.Error("failed to encode outputs", "error", err)
		os.Exit(1)
	}

	_, _ = os.Stdout.Write(raw)
	_, _ = os.Stdout.WriteString("\n")
	slog.Info("node finished", "node", *nodeID, "items", len(outputs), "duration", duration)
}

func readObject(value string) (map[string]any, error) {
	raw := []byte(strings.TrimSpace(value))
	if after, ok := strings.CutPrefix(string(raw), "@"); ok {
		content, err := os.ReadFile(after)
		if err != nil {
			return nil, err
		}
		raw = content
	}

	if len(raw) == 0 {
		return map[string]any{}, nil
	}

	object := map[string]any{}
	if err := json.Unmarshal(raw, &object); err != nil {
		return nil, err
	}

	return object, nil
}
