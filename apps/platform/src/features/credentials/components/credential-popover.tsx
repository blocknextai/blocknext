import { useMemo, useState } from 'react'
import { Button } from '@/components/ui/button'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { useTranslation } from 'react-i18next'
import { KeyRound } from 'lucide-react'
import { SearchInput } from '@/components/shared/search-input'

const CredentialPopover = ({
  children,
  secrets,
  onSelect,
  isUserMode = false,
}) => {
  const { t } = useTranslation()
  const [query, setQuery] = useState('')
  const [open, setOpen] = useState(false)

  const handleOpenChange = (next: boolean) => {
    setOpen(next)
    if (!next) {
      setQuery('')
    }
  }

  const filteredSecrets = useMemo(() => {
    if (!secrets) {
      return []
    }
    const q = query.trim().toLowerCase()
    if (!q) {
      return secrets
    }
    return secrets.filter((secret: any) =>
      t(secret.name).toLowerCase().includes(q),
    )
  }, [secrets, query, t])

  const renderSecrets = () => {
    const s: React.ReactNode[] = []
    for (let i = 0; i < filteredSecrets.length; i++) {
      const secret = filteredSecrets[i]
      const IconComponent = secret.icon || KeyRound
      s.push(
        <Button
          key={i}
          variant={'ghost'}
          className="text-xs w-full justify-start gap-2"
          onClick={() => {
            handleOpenChange(false)
            const secretCopy = { ...secret, isNew: true }
            onSelect(secretCopy)
          }}
        >
          <IconComponent className="text-icon size-4 shrink-0" />
          <span>{t(secret.name)}</span>
        </Button>,
      )
    }
    return s
  }

  return (
    <Popover open={open} onOpenChange={handleOpenChange}>
      <PopoverTrigger asChild>{children}</PopoverTrigger>
      <PopoverContent className="p-0" asChild>
        <div className={`flex flex-col w-80 ${isUserMode ? 'pb-1' : ''}`}>
          <div className="p-2 border-b">
            <SearchInput
              placeholder={t('ui.text.searchCredentials')}
              onSearch={setQuery}
              autoFocus
            />
          </div>
          <div className="flex flex-col gap-1 overflow-y-auto max-h-72 p-1">
            {filteredSecrets.length > 0 ? (
              renderSecrets()
            ) : (
              <div className="text-xs text-muted-foreground text-center p-6">
                {t('ui.text.noCredentials')}
              </div>
            )}
          </div>
        </div>
      </PopoverContent>
    </Popover>
  )
}

CredentialPopover.displayName = 'CredentialPopover'

export { CredentialPopover }
