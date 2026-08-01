import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Loader2, Mail, MailCheck, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { authService } from '@/features/auth'

type Props = {
  onLinked?: () => void | Promise<void>
}

const EmailLinkedAccount = ({ onLinked }: Props) => {
  const { t } = useTranslation()
  const [expanded, setExpanded] = useState(false)
  const [email, setEmail] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [sent, setSent] = useState(false)

  const reset = () => {
    setEmail('')
    setError(null)
    setSent(false)
    setExpanded(false)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!email.trim()) {
      setError(t('ui.text.emailRequired', { ns: 'ui' }))
      return
    }
    setSubmitting(true)
    try {
      const response = await authService.addEmail({ email })
      if (!response.isSuccess) {
        return
      }
      setSent(true)
      await onLinked?.()
    } catch (err) {
      console.error('Add email error:', err)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="p-3 bg-muted/30 rounded-lg border border-dashed">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-primary/10 rounded-full flex items-center justify-center">
            <Mail className="w-4 h-4" />
          </div>
          <div>
            <div className="font-medium text-sm">
              {t('ui.text.emailAddress', { ns: 'ui' })}
            </div>
            <div className="text-xs text-muted-foreground">
              {t('ui.text.add_email_description')}
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
            {t('ui.text.add_email')}
          </Button>
        )}
      </div>

      {expanded && sent && (
        <div className="flex items-start gap-3 p-3 rounded-lg bg-green-500/10 border border-green-500/20 text-sm mt-3">
          <MailCheck className="size-5 text-green-600 dark:text-green-400 shrink-0 mt-0.5" />
          <div className="text-green-700 dark:text-green-300">
            <p className="font-medium">{t('ui.text.add_email_sent')}</p>
            <p className="mt-1 text-xs opacity-90">
              {t('ui.text.add_email_sent_description', { email })}
            </p>
          </div>
        </div>
      )}

      {expanded && !sent && (
        <form onSubmit={handleSubmit} className="space-y-3 mt-3">
          <div className="grid gap-2">
            <Label htmlFor="add-email-input">
              {t('ui.text.emailAddress', { ns: 'ui' })}
            </Label>
            <Input
              id="add-email-input"
              type="email"
              autoComplete="email"
              value={email}
              aria-invalid={!!error || undefined}
              className={error ? 'border-destructive' : undefined}
              onChange={(e) => {
                setError(null)
                setEmail(e.target.value)
              }}
              required
            />
            {error && <p className="text-xs text-destructive">{error}</p>}
          </div>
          <Button type="submit" className="w-full" disabled={submitting}>
            {submitting && <Loader2 className="w-4 h-4 animate-spin mr-2" />}
            {t('ui.text.send_verification_email')}
          </Button>
        </form>
      )}
    </div>
  )
}

EmailLinkedAccount.displayName = 'EmailLinkedAccount'

export { EmailLinkedAccount }
