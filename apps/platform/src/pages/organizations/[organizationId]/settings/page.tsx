import { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { useOrganizationStore } from '@/stores/organization'
import { useNavigate, useParams } from 'react-router'
import { Loading } from '@/components/shared/loading'
import { organizationsService } from '@/features/organizations'
import { useTranslation } from 'react-i18next'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Building2, EllipsisVertical, Trash2 } from 'lucide-react'
import { SettingsBasicInfo } from '@/features/organizations/components/settings/settings-basic-info'
import { SettingsDangerZone } from '@/features/organizations/components/settings/settings-danger-zone'

function OrganizationSettingsPage() {
  const { t } = useTranslation()
  const { organizationId } = useParams()
  const organizations = useOrganizationStore((s) => s.organizations)
  const setOrganizations = useOrganizationStore((s) => s.setOrganizations)
  const [organization, setOrganization] = useState({
    id: '',
    title: '',
    description: '',
  })
  const [formData, setFormData] = useState({
    name: '',
  })

  const [errors, setErrors] = useState({})

  const [loading, setLoading] = useState(true)
  const [delOpen, setDelOpen] = useState(false)
  const navigate = useNavigate()

  const deleteOrganization = () => {
    organizationsService.delete(organization.id).then(() => {
      setTimeout(() => {
        navigate('/organizations')
      }, 250)
    })
  }

  useEffect(() => {
    if (organizationId) {
      setLoading(true)
      organizationsService
        .getById(organizationId)
        .then((orgResponse) => {
          const orgRes = orgResponse.data
          setOrganization({
            id: orgRes.id,
            title: orgRes.title,
            description: orgRes.description,
          })

          setFormData((prev) => ({
            ...prev,
            name: orgRes.title || '',
          }))

          setLoading(false)
        })
        .catch(() => {
          setLoading(false)
        })
    }
  }, [organizationId])

  if (loading) {
    return <Loading />
  }

  return (
    <div className="w-full h-full flex flex-col">
      <div className="min-h-screen bg-linear-to-br mt-4 from-background via-background to-muted/20">
        <div className="container mx-auto">
          <div className="max-w-4xl mx-auto">
            {/* Header */}
            <div className="mb-8">
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-3">
                  <div className="p-2 rounded-lg bg-primary/10">
                    <Building2 className="h-6 w-6 text-primary" />
                  </div>
                  <div>
                    <h1 className="text-3xl font-bold bg-linear-to-r from-foreground to-muted-foreground bg-clip-text text-transparent">
                      {t('ui.text.organizationSettings')}
                    </h1>
                  </div>
                </div>

                {/* Actions Dropdown */}
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="outline"
                      size="sm"
                      className="h-10 w-10 p-0"
                    >
                      <EllipsisVertical className="h-4 w-4" />
                      <span className="sr-only">{t('ui.text.openMenu')}</span>
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuLabel>
                      {t('ui.text.actions')}
                    </DropdownMenuLabel>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem
                      className="text-destructive focus:text-destructive cursor-pointer"
                      onClick={() => setDelOpen(true)}
                    >
                      <Trash2 className="mr-2 h-4 w-4" />
                      {t('ui.text.delete')}
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
              <p className="text-muted-foreground text-lg">
                {t('ui.text.manageOrgDetails')}
              </p>
            </div>

            <div className="space-y-6">
              <SettingsBasicInfo
                organization={organization}
                setOrganization={setOrganization}
                formData={formData}
                setFormData={setFormData}
                errors={errors}
                setErrors={setErrors}
                organizations={organizations}
                setOrganizations={setOrganizations}
              />
            </div>
          </div>
        </div>
      </div>

      <SettingsDangerZone
        delOpen={delOpen}
        setDelOpen={setDelOpen}
        onDelete={deleteOrganization}
      />
    </div>
  )
}

export default OrganizationSettingsPage
