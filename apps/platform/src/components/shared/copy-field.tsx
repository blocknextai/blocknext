import { Copy, Check } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { useCopy } from '@/hooks/use-copy'
import { cn } from '@/lib/utils'

interface CopyFieldProps {
  /** The text written to the clipboard (and shown in the field). */
  value: string
  className?: string
  inputClassName?: string
  iconClassName?: string
}

/**
 * A readonly value field with click-anywhere copy: clicking the input or the
 * copy icon copies the value. The input keeps long values fully readable
 * (scrollable) so they never get truncated or overflow their container.
 */
export function CopyField({
  value,
  className,
  inputClassName,
  iconClassName,
}: CopyFieldProps) {
  const { t } = useTranslation()
  const { copied, copy } = useCopy()
  const label = t('ui.text.copy', 'Copy')

  return (
    <div className={cn('relative w-full', className)}>
      <Input
        readOnly
        value={value}
        onClick={() => copy(value)}
        title={label}
        className={cn(
          'pr-9 bg-muted/50 text-xs cursor-pointer',
          inputClassName,
        )}
      />
      <Button
        type="button"
        variant="ghost"
        size="icon-sm"
        onClick={() => copy(value)}
        aria-label={label}
        title={label}
        className="absolute right-1 top-1/2 -translate-y-1/2"
      >
        {copied ? (
          <Check
            className={cn('text-green-500', iconClassName ?? 'size-3.5')}
          />
        ) : (
          <Copy
            className={cn('text-muted-foreground', iconClassName ?? 'size-3.5')}
          />
        )}
      </Button>
    </div>
  )
}

CopyField.displayName = 'CopyField'
