import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Eye, EyeOff, KeyRound, Loader2, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { authService } from '@/features/auth'

const MIN_PASSWORD_LENGTH = 8

type Props = {
  onLinked?: () => void | Promise<void>
}

const PasswordLinkedAccount = ({ onLinked }: Props) => {
  const { t } = useTranslation()
  const [expanded, setExpanded] = useState(false)
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [errors, setErrors] = useState<Record<string, string>>({})

  const reset = () => {
    setPassword('')
    setConfirmPassword('')
    setShowPassword(false)
    setShowConfirmPassword(false)
    setErrors({})
    setExpanded(false)
  }

  const validate = () => {
    const next: Record<string, string> = {}
    if (password.length < MIN_PASSWORD_LENGTH) {
      next.password = t('ui.text.password_too_short')
    }
    if (password !== confirmPassword) {
      next.confirmPassword = t('ui.text.passwords_do_not_match')
    }
    setErrors(next)
    return Object.keys(next).length === 0
  }

  const clearError = (key: string) => {
    if (!errors[key]) return
    setErrors((prev) => {
      const next = { ...prev }
      delete next[key]
      return next
    })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!validate()) return
    setSubmitting(true)
    try {
      const response = await authService.setPassword({ password })
      if (!response.isSuccess) {
        return
      }
      reset()
      await onLinked?.()
    } catch (error) {
      console.error('Set password error:', error)
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
              {t('ui.text.password', { ns: 'ui' })}
            </div>
            <div className="text-xs text-muted-foreground">
              {t('ui.text.set_password_description')}
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
            {t('ui.text.set_password')}
          </Button>
        )}
      </div>

      {expanded && (
        <form onSubmit={handleSubmit} className="space-y-3 mt-3">
          <div className="grid gap-2">
            <Label htmlFor="set-password">
              {t('ui.text.password', { ns: 'ui' })}
            </Label>
            <div className="relative">
              <Input
                id="set-password"
                type={showPassword ? 'text' : 'password'}
                autoComplete="new-password"
                value={password}
                aria-invalid={!!errors.password || undefined}
                className={`pr-10 ${errors.password ? 'border-destructive' : ''}`.trim()}
                onChange={(e) => {
                  clearError('password')
                  setPassword(e.target.value)
                }}
                required
                minLength={MIN_PASSWORD_LENGTH}
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
            {errors.password && (
              <p className="text-xs text-destructive">{errors.password}</p>
            )}
          </div>

          <div className="grid gap-2">
            <Label htmlFor="set-confirm-password">
              {t('ui.text.confirmPassword', { ns: 'ui' })}
            </Label>
            <div className="relative">
              <Input
                id="set-confirm-password"
                type={showConfirmPassword ? 'text' : 'password'}
                autoComplete="new-password"
                value={confirmPassword}
                aria-invalid={!!errors.confirmPassword || undefined}
                className={`pr-10 ${errors.confirmPassword ? 'border-destructive' : ''}`.trim()}
                onChange={(e) => {
                  clearError('confirmPassword')
                  setConfirmPassword(e.target.value)
                }}
                required
                minLength={MIN_PASSWORD_LENGTH}
              />
              <Button
                type="button"
                variant="ghost"
                size="sm"
                tabIndex={-1}
                aria-label={
                  showConfirmPassword
                    ? t('ui.text.hide', { ns: 'ui' })
                    : t('ui.text.show', { ns: 'ui' })
                }
                className="absolute right-1 top-1/2 -translate-y-1/2 h-7 w-7 p-0"
                onClick={() => setShowConfirmPassword((s) => !s)}
              >
                {showConfirmPassword ? (
                  <EyeOff className="h-3.5 w-3.5" />
                ) : (
                  <Eye className="h-3.5 w-3.5" />
                )}
              </Button>
            </div>
            {errors.confirmPassword && (
              <p className="text-xs text-destructive">
                {errors.confirmPassword}
              </p>
            )}
          </div>

          <Button type="submit" className="w-full" disabled={submitting}>
            {submitting && <Loader2 className="w-4 h-4 animate-spin mr-2" />}
            {t('ui.text.set_password')}
          </Button>
        </form>
      )}
    </div>
  )
}

PasswordLinkedAccount.displayName = 'PasswordLinkedAccount'

export { PasswordLinkedAccount }
