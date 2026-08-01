import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Eye, EyeOff, KeyRound, Loader2, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { authService } from '@/features/auth'

const MIN_PASSWORD_LENGTH = 8

const ChangePasswordForm = () => {
  const { t } = useTranslation()
  const [expanded, setExpanded] = useState(false)
  const [currentPassword, setCurrentPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [showCurrent, setShowCurrent] = useState(false)
  const [showNew, setShowNew] = useState(false)
  const [showConfirm, setShowConfirm] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [errors, setErrors] = useState<Record<string, string>>({})

  const reset = () => {
    setCurrentPassword('')
    setNewPassword('')
    setConfirmPassword('')
    setShowCurrent(false)
    setShowNew(false)
    setShowConfirm(false)
    setErrors({})
    setExpanded(false)
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
    if (!currentPassword) {
      next.currentPassword = t('ui.text.current_password_required')
    }
    if (newPassword.length < MIN_PASSWORD_LENGTH) {
      next.newPassword = t('ui.text.password_too_short')
    }
    if (newPassword !== confirmPassword) {
      next.confirmPassword = t('ui.text.passwords_do_not_match')
    }
    setErrors(next)
    return Object.keys(next).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!validate()) return
    setSubmitting(true)
    try {
      const response = await authService.changePassword({
        currentPassword,
        newPassword,
      })
      if (!response.isSuccess) {
        return
      }
      reset()
    } catch (error) {
      console.error('Change password error:', error)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="p-3 bg-muted/30 rounded-lg border border-dashed">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-primary/10 rounded-full flex items-center justify-center">
            <KeyRound className="w-4 h-4" />
          </div>
          <div>
            <div className="font-medium text-sm">
              {t('ui.text.change_password')}
            </div>
            <div className="text-xs text-muted-foreground">
              {t('ui.text.change_password_description')}
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
            {t('ui.text.change_password')}
          </Button>
        )}
      </div>

      {expanded && (
        <form onSubmit={handleSubmit} className="space-y-3 mt-3">
          <PasswordField
            id="change-current-password"
            label={t('ui.text.current_password')}
            value={currentPassword}
            onChange={(v) => {
              clearError('currentPassword')
              setCurrentPassword(v)
            }}
            error={errors.currentPassword}
            show={showCurrent}
            toggleShow={() => setShowCurrent((s) => !s)}
            autoComplete="current-password"
            t={t}
          />
          <PasswordField
            id="change-new-password"
            label={t('ui.text.new_password')}
            value={newPassword}
            onChange={(v) => {
              clearError('newPassword')
              setNewPassword(v)
            }}
            error={errors.newPassword}
            show={showNew}
            toggleShow={() => setShowNew((s) => !s)}
            autoComplete="new-password"
            minLength={MIN_PASSWORD_LENGTH}
            t={t}
          />
          <PasswordField
            id="change-confirm-password"
            label={t('ui.text.confirmPassword', { ns: 'ui' })}
            value={confirmPassword}
            onChange={(v) => {
              clearError('confirmPassword')
              setConfirmPassword(v)
            }}
            error={errors.confirmPassword}
            show={showConfirm}
            toggleShow={() => setShowConfirm((s) => !s)}
            autoComplete="new-password"
            minLength={MIN_PASSWORD_LENGTH}
            t={t}
          />
          <Button type="submit" disabled={submitting}>
            {submitting && <Loader2 className="w-4 h-4 animate-spin mr-2" />}
            {t('ui.text.change_password')}
          </Button>
        </form>
      )}
    </div>
  )
}

type PasswordFieldProps = {
  id: string
  label: string
  value: string
  onChange: (v: string) => void
  error?: string
  show: boolean
  toggleShow: () => void
  autoComplete: string
  minLength?: number
  t: (k: string, opts?: object) => string
}

const PasswordField = ({
  id,
  label,
  value,
  onChange,
  error,
  show,
  toggleShow,
  autoComplete,
  minLength,
  t,
}: PasswordFieldProps) => (
  <div className="grid gap-2">
    <Label htmlFor={id}>{label}</Label>
    <div className="relative">
      <Input
        id={id}
        type={show ? 'text' : 'password'}
        autoComplete={autoComplete}
        value={value}
        aria-invalid={!!error || undefined}
        className={`pr-10 ${error ? 'border-destructive' : ''}`.trim()}
        onChange={(e) => onChange(e.target.value)}
        required
        minLength={minLength}
      />
      <Button
        type="button"
        variant="ghost"
        size="sm"
        tabIndex={-1}
        aria-label={
          show
            ? t('ui.text.hide', { ns: 'ui' })
            : t('ui.text.show', { ns: 'ui' })
        }
        className="absolute right-1 top-1/2 -translate-y-1/2 h-7 w-7 p-0"
        onClick={toggleShow}
      >
        {show ? (
          <EyeOff className="h-3.5 w-3.5" />
        ) : (
          <Eye className="h-3.5 w-3.5" />
        )}
      </Button>
    </div>
    {error && <p className="text-xs text-destructive">{error}</p>}
  </div>
)

ChangePasswordForm.displayName = 'ChangePasswordForm'

export { ChangePasswordForm }
