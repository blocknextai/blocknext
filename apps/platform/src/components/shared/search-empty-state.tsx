import { Search } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'

interface SearchEmptyStateProps {
  className?: string
}

/**
 * Shared empty state for when a search/filter returns no results.
 * Intentionally does NOT echo the typed query back to the user.
 */
export function SearchEmptyState({ className }: SearchEmptyStateProps) {
  const { t } = useTranslation()
  return (
    <div
      className={cn(
        'flex w-full h-full flex-col items-center justify-center text-center py-10',
        className,
      )}
    >
      <div className="flex flex-col items-center gap-3">
        <div className="p-3 rounded-full w-12 h-12 flex items-center justify-center">
          <Search className="size-6 text-muted-foreground/80" />
        </div>
        <div className="flex flex-col gap-2">
          <span className="text-md font-semibold">
            {t('ui.text.noSearchResults')}
          </span>
          <span className="text-muted-foreground text-xs max-w-xs">
            {t('ui.text.noSearchResultsDescription')}
          </span>
        </div>
      </div>
    </div>
  )
}

SearchEmptyState.displayName = 'SearchEmptyState'
