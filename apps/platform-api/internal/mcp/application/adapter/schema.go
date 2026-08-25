package adapter

import (
	"maps"

	nodeEngineDomainNodes "github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func mcpAnnotations(a nodeEngineDomainNodes.NodeAnnotations) *mcpsdk.ToolAnnotations {
	return &mcpsdk.ToolAnnotations{
		ReadOnlyHint:    a.ReadOnly,
		DestructiveHint: a.Destructive,
		IdempotentHint:  a.Idempotent,
		OpenWorldHint:   a.OpenWorld,
	}
}

func augmentInputSchema(original *gjs.Schema, credentialKeys []string) *gjs.Schema {
	if len(credentialKeys) == 0 {
		return original
	}

	base := &gjs.Schema{Type: "object", Properties: map[string]*gjs.Schema{}}
	if original != nil {
		base = new(*original)
	}

	properties := make(map[string]*gjs.Schema, len(base.Properties)+1)
	maps.Copy(properties, base.Properties)

	credentialProps := make(map[string]*gjs.Schema, len(credentialKeys))
	for _, key := range credentialKeys {
		credentialProps[key] = &gjs.Schema{
			Type:        "string",
			Description: "Reference to a saved credential. Format: credential:organization:<uuid>.",
			Pattern:     "^credential:organization:[0-9a-fA-F-]+$",
		}
	}

	properties[credentialsKey] = &gjs.Schema{
		Type:        "object",
		Title:       "Credentials",
		Description: "References to saved credentials required by this tool.",
		Properties:  credentialProps,
		Required:    append([]string{}, credentialKeys...),
	}

	required := append([]string{}, base.Required...)
	required = append(required, credentialsKey)

	base.Properties = properties
	base.Required = required
	return base
}

func wrapOutputSchema(original *gjs.Schema) *gjs.Schema {
	if original == nil {
		return nil
	}
	return &gjs.Schema{
		Type: "object",
		Properties: map[string]*gjs.Schema{
			itemsKey: original,
		},
		Required: []string{
			itemsKey,
		},
	}
}
