import { useTranslation } from 'react-i18next'
import { useNavigate, useSearchParams } from 'react-router'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { PageLoading } from '@/components/shared/loading'
import { PlatformLogo } from '@/features/navigation/components/logo'
import { ModeToggle } from '@/features/theme/components/mode-toggle'
import {
  McpConsentCard,
  useMcpAuthorizationConsent,
} from '@/features/mcp-oauth'

const OAuthMcpConsentPage = () => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const consent = useMcpAuthorizationConsent(searchParams.get('requestId'))

  if (consent.isLoading) {
    return <PageLoading />
  }

  return (
    <div className="flex min-h-app-screen flex-col">
      <header className="flex h-16 shrink-0 items-center justify-between px-4">
        <PlatformLogo />
        <ModeToggle aria-label={t('ui.text.chooseTheme')} />
      </header>

      <div className="flex flex-1 items-center justify-center p-4">
        {consent.isUnavailable || !consent.authorizationRequest ? (
          <Card className="w-full max-w-md">
            <CardHeader>
              <CardTitle>{t('ui.text.mcpConsentUnavailableTitle')}</CardTitle>
              <CardDescription>
                {t('ui.text.mcpConsentUnavailableDescription')}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Button
                onClick={() => navigate('/organizations', { replace: true })}
              >
                {t('ui.text.backToHome')}
              </Button>
            </CardContent>
          </Card>
        ) : (
          <McpConsentCard
            authorizationRequest={consent.authorizationRequest}
            organizations={consent.organizations}
            organizationId={consent.organizationId}
            isSubmitting={consent.isSubmitting}
            onOrganizationChange={consent.setOrganizationId}
            onApprove={consent.approve}
            onDeny={consent.deny}
          />
        )}
      </div>
    </div>
  )
}

export default OAuthMcpConsentPage
