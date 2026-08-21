import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router'
import { Building } from 'lucide-react'
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
import { CreateOrganizationForm } from '@/features/organizations/components/create-organization-form'
import { useOrganizations } from '@/features/organizations'
import { useOrganizationStore } from '@/stores/organization'
import { logoutAndRedirect } from '@/features/auth'

const OrganizationCreatePage = () => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { organizations, isLoading } = useOrganizations()
  const setOrganizationId = useOrganizationStore((s) => s.setOrganizationId)

  const hasOrganizations = organizations.length > 0

  const handleCreated = (organization) => {
    setOrganizationId(organization.id)
    navigate(`/organizations/${organization.id}`, { replace: true })
  }

  if (isLoading) {
    return <PageLoading />
  }

  return (
    <div className="flex min-h-app-screen flex-col">
      <header className="flex h-16 shrink-0 items-center justify-between px-4">
        <PlatformLogo />
        <ModeToggle aria-label={t('ui.text.chooseTheme')} />
      </header>

      <div className="flex flex-1 items-center justify-center p-4">
        <Card className="w-full max-w-md">
          <CardHeader>
            <div className="mb-2 flex size-10 items-center justify-center rounded-lg bg-muted">
              <Building className="size-5 text-muted-foreground" />
            </div>
            <CardTitle>
              {hasOrganizations
                ? t('ui.text.createOrganization')
                : t('ui.text.createYourFirstOrganization')}
            </CardTitle>
            <CardDescription>
              {t('ui.text.createFirstOrgDescription')}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <CreateOrganizationForm
              autoFocus
              onCreated={handleCreated}
              onCancel={
                hasOrganizations ? () => navigate('/organizations') : undefined
              }
            />
            {!hasOrganizations && (
              <div className="mt-6 text-center">
                <Button
                  variant="link"
                  className="text-muted-foreground"
                  onClick={logoutAndRedirect}
                >
                  {t('ui.text.logOut')}
                </Button>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export default OrganizationCreatePage
