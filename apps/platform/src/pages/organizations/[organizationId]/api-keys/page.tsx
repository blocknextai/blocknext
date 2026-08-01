import { useParams } from 'react-router'
import { useApiKeyScopes, useOrganizationApiKeys } from '@/features/api-keys'
import APIKeyTab from '@/features/api-keys/components/api-key-tab'

function OrganizationAPIKeysPage() {
  const { organizationId } = useParams()

  const {
    loading,
    apiKeys,
    pagination,
    fetchKeys,
    handleCreate,
    handleRegenerate,
    handleDelete,
  } = useOrganizationApiKeys(organizationId)

  const { scopes, loadingScopes } = useApiKeyScopes()

  return (
    <APIKeyTab
      loading={loading}
      apiKeys={apiKeys}
      pagination={pagination}
      scopes={scopes}
      loadingScopes={loadingScopes}
      onFetch={fetchKeys}
      onCreate={handleCreate}
      onRegenerate={handleRegenerate}
      onDelete={handleDelete}
    />
  )
}

export default OrganizationAPIKeysPage
