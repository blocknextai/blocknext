import { useState } from 'react'
import { Link, useSearchParams, useNavigate } from 'react-router'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Eye, EyeOff, Loader2, AlertCircle } from 'lucide-react'
import { authService } from '@/features/auth'
import { BackToSignInLink } from '@/features/auth/components/back-to-sign-in-link'

const MIN_PASSWORD_LENGTH = 8

function AuthResetPasswordPage() {
  const { t } = useTranslation()
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const token = searchParams.get('token') ?? ''

  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [errors, setErrors] = useState<Record<string, string>>({})

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

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!validate()) return
    setSubmitting(true)
    try {
      const response = await authService.resetPassword({
        token,
        newPassword: password,
      })
      if (!response.isSuccess) return
      navigate('/preferences/profile')
    } catch (error) {
      console.error('Password reset error:', error)
    } finally {
      setSubmitting(false)
    }
  }

  if (!token) {
    return (
      <>
        <div className="mb-6 flex flex-col gap-y-3">
          <h1 className="text-2xl font-bold md:text-3xl">
            {t('ui.text.reset_password_title')}
          </h1>
        </div>
        <BackToSignInLink to="/auth/login" />
        <div className="flex items-start gap-3 p-4 rounded-lg bg-red-500/10 border border-red-500/20 text-sm">
          <AlertCircle className="size-5 text-red-600 dark:text-red-400 shrink-0 mt-0.5" />
          <span className="text-red-700 dark:text-red-300">
            {t('ui.text.reset_token_missing')}
          </span>
        </div>
        <Button asChild variant="outline" className="w-full mt-4">
          <Link to="/auth/forgot-password">
            {t('ui.text.request_new_link')}
          </Link>
        </Button>
      </>
    )
  }

  return (
    <>
      <div className="mb-6 flex flex-col gap-y-3">
        <h1 className="text-2xl font-bold md:text-3xl">
          {t('ui.text.reset_password_title')}
        </h1>
        <p className="text-sm text-muted-foreground">
          {t('ui.text.reset_password_description')}
        </p>
      </div>

      <BackToSignInLink to="/auth/login" />

      <form onSubmit={handleSubmit} className="space-y-3">
        <div className="grid gap-2">
          <Label htmlFor="reset-password">
            {t('ui.text.password', { ns: 'ui' })}
          </Label>
          <div className="relative">
            <Input
              id="reset-password"
              type={showPassword ? 'text' : 'password'}
              autoComplete="new-password"
              value={password}
              aria-invalid={!!errors.password || undefined}
              className={`pr-10 ${errors.password ? 'border-destructive' : ''}`.trim()}
              onChange={(e) => {
                setPassword(e.target.value)
                if (errors.password) setErrors({ ...errors, password: '' })
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
          <Label htmlFor="reset-confirm">
            {t('ui.text.confirmPassword', { ns: 'ui' })}
          </Label>
          <Input
            id="reset-confirm"
            type={showPassword ? 'text' : 'password'}
            autoComplete="new-password"
            value={confirmPassword}
            aria-invalid={!!errors.confirmPassword || undefined}
            className={
              errors.confirmPassword ? 'border-destructive' : undefined
            }
            onChange={(e) => {
              setConfirmPassword(e.target.value)
              if (errors.confirmPassword)
                setErrors({ ...errors, confirmPassword: '' })
            }}
            required
            minLength={MIN_PASSWORD_LENGTH}
          />
          {errors.confirmPassword && (
            <p className="text-xs text-destructive">{errors.confirmPassword}</p>
          )}
        </div>
        <Button type="submit" className="w-full" disabled={submitting}>
          {submitting ? (
            <Loader2 className="w-4 h-4 animate-spin mr-2" />
          ) : null}
          {t('ui.text.reset_password')}
        </Button>
      </form>
    </>
  )
}

export default AuthResetPasswordPage
