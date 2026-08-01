package organizationusers

import (
	"context"

	"github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
	"github.com/google/uuid"
)

type OrganizationUserService interface {
	GetAllByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]*organizationusers.OrganizationUser, error)
	GetAllByIDs(ctx context.Context, ids []uuid.UUID) ([]*organizationusers.OrganizationUser, error)
	GetAllByOrganizationIDAndUserIDs(ctx context.Context, organizationID uuid.UUID, userIDs []uuid.UUID) ([]*organizationusers.OrganizationUser, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*organizationusers.OrganizationUser, error)
	GetByID(ctx context.Context, id uuid.UUID) (*organizationusers.OrganizationUser, error)
	GetByOrganizationIDAndUserID(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID) (*organizationusers.OrganizationUser, error)
	CountByOrganizationID(ctx context.Context, organizationID uuid.UUID) (int64, error)
}

type organizationUserService struct {
	organizationUserRepository organizationusers.OrganizationUserRepository
}

func NewOrganizationUserService(organizationUserRepository organizationusers.OrganizationUserRepository) OrganizationUserService {
	return &organizationUserService{organizationUserRepository: organizationUserRepository}
}

func (s *organizationUserService) GetAllByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]*organizationusers.OrganizationUser, error) {
	// TODO: hardcoded 1000 cap silently drops members in orgs with >1000 users.
	members, _, err := s.organizationUserRepository.GetAllByOrganizationID(ctx, organizationID, "", 0, 1000)
	return members, err
}

func (s *organizationUserService) GetAllByIDs(ctx context.Context, ids []uuid.UUID) ([]*organizationusers.OrganizationUser, error) {
	return s.organizationUserRepository.GetAllByIDs(ctx, ids)
}

func (s *organizationUserService) GetAllByOrganizationIDAndUserIDs(ctx context.Context, organizationID uuid.UUID, userIDs []uuid.UUID) ([]*organizationusers.OrganizationUser, error) {
	return s.organizationUserRepository.GetAllByOrganizationIDAndUserIDs(ctx, organizationID, userIDs)
}

func (s *organizationUserService) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*organizationusers.OrganizationUser, error) {
	return s.organizationUserRepository.GetAllByUserID(ctx, userID)
}

func (s *organizationUserService) GetByID(ctx context.Context, id uuid.UUID) (*organizationusers.OrganizationUser, error) {
	return s.organizationUserRepository.GetByID(ctx, id)
}

func (s *organizationUserService) GetByOrganizationIDAndUserID(ctx context.Context, organizationID, userID uuid.UUID) (*organizationusers.OrganizationUser, error) {
	return s.organizationUserRepository.GetByOrganizationIDAndUserID(ctx, organizationID, userID)
}

func (s *organizationUserService) CountByOrganizationID(ctx context.Context, organizationID uuid.UUID) (int64, error) {
	return s.organizationUserRepository.CountByOrganizationID(ctx, organizationID)
}
