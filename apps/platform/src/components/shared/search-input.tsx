import { useEffect, useRef, useState } from 'react'
import type { ChangeEvent } from 'react'
import { useTranslation } from 'react-i18next'
import { Search, X } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

interface SearchInputProps {
  placeholder?: string
  defaultValue?: string
  onSearch: (value: string) => void
  debounceMs?: number
  className?: string
  autoFocus?: boolean
}

export function SearchInput({
  placeholder,
  defaultValue = '',
  onSearch,
  debounceMs = 300,
  className,
  autoFocus = false,
}: SearchInputProps) {
  const { t } = useTranslation()
  const [value, setValue] = useState(defaultValue)
  const onSearchRef = useRef(onSearch)
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  useEffect(() => {
    onSearchRef.current = onSearch
  }, [onSearch])

  useEffect(() => {
    return () => {
      if (timeoutRef.current) clearTimeout(timeoutRef.current)
    }
  }, [])

  const fireSearch = (next: string) => {
    if (timeoutRef.current) clearTimeout(timeoutRef.current)
    timeoutRef.current = setTimeout(() => {
      onSearchRef.current(next)
    }, debounceMs)
  }

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const next = e.target.value
    setValue(next)
    fireSearch(next)
  }

  const handleClear = () => {
    setValue('')
    if (timeoutRef.current) clearTimeout(timeoutRef.current)
    onSearchRef.current('')
  }

  return (
    <div className={cn('relative w-full', className)}>
      <Search className="pointer-events-none absolute left-2.5 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
      <Input
        type="text"
        autoFocus={autoFocus}
        placeholder={placeholder}
        value={value}
        onChange={handleChange}
        className="pl-8 pr-8"
      />
      {value && (
        <Button
          type="button"
          variant="ghost"
          size="icon"
          tabIndex={-1}
          aria-label={t('ui.text.clearSearch')}
          onClick={handleClear}
          className="absolute right-0.5 top-1/2 -translate-y-1/2 size-7"
        >
          <X className="size-3.5" />
        </Button>
      )}
    </div>
  )
}
