import { useCallback, useMemo, useState } from 'react'
import {
  useOrganizationMembers,
  useOrganizationMemberActions,
  useOrganizationRoles,
} from '@/features/organizations'
import { useAuthMethods } from '@/features/auth'
import { Loading } from '@/components/shared/loading'
import { ProviderIcon } from '@/features/auth/components/provider-icon'
import { useTranslation } from 'react-i18next'
import { useParams } from 'react-router'
import { Crown, Shield, Brush, Eye, Mail, Wallet } from 'lucide-react'
import { MemberTable } from '@/features/organizations/components/members/member-table'
import { InviteMemberDialog } from '@/features/organizations/components/members/invite-member-dialog'
import { MemberRoleInfo } from '@/features/organizations/components/members/member-role-select'
import { AppPagination } from '@/components/shared/app-pagination'
import { SearchInput } from '@/components/shared/search-input'
import { SearchEmptyState } from '@/components/shared/search-empty-state'

const roleColorMap = {
  'organization:viewer': {
    iconColor: 'text-gray-500',
    badgeBg: 'bg-gray-100 dark:bg-gray-800',
    badgeText: 'text-gray-800 dark:text-gray-200',
  },
  'organization:editor': {
    iconColor: 'text-blue-500',
    badgeBg: 'bg-blue-100 dark:bg-blue-900',
    badgeText: 'text-blue-800 dark:text-blue-200',
  },
  'organization:admin': {
    iconColor: 'text-emerald-500',
    badgeBg: 'bg-emerald-100 dark:bg-emerald-900',
    badgeText: 'text-emerald-800 dark:text-emerald-200',
  },
  'organization:owner': {
    iconColor: 'text-amber-500',
    badgeBg: 'bg-amber-100 dark:bg-amber-900',
    badgeText: 'text-amber-800 dark:text-amber-200',
  },
}

const roleIcons = {
  'organization:viewer': (
    <Eye size={14} className={roleColorMap['organization:viewer'].iconColor} />
  ),
  'organization:editor': (
    <Brush
      size={14}
      className={roleColorMap['organization:editor'].iconColor}
    />
  ),
  'organization:admin': (
    <Shield
      size={14}
      className={roleColorMap['organization:admin'].iconColor}
    />
  ),
  'organization:owner': (
    <Crown size={14} className={roleColorMap['organization:owner'].iconColor} />
  ),
}

function MembersPage() {
  const { t } = useTranslation()
  const { organizationId } = useParams()
  const [page, setPage] = useState(1)
  const [query, setQuery] = useState('')
  const limit = 10
  const offset = (page - 1) * limit

  const handleSearch = useCallback((value: string) => {
    setQuery(value)
    setPage(1)
  }, [])

  const {
    members,
    pagination,
    isLoading: membersLoading,
  } = useOrganizationMembers(organizationId, { offset, limit, query })
  const { roles } = useOrganizationRoles()
  const { methods: authMethods, isLoading: providersLoading } = useAuthMethods()
  const providers = useMemo(
    () => [...(authMethods.crypto ?? []), ...(authMethods.oauth ?? [])],
    [authMethods],
  )
  const { invite, updateRole, updateInfo, remove } =
    useOrganizationMemberActions(organizationId)

  const [openPopoverId, setOpenPopoverId] = useState(null)
  const [alias, setAlias] = useState<Record<string, string>>({})
  const [, setAliasCache] = useState<string | undefined>()
  const [inviteProvider, setInviteProvider] = useState('wallet')
  const [inviteIdentifier, setInviteIdentifier] = useState('')
  const [inviteError, setInviteError] = useState('')
  const [inviteLoading, setInviteLoading] = useState(false)

  const getProviderConfig = (providerId: string) => {
    const configs: Record<string, any> = {
      wallet: {
        id: 'wallet',
        name: t('ui.text.walletAddress'),
        icon: <Wallet className="w-4 h-4" />,
        placeholder: t('ui.text.walletPlaceholder'),
        description: t('ui.text.walletDescription'),
        validation: (value: string) =>
          value.startsWith('0x') && value.length >= 42,
        errorMessage: t('ui.text.invalidWallet'),
      },
      email: {
        id: 'email',
        name: t('ui.text.emailAddress'),
        icon: <Mail className="w-4 h-4" />,
        placeholder: t('ui.text.emailPlaceholder'),
        description: t('ui.text.emailDescription'),
        validation: (value: string) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value),
        errorMessage: t('ui.text.invalidEmail'),
      },
      google: {
        id: 'google',
        name: t('ui.text.googleAccount'),
        icon: <ProviderIcon provider="google" className="w-4 h-4" />,
        placeholder: t('ui.text.googlePlaceholder'),
        description: t('ui.text.googleDescription'),
        validation: (value: string) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value),
        errorMessage: t('ui.text.invalidGoogle'),
      },
      github: {
        id: 'github',
        name: t('ui.text.githubAccount'),
        icon: <ProviderIcon provider="github" className="w-4 h-4" />,
        placeholder: t('ui.text.githubPlaceholder'),
        description: t('ui.text.githubDescription'),
        validation: (value: string) =>
          /^[a-zA-Z0-9-]+$/.test(value) && value.length >= 1,
        errorMessage: t('ui.text.invalidGithub'),
      },
      metamask: {
        id: 'metamask',
        name: t('ui.text.metamaskAccount'),
        icon: <ProviderIcon provider="metamask" className="w-4 h-4" />,
        placeholder: t('ui.text.metamaskPlaceholder'),
        description: t('ui.text.metamaskDescription'),
        validation: (value: string) =>
          value.startsWith('0x') && value.length >= 42,
        errorMessage: t('ui.text.invalidWallet'),
      },
    }
    return configs[providerId]
  }

  const availableProviders = useMemo(() => {
    return providers.map((providerId: string) => {
      const config = getProviderConfig(providerId)
      return (
        config || {
          id: providerId,
          name: providerId.charAt(0).toUpperCase() + providerId.slice(1),
          icon: <ProviderIcon provider={providerId} className="w-4 h-4" />,
          placeholder: `Enter ${providerId}`,
          description: `Invite by ${providerId}`,
          validation: (value: string) => value.trim().length > 0,
          errorMessage: `Please enter a valid ${providerId}`,
        }
      )
    })
  }, [providers, t])

  const handleUpdateRole = async (uid: string, role: string) => {
    await updateRole(uid, { role })
    setOpenPopoverId(null)
  }

  const validateInviteInput = () => {
    const selectedProvider = availableProviders.find(
      (p: any) => p.id === inviteProvider,
    )
    if (!selectedProvider) {
      setInviteError(t('ui.text.invalidProvider'))
      return false
    }
    if (!inviteIdentifier.trim()) {
      setInviteError(t('ui.text.enterIdentifier'))
      return false
    }
    if (!selectedProvider.validation(inviteIdentifier.trim())) {
      setInviteError(selectedProvider.errorMessage)
      return false
    }
    setInviteError('')
    return true
  }

  const handleInvite = async () => {
    if (!validateInviteInput()) return
    setInviteLoading(true)
    try {
      await invite({
        identifier: inviteIdentifier.trim(),
        alias: '',
      })
      setInviteIdentifier('')
      setInviteProvider(availableProviders[0]?.id || 'wallet')
      setInviteError('')
    } catch (error) {
      console.error('Error inviting member:', error)
      setInviteError(t('ui.text.inviteFailed'))
    } finally {
      setInviteLoading(false)
    }
  }

  const handleProviderChange = (provider: string) => {
    setInviteProvider(provider)
    setInviteIdentifier('')
    setInviteError('')
  }

  const handleIdentifierChange = (value: string) => {
    setInviteIdentifier(value)
    setInviteError('')
  }

  const updateAlias = async (u: any) => {
    const a = alias[u.id]
    if (a === undefined) return
    const trimmed = a.trim()
    await updateInfo(u.id, { alias: trimmed })
    setAlias((prev) => {
      const next = { ...prev }
      delete next[u.id]
      return next
    })
    setAliasCache(undefined)
  }

  const handleRemoveMember = async (uid: string) => {
    await remove(uid)
  }

  const isInitialLoading = membersLoading && pagination === undefined

  return (
    <div className="w-full py-4 px-6">
      <div className="flex items-center gap-2 justify-between">
        <div className="flex flex-row items-center gap-2">
          <h2 className="text-2xl font-bold">
            {t('ui.text.organizationMembers')}
          </h2>
          <MemberRoleInfo />
        </div>
        <div className="flex p-0">
          <InviteMemberDialog
            availableProviders={availableProviders}
            providersLoading={providersLoading}
            inviteProvider={inviteProvider}
            inviteIdentifier={inviteIdentifier}
            inviteError={inviteError}
            inviteLoading={inviteLoading}
            onProviderChange={handleProviderChange}
            onIdentifierChange={handleIdentifierChange}
            onInvite={handleInvite}
          />
        </div>
      </div>
      <div className="mt-4 mb-3 max-w-sm">
        <SearchInput
          placeholder={t('ui.text.searchMembers')}
          onSearch={handleSearch}
        />
      </div>
      {isInitialLoading ? (
        <Loading />
      ) : members.length === 0 && query ? (
        <SearchEmptyState />
      ) : (
        <>
          <MemberTable
            members={members}
            roles={roles}
            roleColorMap={roleColorMap}
            roleIcons={roleIcons}
            openPopoverId={openPopoverId}
            alias={alias}
            onOpenPopoverChange={setOpenPopoverId}
            onAliasChange={(memberId: string, value: string) =>
              setAlias({ ...alias, [memberId]: value })
            }
            onAliasFocus={(value: string) => setAliasCache(value)}
            onAliasBlur={(originalAlias: string) =>
              setAliasCache(originalAlias)
            }
            onUpdateAlias={updateAlias}
            onUpdateRole={handleUpdateRole}
            onRemoveMember={handleRemoveMember}
          />
          {pagination && (
            <AppPagination
              pagination={{
                total: pagination.total,
                page,
                limit,
                offset,
                hasPrev: pagination.hasPrev,
                hasNext: pagination.hasNext,
              }}
              onPageChange={setPage}
            />
          )}
        </>
      )}
    </div>
  )
}

export default MembersPage
