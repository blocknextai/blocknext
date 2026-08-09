import { Check, Copy } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { useCopy } from '@/hooks/use-copy'

export function CodeBlock({ code }: { code: string }) {
  const { t } = useTranslation()
  const { copied, copy } = useCopy()
  const label = t('ui.text.copy', 'Copy')

  return (
    <div className="overflow-hidden rounded-lg border bg-muted/50">
      <div className="flex items-center justify-end border-b bg-muted/50 px-2 py-1">
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={() => copy(code)}
          aria-label={label}
          title={label}
          className="h-6 gap-1.5 px-2 text-xs text-muted-foreground"
        >
          {copied ? (
            <Check className="size-3.5 text-green-500" />
          ) : (
            <Copy className="size-3.5" />
          )}
          {label}
        </Button>
      </div>
      <pre className="max-h-[55vh] overflow-auto p-3 font-mono text-xs leading-relaxed">
        <code>{code}</code>
      </pre>
    </div>
  )
}
