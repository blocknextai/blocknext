import { useState } from 'react'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { useTranslation } from 'react-i18next'
import { Type, Sparkles } from 'lucide-react'
import { CopyField } from '@/components/shared/copy-field'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import organizationsService from '@/features/organizations/services/organizations'

const SettingsBasicInfo = ({
  organization,
  setOrganization,
  setFormData,
  errors,
  setErrors,
  organizations,
  setOrganizations,
}) => {
  const { t } = useTranslation()
  const [isSubmitting, setIsSubmitting] = useState(false)

  const validateForm = () => {
    const newErrors = {}

    if (!organization.title.trim()) {
      newErrors.name = t('ui.text.orgNameRequired')
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSave = async () => {
    if (!validateForm()) {
      return
    }

    setIsSubmitting(true)

    try {
      const updateData = {
        title: organization.title,
        description: organization.description || '',
      }

      await organizationsService.update(organization.id, updateData)

      const orgs = [...organizations]
      const orgIndex = orgs.findIndex((o) => o.id === organization.id)
      if (orgIndex !== -1) {
        orgs[orgIndex] = {
          ...orgs[orgIndex],
          title: organization.title,
          id: organization.id,
        }
        setOrganizations(orgs)
      }
    } catch (error) {
      console.error('Update error:', error)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <>
      {/* Basic Information Card */}
      <Card className="border-0 shadow-lg bg-card/50 backdrop-blur-sm">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Type className="h-5 w-5 text-primary" />
            {t('ui.text.basicInformation')}
          </CardTitle>
          <CardDescription>{t('ui.text.provideOrgDetails')}</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="org-id" className="text-sm font-medium">
              {t('ui.text.organizationId')}
            </Label>
            <CopyField value={organization.id || ''} inputClassName="h-12" />
          </div>
          <div className="space-y-2">
            <Label htmlFor="name" className="text-sm font-medium">
              {t('ui.text.organizationName')}
            </Label>
            <Input
              id="name"
              type="text"
              placeholder={t('ui.text.enterOrgName')}
              value={organization.title}
              onChange={(e) => {
                setOrganization({
                  ...organization,
                  title: e.target.value,
                })
                setFormData((prev) => ({
                  ...prev,
                  name: e.target.value,
                }))
                if (errors.name) {
                  setErrors((prev) => ({
                    ...prev,
                    name: '',
                  }))
                }
              }}
              className="h-12"
            />
            {errors.name && (
              <p className="text-sm text-destructive">{errors.name}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="description" className="text-sm font-medium">
              {t('ui.text.description')}
            </Label>
            <Textarea
              id="description"
              placeholder={
                organization?.title
                  ? t('ui.text.briefDescriptionFor', {
                      name: organization.title,
                    })
                  : t('ui.text.summaryPlaceholder')
              }
              value={organization?.description || ''}
              onChange={(e) =>
                setOrganization({
                  ...organization,
                  description: e.target.value,
                })
              }
              className="min-h-[80px] resize-none"
            />
          </div>
        </CardContent>
      </Card>

      {/* Save Button */}
      <div className="flex justify-end">
        <Button type="button" disabled={isSubmitting} onClick={handleSave}>
          <Sparkles className="h-4 w-4 mr-2" />
          {isSubmitting
            ? t('ui.text.saving')
            : t('ui.text.saveGeneralSettings')}
        </Button>
      </div>
    </>
  )
}

SettingsBasicInfo.displayName = 'SettingsBasicInfo'

export { SettingsBasicInfo }
