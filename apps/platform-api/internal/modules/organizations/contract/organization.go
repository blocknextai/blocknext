package contract

import (
	organizationsApplicationOrganizations "github.com/blocknextai/platform-api/internal/modules/organizations/internal/application/organizations"
	organizationsDomainOrganizations "github.com/blocknextai/platform-api/internal/modules/organizations/internal/domain/organizations"
)

var ErrOrganizationNotFound = organizationsDomainOrganizations.ErrOrganizationNotFound

type Organization = organizationsDomainOrganizations.Organization

type OrganizationService = organizationsApplicationOrganizations.OrganizationService
