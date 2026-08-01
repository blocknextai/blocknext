package credentials

import (
	"context"

	commonDomainOAuth2 "github.com/blocknextai/platform-api/internal/common/domain/oauth2"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type Credential struct {
	ID                string         `json:"id,omitempty"`
	PlatformID        string         `json:"-"`
	Name              string         `json:"name,omitempty"`
	Description       string         `json:"description,omitempty"`
	Icon              CredentialIcon `json:"icon"`
	IsOAuth1          bool           `json:"isOAuth1,omitempty"`
	IsOAuth2          bool           `json:"isOAuth2,omitempty"`
	IsSupportPlatform bool           `json:"isSupportPlatform,omitempty"`
	Schema            *gjs.Schema    `json:"schema,omitempty"`
	SupportedNodes    *[]string      `json:"supportedNodes,omitempty"`
	Disabled          bool           `json:"disabled,omitempty"`
}

type CredentialManager interface {
	GetID() string
	GetPlatformID() string
	GetName() string
	GetDescription() string
	GetIcon() CredentialIcon
	GetIsOAuth1() bool
	SetIsOAuth1(value bool)
	GetIsOAuth2() bool
	SetIsOAuth2(value bool)
	GetIsSupportPlatform() bool
	SetIsSupportPlatform(value bool)
	GetSchema() *gjs.Schema
	GetSupportedNodes() *[]string
	GetDisabled() bool
}

type RefreshableCredential interface {
	RefreshToken(ctx context.Context, credential commonDomainOAuth2.Credential) (*commonDomainOAuth2.Token, error)
}

func (c *Credential) GetID() string {
	return c.ID
}

func (c *Credential) GetPlatformID() string {
	return c.PlatformID
}

func (c *Credential) GetName() string {
	return c.Name
}

func (c *Credential) GetDescription() string {
	return c.Description
}

func (c *Credential) GetIcon() CredentialIcon {
	return c.Icon
}

func (c *Credential) GetIsOAuth1() bool {
	return c.IsOAuth1
}

func (c *Credential) SetIsOAuth1(value bool) {
	c.IsOAuth1 = value
}

func (c *Credential) GetIsOAuth2() bool {
	return c.IsOAuth2
}

func (c *Credential) SetIsOAuth2(value bool) {
	c.IsOAuth2 = value
}

func (c *Credential) GetIsSupportPlatform() bool {
	return c.IsSupportPlatform
}

func (c *Credential) SetIsSupportPlatform(value bool) {
	c.IsSupportPlatform = value
}

func (c *Credential) GetSchema() *gjs.Schema {
	return c.Schema
}

func (c *Credential) GetSupportedNodes() *[]string {
	return c.SupportedNodes
}

func (c *Credential) GetDisabled() bool {
	return c.Disabled
}
