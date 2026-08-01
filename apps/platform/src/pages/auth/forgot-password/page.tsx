import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Loader2, CheckCircle2 } from 'lucide-react'
import { authService } from '@/features/auth'
import { BackToSignInLink } from '@/features/auth/components/back-to-sign-in-link'

function AuthForgotPasswordPage() {
  const { t } = useTranslation()
  const [email, setEmail] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [submitted, setSubmitted] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!email.trim()) return
    setSubmitting(true)
    try {
      const response = await authService.requestPasswordReset({ email })
      if (!response.isSuccess) return
      setSubmitted(true)
    } catch (error) {
      console.error('Password reset request error:', error)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <>
      <div className="mb-6 flex flex-col gap-y-3">
        <h1 className="text-2xl font-bold md:text-3xl">
          {t('ui.text.forgot_password_title')}
        </h1>
        <p className="text-sm text-muted-foreground">
          {t('ui.text.forgot_password_description')}
        </p>
      </div>

      <BackToSignInLink to="/auth/login" />

      {submitted ? (
        <div className="flex items-start gap-3 p-4 rounded-lg bg-green-500/10 border border-green-500/20 text-sm">
          <CheckCircle2 className="size-5 text-green-600 dark:text-green-400 shrink-0 mt-0.5" />
          <span className="text-green-700 dark:text-green-300">
            {t('ui.text.password_reset_sent')}
          </span>
        </div>
      ) : (
        <form onSubmit={handleSubmit} className="space-y-3">
          <div className="grid gap-2">
            <Label htmlFor="forgot-email">
              {t('ui.text.emailAddress', { ns: 'ui' })}
            </Label>
            <Input
              id="forgot-email"
              type="email"
              autoComplete="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>
          <Button type="submit" className="w-full" disabled={submitting}>
            {submitting ? (
              <Loader2 className="w-4 h-4 animate-spin mr-2" />
            ) : null}
            {t('ui.text.send_reset_link')}
          </Button>
        </form>
      )}
    </>
  )
}

export default AuthForgotPasswordPage
