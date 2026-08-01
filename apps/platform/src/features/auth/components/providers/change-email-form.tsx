import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { AtSign, Eye, EyeOff, Loader2, MailCheck, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { authService } from '@/features/auth'

const ChangeEmailForm = () => {
  const { t } = useTranslation()
  const [expanded, setExpanded] = useState(false)
  const [newEmail, setNewEmail] = useState('')
  const [currentPassword, setCurrentPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [sent, setSent] = useState(false)

  const reset = () => {
    setNewEmail('')
    setCurrentPassword('')
    setShowPassword(false)
    setErrors({})
    setExpanded(false)
    setSent(false)
  }

  const clearError = (key: string) => {
    if (!errors[key]) return
    setErrors((prev) => {
      const next = { ...prev }
      delete next[key]
      return next
    })
  }

  const validate = () => {
    const next: Record<string, string> = {}
    if (!newEmail.trim()) {
      next.newEmail = t('ui.text.emailRequired', { ns: 'ui' })
    }
    if (!currentPassword) {
      next.currentPassword = t('ui.text.current_password_required')
    }
    setErrors(next)
    return Object.keys(next).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!validate()) return
    setSubmitting(true)
    try {
      const response = await authService.changeEmail({
        newEmail,
        currentPassword,
      })
      if (!response.isSuccess) {
        return
      }
      setSent(true)
    } catch (error) {
      console.error('Change email error:', error)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="p-3 bg-muted/30 rounded-lg border border-dashed">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-primary/10 rounded-full flex items-center justify-center">
            <AtSign className="w-4 h-4" />
          </div>
          <div>
            <div className="font-medium text-sm">
              {t('ui.text.change_email')}
            </div>
            <div className="text-xs text-muted-foreground">
              {t('ui.text.change_email_description')}
            </div>
          </div>
        </div>
        {expanded ? (
          <Button
            variant="ghost"
            size="sm"
            onClick={reset}
            aria-label={t('ui.text.cancel', { ns: 'ui' })}
            disabled={submitting}
          >
            <X className="w-4 h-4" />
          </Button>
        ) : (
          <Button variant="outline" onClick={() => setExpanded(true)}>
            {t('ui.text.change_email')}
          </Button>
        )}
      </div>

      {expanded && sent && (
        <div className="flex items-start gap-3 p-3 rounded-lg bg-green-500/10 border border-green-500/20 text-sm mt-3">
          <MailCheck className="size-5 text-green-600 dark:text-green-400 shrink-0 mt-0.5" />
          <div className="text-green-700 dark:text-green-300">
            <p className="font-medium">{t('ui.text.change_email_sent')}</p>
            <p className="mt-1 text-xs opacity-90">
              {t('ui.text.change_email_sent_description', {
                email: newEmail,
              })}
            </p>
          </div>
        </div>
      )}

      {expanded && !sent && (
        <form onSubmit={handleSubmit} className="space-y-3 mt-3">
          <div className="grid gap-2">
            <Label htmlFor="change-email-new">{t('ui.text.new_email')}</Label>
            <Input
              id="change-email-new"
              type="email"
              autoComplete="email"
              value={newEmail}
              aria-invalid={!!errors.newEmail || undefined}
              className={errors.newEmail ? 'border-destructive' : undefined}
              onChange={(e) => {
                clearError('newEmail')
                setNewEmail(e.target.value)
              }}
              required
            />
            {errors.newEmail && (
              <p className="text-xs text-destructive">{errors.newEmail}</p>
            )}
          </div>
          <div className="grid gap-2">
            <Label htmlFor="change-email-password">
              {t('ui.text.current_password')}
            </Label>
            <div className="relative">
              <Input
                id="change-email-password"
                type={showPassword ? 'text' : 'password'}
                autoComplete="current-password"
                value={currentPassword}
                aria-invalid={!!errors.currentPassword || undefined}
                className={`pr-10 ${errors.currentPassword ? 'border-destructive' : ''}`.trim()}
                onChange={(e) => {
                  clearError('currentPassword')
                  setCurrentPassword(e.target.value)
                }}
                required
              />
              <Button
                type="button"
                variant="ghost"
                size="sm"
                tabIndex={-1}
                aria-label={
                  showPassword
                    ? t('ui.text.hide', { ns: 'ui' })
                    : t('ui.text.show', { ns: 'ui' })
                }
                className="absolute right-1 top-1/2 -translate-y-1/2 h-7 w-7 p-0"
                onClick={() => setShowPassword((s) => !s)}
              >
                {showPassword ? (
                  <EyeOff className="h-3.5 w-3.5" />
                ) : (
                  <Eye className="h-3.5 w-3.5" />
                )}
              </Button>
            </div>
            {errors.currentPassword && (
              <p className="text-xs text-destructive">
                {errors.currentPassword}
              </p>
            )}
          </div>
          <Button type="submit" disabled={submitting}>
            {submitting && <Loader2 className="w-4 h-4 animate-spin mr-2" />}
            {t('ui.text.send_confirmation_email')}
          </Button>
        </form>
      )}
    </div>
  )
}

ChangeEmailForm.displayName = 'ChangeEmailForm'

export { ChangeEmailForm }
