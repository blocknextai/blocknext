import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Building } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useOrganizations } from '@/features/organizations'
import { Loading } from '@/components/shared/loading'
import { OrganizationCard } from '@/features/organizations/components/organization-card'

const OrganizationsPage = () => {
  const { t } = useTranslation()
  const { organizations, isLoading } = useOrganizations()

  if (isLoading) {
    return <Loading />
  }

  return (
    <div className="w-full h-full p-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold">{t('ui.text.myOrganizations')}</h1>
      </div>

      {organizations.length === 0 ? (
        <div className="flex flex-col items-center justify-center h-full text-center -mt-10">
          <Building className="size-10 text-muted-foreground/80 mb-4" />
          <h3 className="text-lg font-semibold mb-2">
            {t('ui.text.noOrganizationsYet')}
          </h3>
          <p className="text-muted-foreground mb-4 max-w-xs">
            {t('ui.text.createFirstOrgDescription')}
          </p>
          <Button asChild>
            <Link to="/organizations/new">
              {t('ui.text.createOrganization')}
            </Link>
          </Button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {organizations.map((organization) => (
            <OrganizationCard
              key={organization.id}
              organization={organization}
            />
          ))}
        </div>
      )}
    </div>
  )
}

export default OrganizationsPage
