package organizations

import (
	"context"

	organizationsDomainOrganizations "github.com/blocknextai/platform-api/internal/organizations/domain/organizations"
	"github.com/google/uuid"
)

type OrganizationService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*organizationsDomainOrganizations.Organization, error)
	GetAllByIDs(ctx context.Context, ids []uuid.UUID) ([]*organizationsDomainOrganizations.Organization, error)
}

type organizationService struct {
	organizationRepository organizationsDomainOrganizations.OrganizationRepository
}

func NewOrganizationService(
	organizationRepository organizationsDomainOrganizations.OrganizationRepository,
) OrganizationService {
	return &organizationService{
		organizationRepository: organizationRepository,
	}
}

func (s *organizationService) GetByID(ctx context.Context, id uuid.UUID) (*organizationsDomainOrganizations.Organization, error) {
	return s.organizationRepository.GetByID(ctx, id)
}

func (s *organizationService) GetAllByIDs(ctx context.Context, ids []uuid.UUID) ([]*organizationsDomainOrganizations.Organization, error) {
	return s.organizationRepository.GetAllByIDs(ctx, ids)
}
