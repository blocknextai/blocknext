import { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import {
  Link as LinkIcon,
  Trash2,
  AlertTriangle,
  BadgeCheck,
  CircleAlert,
  MoreHorizontal,
  MailCheck,
  Loader2,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { CopyField } from '@/components/shared/copy-field'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { ProviderLoginButton } from '@/features/auth/components/provider-login-button'
import { EmailLinkedAccount } from '@/features/auth/components/providers/email-linked-account'
import { PasswordLinkedAccount } from '@/features/auth/components/providers/password-linked-account'
import { ChangePasswordForm } from '@/features/auth/components/providers/change-password-form'
import { ChangeEmailForm } from '@/features/auth/components/providers/change-email-form'
import { ProviderIcon } from '@/features/auth/components/provider-icon'
import { getProviderName } from '@/features/auth/components/provider-utils'

const ProfileSection = ({
  userId,
  loading,
  linkedAccounts,
  deletingAccountId,
  deleteLinkedAccount,
  getUnlinkedProviders,
  handleAccountLinked,
  passwordEnabled,
  resendingEmail,
  onResendVerification,
}) => {
  const { t } = useTranslation()
  const [confirmRemoveId, setConfirmRemoveId] = useState<string | null>(null)

  const hasEmail = useMemo(
    () => linkedAccounts.some((a) => a.authProvider === 'email'),
    [linkedAccounts],
  )
  const hasPassword = useMemo(
    () => linkedAccounts.some((a) => a.authProvider === 'password'),
    [linkedAccounts],
  )

  return (
    <div className="flex flex-col gap-4 w-full min-w-0 p-2 sm:p-4">
      {userId && (
        <>
          <div className="flex justify-between items-center">
            <div className="grid gap-3 flex-1">
              <Label htmlFor="user-id">{t('ui.text.userId')}</Label>
              <CopyField value={userId} />
            </div>
          </div>
          <Separator className="my-2 p-0" decorative={true} />
        </>
      )}

      <div className="flex flex-col gap-4">
        <div className="grid gap-3">
          <Label htmlFor={`linked-accounts`}>
            {t('ui.text.linkedAccounts')}
          </Label>
          <span id="linked-accounts" className="text-muted-foreground text-sm!">
            {t('ui.text.manageConnectedAccounts')}
          </span>
        </div>

        {loading ? (
          <div className="flex items-center justify-center py-8">
            <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary"></div>
          </div>
        ) : linkedAccounts.length > 0 ? (
          <div className="space-y-3">
            {linkedAccounts.map((account, index) => (
              <div
                key={index}
                className="flex items-center justify-between p-3 bg-muted/50 rounded-lg"
              >
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 bg-primary/10 rounded-full flex items-center justify-center">
                    <ProviderIcon
                      provider={account.authProvider}
                      className="w-4 h-4"
                    />
                  </div>
                  <div>
                    <div className="font-medium text-sm capitalize">
                      {getProviderName(account.authProvider)}
                    </div>
                    {(() => {
                      const { displayName, identifier, authProvider } = account
                      const isPlaceholder =
                        !identifier || identifier === authProvider
                      let subtitle: string | null = null
                      if (
                        displayName &&
                        identifier &&
                        displayName !== identifier
                      ) {
                        subtitle = `${displayName} • ${identifier}`
                      } else if (displayName) {
                        subtitle = displayName
                      } else if (!isPlaceholder) {
                        subtitle = identifier
                      }
                      if (!subtitle) {
                        return null
                      }
                      return (
                        <div className="text-xs text-muted-foreground">
                          {subtitle}
                        </div>
                      )
                    })()}
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  {account.isPrimary && (
                    <span className="px-2 py-1 text-xs bg-primary text-primary-foreground rounded-full font-medium">
                      {t('ui.text.primary')}
                    </span>
                  )}
                  {account.isVerified && (
                    <span className="inline-flex items-center gap-1 px-2 py-1 text-xs rounded-full font-medium bg-green-500/10 text-green-700 dark:text-green-300">
                      <BadgeCheck className="w-3 h-3" />
                      {t('ui.text.verified')}
                    </span>
                  )}
                  {account.isVerified === false && (
                    <span className="inline-flex items-center gap-1 px-2 py-1 text-xs rounded-full font-medium bg-yellow-500/10 text-yellow-700 dark:text-yellow-300">
                      <CircleAlert className="w-3 h-3" />
                      {t('ui.text.unverified')}
                    </span>
                  )}
                  {(() => {
                    const canResend =
                      account.authProvider === 'email' && !account.isVerified
                    const canRemove = !account.isPrimary
                    if (!canResend && !canRemove) {
                      return null
                    }
                    const isBusy =
                      resendingEmail === account.identifier ||
                      deletingAccountId === account.id
                    return (
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button
                            variant="ghost"
                            size="sm"
                            className="h-8 w-8 p-0"
                            disabled={isBusy}
                          >
                            {isBusy ? (
                              <Loader2 className="w-4 h-4 animate-spin" />
                            ) : (
                              <MoreHorizontal className="w-4 h-4" />
                            )}
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          {canResend && (
                            <DropdownMenuItem
                              onClick={() =>
                                onResendVerification?.(account.identifier)
                              }
                            >
                              <MailCheck className="w-4 h-4 mr-2" />
                              {t('ui.text.resendVerification')}
                            </DropdownMenuItem>
                          )}
                          {canRemove && canResend && <DropdownMenuSeparator />}
                          {canRemove && (
                            <DropdownMenuItem
                              variant="destructive"
                              onClick={() => setConfirmRemoveId(account.id)}
                            >
                              <Trash2 className="w-4 h-4 mr-2" />
                              {t('ui.text.remove')}
                            </DropdownMenuItem>
                          )}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    )
                  })()}
                </div>
              </div>
            ))}

            <AlertDialog
              open={confirmRemoveId !== null}
              onOpenChange={(open) => !open && setConfirmRemoveId(null)}
            >
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle className="flex items-center gap-2">
                    <AlertTriangle className="w-5 h-5 text-destructive" />
                    {t('ui.text.removeLinkedAccount')}
                  </AlertDialogTitle>
                  <AlertDialogDescription>
                    {(() => {
                      const target = linkedAccounts.find(
                        (a) => a.id === confirmRemoveId,
                      )
                      if (!target) {
                        return null
                      }
                      return t('ui.text.removeLinkedAccountConfirmation', {
                        name: target.displayName ?? target.identifier,
                        provider: target.authProvider,
                      })
                    })()}
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>{t('ui.text.cancel')}</AlertDialogCancel>
                  <AlertDialogAction
                    onClick={() => {
                      if (confirmRemoveId) {
                        deleteLinkedAccount(confirmRemoveId)
                      }
                      setConfirmRemoveId(null)
                    }}
                    className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                  >
                    {t('ui.text.remove')}
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </div>
        ) : (
          <div className="text-center py-8 text-muted-foreground">
            <LinkIcon className="w-8 h-8 mx-auto mb-2 opacity-50" />
            <p className="text-sm">{t('ui.text.noLinkedAccounts')}</p>
          </div>
        )}

        {(() => {
          const unlinked = getUnlinkedProviders()
          const showEmailRow = !hasEmail
          const showPasswordRow = passwordEnabled && hasEmail && !hasPassword
          if (unlinked.length === 0 && !showEmailRow && !showPasswordRow) {
            return null
          }
          return (
            <>
              <Separator className="my-4" />
              <div className="grid gap-3">
                <Label htmlFor={`unlinked-accounts`}>
                  {t('ui.text.connectAdditionalAccounts')}
                </Label>
                <span
                  id="unlinked-accounts"
                  className="text-muted-foreground text-sm!"
                >
                  {t('ui.text.connectOtherAccounts')}
                </span>
              </div>
              <div className="space-y-3">
                {unlinked.map((provider) => (
                  <div
                    key={provider}
                    className="flex items-center justify-between p-3 bg-muted/30 rounded-lg border border-dashed"
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 bg-primary/10 rounded-full flex items-center justify-center">
                        <ProviderIcon provider={provider} className="w-4 h-4" />
                      </div>
                      <div>
                        <div className="font-medium text-sm capitalize">
                          {getProviderName(provider)}
                        </div>
                        <div className="text-xs text-muted-foreground">
                          {t('ui.text.connectYourAccount', {
                            provider: getProviderName(provider),
                          })}
                        </div>
                      </div>
                    </div>
                    <ProviderLoginButton
                      provider={provider}
                      mode="linked_accounts"
                    />
                  </div>
                ))}
                {showEmailRow && (
                  <EmailLinkedAccount onLinked={handleAccountLinked} />
                )}
                {showPasswordRow && (
                  <PasswordLinkedAccount onLinked={handleAccountLinked} />
                )}
              </div>
            </>
          )
        })()}

        {hasPassword && (
          <>
            <Separator className="my-4" />
            <div className="grid gap-3">
              <Label htmlFor="security-section">{t('ui.text.security')}</Label>
              <span
                id="security-section"
                className="text-muted-foreground text-sm!"
              >
                {t('ui.text.securityDescription')}
              </span>
            </div>
            <div className="space-y-3">
              <ChangePasswordForm />
              <ChangeEmailForm />
            </div>
          </>
        )}
      </div>
    </div>
  )
}

ProfileSection.displayName = 'ProfileSection'

export { ProfileSection }
