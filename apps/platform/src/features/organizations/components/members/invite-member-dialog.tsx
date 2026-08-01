import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { Users } from 'lucide-react'
import { Loading } from '@/components/shared/loading'
import { useTranslation } from 'react-i18next'

const InviteMemberDialog = ({
  availableProviders,
  providersLoading,
  inviteProvider,
  inviteIdentifier,
  inviteError,
  inviteLoading,
  onProviderChange,
  onIdentifierChange,
  onInvite,
}) => {
  const { t } = useTranslation()

  const selectedProvider = availableProviders.find(
    (p) => p.id === inviteProvider,
  )

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button className="cursor-pointer" size={'sm'}>
          <Users />
          {t('ui.text.invite')}
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{t('ui.text.inviteMember')}</DialogTitle>
          <DialogDescription>
            {t('ui.text.inviteMemberDescription')}
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 mt-4">
          <div className="grid gap-3">
            <Label htmlFor="provider">{t('ui.text.invitationMethod')}</Label>
            {providersLoading ? (
              <div className="flex items-center justify-center p-4">
                <Loading />
              </div>
            ) : (
              <Select value={inviteProvider} onValueChange={onProviderChange}>
                <SelectTrigger>
                  <SelectValue
                    placeholder={t('ui.text.selectInvitationMethod')}
                  />
                </SelectTrigger>
                <SelectContent>
                  {availableProviders.map((provider) => (
                    <SelectItem key={provider.id} value={provider.id}>
                      <div className="flex items-center gap-2">
                        {provider.icon}
                        <span>{provider.name}</span>
                      </div>
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          </div>
          <div className="grid gap-3">
            <Label htmlFor="identifier">{selectedProvider?.name}</Label>
            <Input
              id="identifier"
              name="identifier"
              placeholder={selectedProvider?.placeholder}
              value={inviteIdentifier}
              onChange={(e) => onIdentifierChange(e.target.value)}
              className={inviteError ? 'border-red-500' : ''}
            />
            {inviteError && (
              <p className="text-sm text-red-500">{inviteError}</p>
            )}
            <p className="text-xs text-muted-foreground">
              {selectedProvider?.description}
            </p>
          </div>
        </div>
        <DialogFooter className="-mb-2">
          <DialogClose asChild>
            <Button variant="outline">{t('ui.text.cancel')}</Button>
          </DialogClose>
          <Button
            onClick={onInvite}
            disabled={
              !inviteIdentifier.trim() || inviteLoading || providersLoading
            }
          >
            {inviteLoading ? t('ui.text.inviting') : t('ui.text.invite')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

InviteMemberDialog.displayName = 'InviteMemberDialog'

export { InviteMemberDialog }
