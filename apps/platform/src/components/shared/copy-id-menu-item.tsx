import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Copy, Check } from 'lucide-react'
import { DropdownMenuItem } from '@/components/ui/dropdown-menu'

const CopyIdMenuItem = ({ id }: { id: string }) => {
  const { t } = useTranslation()
  const [copied, setCopied] = useState(false)

  const handleCopy = (e: Event) => {
    e.preventDefault()
    navigator.clipboard.writeText(id)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <DropdownMenuItem onSelect={handleCopy}>
      {copied ? <Check className="text-green-500" /> : <Copy />}
      {t('ui.text.copyId')}
    </DropdownMenuItem>
  )
}

CopyIdMenuItem.displayName = 'CopyIdMenuItem'

export { CopyIdMenuItem }
