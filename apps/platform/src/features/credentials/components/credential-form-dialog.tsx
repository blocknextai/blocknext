import { useEffect, useState } from 'react'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Checkbox } from '@/components/ui/checkbox'
import { Badge } from '@/components/ui/badge'
import { Building2, Eye, EyeOff, User } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { CopyField } from '@/components/shared/copy-field'
import { Loading } from '@/components/shared/loading'

const CredentialFormDialog = ({
  open,
  onOpenChange,
  selectedSecret,
  secretDetailLoading,
  usePlatformCredential,
  setUsePlatformCredential,
  updateSelectedSecret,
  saveCredentials,
  authorizeOAuth2,
  reauthorizeOAuth2,
}) => {
  const { t } = useTranslation()
  const [visibleFields, setVisibleFields] = useState<Record<string, boolean>>(
    {},
  )
  const [errors, setErrors] = useState<Record<string, boolean>>({})

  const toggleVisibility = (key: string) => {
    setVisibleFields((prev) => ({ ...prev, [key]: !prev[key] }))
  }

  const clearError = (key: string) => {
    setErrors((prev) => {
      if (!prev[key]) {
        return prev
      }
      const next = { ...prev }
      delete next[key]
      return next
    })
  }

  useEffect(() => {
    if (!open) {
      setErrors({})
    }
  }, [open])

  useEffect(() => {
    setErrors({})
  }, [usePlatformCredential])

  const getFieldValue = (schema: any) => {
    if (selectedSecret?.data && selectedSecret.data[schema.key]) {
      return selectedSecret.data[schema.key]
    }
    if (
      selectedSecret?.secretDetails?.data &&
      selectedSecret.secretDetails.data[schema.key]
    ) {
      return selectedSecret.secretDetails.data[schema.key]
    }
    if (schema?.default) {
      return schema.default
    }
    return ''
  }

  const validate = () => {
    const errs: Record<string, boolean> = {}
    const titleValue = (selectedSecret?.label || '').toString().trim()
    if (!titleValue) {
      errs['s-t'] = true
    }

    if (!usePlatformCredential && selectedSecret?.provider) {
      for (const schema of selectedSecret.provider) {
        if (schema.hidden || schema.readOnly || !schema.required) {
          continue
        }
        const value = (getFieldValue(schema) || '').toString().trim()
        if (!value) {
          errs[schema.key] = true
        }
      }
    }

    setErrors(errs)
    return Object.keys(errs).length === 0
  }

  const handleSave = () => {
    if (!validate()) {
      return
    }
    saveCredentials()
    onOpenChange(false)
  }

  const handleAuthorize = () => {
    if (!validate()) {
      return
    }
    authorizeOAuth2()
  }

  const renderCFields = () => {
    if (secretDetailLoading) {
      return (
        <div className="py-4">
          <Loading />
        </div>
      )
    }

    if (!selectedSecret) {
      return null
    }

    const fields: React.ReactNode[] = []

    if (selectedSecret?.value && !selectedSecret?.isNew) {
      fields.push(
        <div key="ui-key" className="grid gap-3">
          <Label htmlFor="s-uk">{t('ui.text.uiKey')}</Label>
          <CopyField value={selectedSecret.value} />
        </div>,
      )
    }

    fields.push(
      <div key={selectedSecret?.id || 0} className="grid gap-3">
        <Label htmlFor="s-t">
          {t('ui.text.title')}
          <span className="text-destructive">*</span>
        </Label>
        <Input
          id="s-t"
          name="s-t"
          placeholder={t('ui.text.myApiKey')}
          defaultValue={
            (selectedSecret?.label ? t(selectedSecret.label) : '') ||
            (selectedSecret?.title ? t(selectedSecret.title) : '') ||
            (selectedSecret?.name ? t(selectedSecret.name) : '') ||
            ''
          }
          aria-invalid={errors['s-t'] || undefined}
          className={errors['s-t'] ? 'border-destructive' : undefined}
          onChange={(e) => {
            clearError('s-t')
            updateSelectedSecret({
              key: 'label',
              value: e.target.value,
            })
          }}
        />
      </div>,
    )

    const canUsePlatform = selectedSecret?.isSupportPlatform
    const isNewCredential = selectedSecret?.isNew

    if (isNewCredential && canUsePlatform) {
      fields.push(
        <div
          key="platform-option"
          className={`flex items-center space-x-3 py-2 px-3 rounded-lg border ${usePlatformCredential ? 'bg-primary/10 border-primary/30' : 'bg-primary/5 border-primary/20'}`}
        >
          <Checkbox
            id="use-platform"
            checked={usePlatformCredential}
            onCheckedChange={setUsePlatformCredential}
          />
          <div className="flex flex-col gap-0.5">
            <Label
              htmlFor="use-platform"
              className="flex items-center gap-2 font-medium cursor-pointer"
            >
              <Building2 className="size-4 text-primary" />
              {t('ui.text.usePlatformCredential')}
            </Label>
            <span className="text-xs text-muted-foreground">
              {t('ui.text.usePlatformCredentialDescription')}
            </span>
          </div>
        </div>,
      )
    } else if (!isNewCredential && selectedSecret?.sourceType) {
      const isPlatform = selectedSecret.sourceType === 'platform'
      fields.push(
        <div key="source-type" className="grid gap-3">
          <Label>{t('ui.text.credentialSourceType')}</Label>
          <Badge
            variant={isPlatform ? 'secondary' : 'outline'}
            className={`gap-1 w-fit ${isPlatform ? 'bg-primary/10 text-primary border-primary/20' : ''}`}
          >
            {isPlatform ? (
              <Building2 className="size-3" />
            ) : (
              <User className="size-3" />
            )}
            {isPlatform
              ? t('ui.text.credentialSourceTypePlatform')
              : t('ui.text.credentialSourceTypeOwner')}
          </Badge>
        </div>,
      )
    }

    if (!usePlatformCredential && selectedSecret?.provider) {
      for (let i = 0; i < selectedSecret.provider.length; i++) {
        const schema = selectedSecret.provider[i]
        if (schema.hidden) {
          continue
        }

        let currentValue = ''
        if (selectedSecret.data && selectedSecret.data[schema.key]) {
          currentValue = selectedSecret.data[schema.key]
        } else if (
          selectedSecret.secretDetails?.data &&
          selectedSecret.secretDetails.data[schema.key]
        ) {
          currentValue = selectedSecret.secretDetails.data[schema.key]
        } else if (schema.default) {
          currentValue = schema.default
        }

        const isPassword = !!schema.writeOnly
        const isCopyable = !!schema.copyable
        const isVisible = visibleFields[schema.key]
        const hasAdornment = isPassword || isCopyable

        fields.push(
          <div key={schema.key} className="grid gap-3">
            <Label htmlFor={schema.key}>
              {t(schema.title)}
              {schema.required && <span className="text-destructive">*</span>}
            </Label>
            {isCopyable && !isPassword ? (
              <CopyField value={currentValue} />
            ) : (
              <div className={hasAdornment ? 'relative' : undefined}>
                <Input
                  id={schema.key}
                  name={schema.key}
                  defaultValue={currentValue}
                  disabled={schema.readOnly}
                  type={isPassword && !isVisible ? 'password' : 'text'}
                  autoComplete="off"
                  aria-invalid={errors[schema.key] || undefined}
                  className={
                    `${hasAdornment ? 'pr-10' : ''} ${errors[schema.key] ? 'border-destructive' : ''}`.trim() ||
                    undefined
                  }
                  onChange={(e) => {
                    clearError(schema.key)
                    updateSelectedSecret(
                      {
                        key: schema.key,
                        value: e.target.value,
                      },
                      true,
                    )
                  }}
                />
                {isPassword && (
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    tabIndex={-1}
                    aria-label={
                      isVisible ? t('ui.text.hide') : t('ui.text.show')
                    }
                    className="absolute right-1 top-1/2 -translate-y-1/2 h-7 w-7 p-0"
                    onClick={() => toggleVisibility(schema.key)}
                  >
                    {isVisible ? (
                      <EyeOff className="h-3.5 w-3.5" />
                    ) : (
                      <Eye className="h-3.5 w-3.5" />
                    )}
                  </Button>
                )}
              </div>
            )}
          </div>,
        )
      }
    }

    return fields
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{t(selectedSecret?.name)}</DialogTitle>
          <DialogDescription>
            {t(selectedSecret?.description)}
          </DialogDescription>
        </DialogHeader>

        {renderCFields()}

        <DialogFooter>
          <DialogClose asChild>
            <Button
              variant="outline"
              onClick={() => {
                onOpenChange(false)
              }}
            >
              {t('ui.text.cancel')}
            </Button>
          </DialogClose>
          {selectedSecret?.isOAuth2 ? (
            <>
              <Button variant="outline" onClick={handleSave}>
                {t('ui.text.save')}
              </Button>
              <Button
                variant="secondary"
                className="bg-primary/10 hover:bg-primary/20 text-primary"
                onClick={
                  selectedSecret?.isNew ? handleAuthorize : reauthorizeOAuth2
                }
              >
                {selectedSecret?.isNew
                  ? t('ui.text.authorize')
                  : t('ui.text.reauthorize')}
              </Button>
            </>
          ) : (
            <Button onClick={handleSave}>{t('ui.text.save')}</Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

CredentialFormDialog.displayName = 'CredentialFormDialog'

export { CredentialFormDialog }
