package credentials

import (
	"strings"
	"time"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type Credential struct {
	database.BaseEntity

	OwnerType  commonDomain.OwnerType
	OwnerID    uuid.UUID
	SourceType SourceType

	Key   string
	Title string
	Data  string
}

func New(
	ownerType commonDomain.OwnerType,
	ownerId uuid.UUID,
	sourceType SourceType,
	key string,
	title string,
	data string,
) (*Credential, error) {
	utcNow := time.Now().UTC()

	credential := &Credential{
		ID:         bnuuid.NewV7(),
		CreatedAt:  utcNow,
		UpdatedAt:  utcNow,
		DeletedAt:  nil,
		OwnerType:  ownerType,
		OwnerID:    ownerId,
		SourceType: sourceType,
		Key:        key,
		Title:      title,
		Data:       data,
	}

	return credential.validateThenReturn()
}

func (c *Credential) Update(
	key string,
	title string,
	data string,
) (*Credential, error) {
	c.UpdatedAt = time.Now().UTC()

	c.Key = key
	c.Title = title
	c.Data = data

	return c.validateThenReturn()
}

func (c *Credential) Delete() (*Credential, error) {
	utcNow := time.Now().UTC()

	c.UpdatedAt = utcNow
	c.DeletedAt = new(utcNow)
	return c.validateThenReturn()
}

func (c *Credential) validateThenReturn() (*Credential, error) {
	if !c.OwnerType.IsValid() {
		return nil, ErrOwnerTypeIsInvalid
	}

	if c.OwnerID == uuid.Nil {
		return nil, ErrOwnerIDIsRequired
	}

	if strings.TrimSpace(c.Title) == "" {
		return nil, ErrTitleIsRequired
	}

	if strings.TrimSpace(c.Key) == "" {
		return nil, ErrKeyIsRequired
	}

	if strings.TrimSpace(c.Data) == "" {
		return nil, ErrDataIsRequired
	}

	return c, nil
}
