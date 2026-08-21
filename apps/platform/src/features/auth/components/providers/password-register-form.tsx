import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Eye, EyeOff, Loader2 } from 'lucide-react'
import { authService } from '@/features/auth'
import tokenManager from '@/lib/token-manager'
import { getReturnUrl } from '@/lib/auth-redirect'
import { useAuthRedirect } from '@/features/auth/hooks/use-auth-redirect'

const MIN_PASSWORD_LENGTH = 8

type Props = {
  onSubmitted?: (email: string) => void
}

const PasswordRegisterForm = ({ onSubmitted }: Props) => {
  const { t } = useTranslation()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const redirectAfterAuth = useAuthRedirect()
  const [errors, setErrors] = useState<Record<string, string>>({})

  const validate = () => {
    const next: Record<string, string> = {}
    if (!email.trim()) {
      next.email = t('ui.text.emailRequired')
    }
    if (password.length < MIN_PASSWORD_LENGTH) {
      next.password = t('ui.text.password_too_short')
    }
    if (password !== confirmPassword) {
      next.confirmPassword = t('ui.text.passwords_do_not_match')
    }
    setErrors(next)
    return Object.keys(next).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!validate()) {
      return
    }
    setSubmitting(true)
    try {
      const response = await authService.register({ email, password })
      if (!response.isSuccess) {
        return
      }
      const tokens = response.data
      if (tokens?.accessToken) {
        tokenManager.setTokens(tokens.accessToken, tokens.refreshToken)
        await redirectAfterAuth(getReturnUrl())
        return
      }
      onSubmitted?.(email)
    } catch (error) {
      console.error('Password register error:', error)
    } finally {
      setSubmitting(false)
    }
  }

  const clearError = (key: string) => {
    if (!errors[key]) {
      return
    }
    setErrors((prev) => {
      const next = { ...prev }
      delete next[key]
      return next
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <div className="grid gap-2">
        <Label htmlFor="register-email">{t('ui.text.emailAddress')}</Label>
        <Input
          id="register-email"
          type="email"
          autoComplete="email"
          value={email}
          aria-invalid={!!errors.email || undefined}
          className={errors.email ? 'border-destructive' : undefined}
          onChange={(e) => {
            clearError('email')
            setEmail(e.target.value)
          }}
          required
        />
        {errors.email && (
          <p className="text-xs text-destructive">{errors.email}</p>
        )}
      </div>
      <div className="grid gap-2">
        <Label htmlFor="register-password">{t('ui.text.password')}</Label>
        <div className="relative">
          <Input
            id="register-password"
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
            aria-label={showPassword ? t('ui.text.hide') : t('ui.text.show')}
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
        <Label htmlFor="register-confirm">{t('ui.text.confirmPassword')}</Label>
        <div className="relative">
          <Input
            id="register-confirm"
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
              showConfirmPassword ? t('ui.text.hide') : t('ui.text.show')
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
          <p className="text-xs text-destructive">{errors.confirmPassword}</p>
        )}
      </div>
      <Button type="submit" className="w-full" disabled={submitting}>
        {submitting ? <Loader2 className="w-4 h-4 animate-spin mr-2" /> : null}
        {t('ui.text.create_account')}
      </Button>
    </form>
  )
}

PasswordRegisterForm.displayName = 'PasswordRegisterForm'

export { PasswordRegisterForm }
