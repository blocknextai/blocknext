import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { organizationBasicInfoSchema } from '@/features/organizations/schemas'
import { useOrganizationActions } from '@/features/organizations'
import toast from '@/lib/toast'

const CreateOrganizationForm = ({
  onSubmit,
  onCreated,
  onCancel,
  submitLabel,
  autoFocus = false,
}: {
  onSubmit?: (values: { title: string; description?: string }) => Promise<any>
  onCreated?: (organization: any) => void
  onCancel?: () => void
  submitLabel?: string
  autoFocus?: boolean
}) => {
  const { t } = useTranslation()
  const { create } = useOrganizationActions()
  const [values, setValues] = useState({ title: '', description: '' })
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [isSubmitting, setIsSubmitting] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    const result = organizationBasicInfoSchema.safeParse(values)
    if (!result.success) {
      const fieldErrors: Record<string, string> = {}
      for (const issue of result.error.issues) {
        const field = String(issue.path[0])
        if (fieldErrors[field]) {
          continue
        }
        fieldErrors[field] = issue.message.startsWith('ui.')
          ? t(issue.message)
          : issue.message
      }
      setErrors(fieldErrors)
      return
    }

    setIsSubmitting(true)
    try {
      const organization = onSubmit
        ? await onSubmit(result.data)
        : await create(result.data)
      setValues({ title: '', description: '' })
      setErrors({})
      onCreated?.(organization)
    } catch (error) {
      console.error('Create organization error:', error)
      toast.error(t('ui.text.createOrganizationFailed'))
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="grid gap-4">
      <div className="grid gap-3">
        <Label htmlFor="organization-title">{t('ui.text.title')}</Label>
        <Input
          id="organization-title"
          name="title"
          autoFocus={autoFocus}
          placeholder={t('ui.text.myOrganization')}
          value={values.title}
          onChange={(e) => {
            setValues({ ...values, title: e.target.value })
            if (errors.title) {
              setErrors({ ...errors, title: '' })
            }
          }}
        />
        {errors.title && (
          <p className="text-sm text-destructive">{errors.title}</p>
        )}
      </div>

      <div className="grid gap-3">
        <Label htmlFor="organization-description">
          {t('ui.text.description')}
        </Label>
        <Textarea
          id="organization-description"
          name="description"
          placeholder={t('ui.text.descriptionOfOrganization')}
          value={values.description}
          onChange={(e) => {
            setValues({ ...values, description: e.target.value })
            if (errors.description) {
              setErrors({ ...errors, description: '' })
            }
          }}
        />
        {errors.description && (
          <p className="text-sm text-destructive">{errors.description}</p>
        )}
      </div>

      <div className="flex justify-end gap-2">
        {onCancel && (
          <Button type="button" variant="outline" onClick={onCancel}>
            {t('ui.text.cancel')}
          </Button>
        )}
        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting
            ? t('ui.text.creating')
            : (submitLabel ?? t('ui.text.create'))}
        </Button>
      </div>
    </form>
  )
}

CreateOrganizationForm.displayName = 'CreateOrganizationForm'

export { CreateOrganizationForm }
