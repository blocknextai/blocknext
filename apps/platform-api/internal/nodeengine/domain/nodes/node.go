package nodes

import (
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type Node struct {
	ID                   string          `json:"id,omitempty"`
	Kind                 NodeKind        `json:"kind"`
	Version              string          `json:"version,omitempty"`
	Name                 string          `json:"name,omitempty"`
	Description          string          `json:"description,omitempty"`
	Icon                 NodeIcon        `json:"icon"`
	Inputs               []NodeHandle    `json:"inputs"`
	Outputs              []NodeHandle    `json:"outputs"`
	InputSchema          *gjs.Schema     `json:"inputSchema,omitempty"`
	OutputSchema         *gjs.Schema     `json:"outputSchema,omitempty"`
	Categories           []string        `json:"categories,omitempty"`
	SubCategories        []string        `json:"subCategories,omitempty"`
	Tags                 []string        `json:"tags,omitempty"`
	SupportedCredentials []string        `json:"supportedCredentials,omitempty"`
	Annotations          NodeAnnotations `json:"annotations"`
	Disabled             bool            `json:"disabled,omitempty"`
	HasNaturalLanguage   bool            `json:"hasNaturalLanguage"`
}

type NodeManager interface {
	GetID() string
	GetKind() NodeKind
	GetVersion() string
	GetName() string
	GetDescription() string
	GetIcon() NodeIcon
	GetInputs() []NodeHandle
	SetInputs(handles []NodeHandle)
	GetOutputs() []NodeHandle
	SetOutputs(handles []NodeHandle)
	GetInputSchema() *gjs.Schema
	GetOutputSchema() *gjs.Schema
	GetCategories() []string
	GetSubCategories() []string
	GetTags() []string
	GetSupportedCredentials() []string
	GetAnnotations() NodeAnnotations
	GetDisabled() bool
	GetHasNaturalLanguage() bool
}

func (n *Node) GetID() string {
	return n.ID
}

func (n *Node) GetKind() NodeKind {
	return n.Kind
}

func (n *Node) GetVersion() string {
	return n.Version
}

func (n *Node) GetName() string {
	return n.Name
}

func (n *Node) GetDescription() string {
	return n.Description
}

func (n *Node) GetIcon() NodeIcon {
	return n.Icon
}

func (n *Node) GetInputs() []NodeHandle {
	return n.Inputs
}

func (n *Node) SetInputs(handles []NodeHandle) {
	n.Inputs = handles
}

func (n *Node) GetOutputs() []NodeHandle {
	return n.Outputs
}

func (n *Node) SetOutputs(handles []NodeHandle) {
	n.Outputs = handles
}

func (n *Node) GetInputSchema() *gjs.Schema {
	return n.InputSchema
}

func (n *Node) GetOutputSchema() *gjs.Schema {
	return n.OutputSchema
}

func (n *Node) GetCategories() []string {
	return n.Categories
}

func (n *Node) GetSubCategories() []string {
	return n.SubCategories
}

func (n *Node) GetTags() []string {
	return n.Tags
}

func (n *Node) GetSupportedCredentials() []string {
	return n.SupportedCredentials
}

func (n *Node) GetAnnotations() NodeAnnotations {
	return n.Annotations
}

func (n *Node) GetDisabled() bool {
	return n.Disabled
}

func (n *Node) GetHasNaturalLanguage() bool {
	return n.HasNaturalLanguage
}
