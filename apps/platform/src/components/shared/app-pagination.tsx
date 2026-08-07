import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination'

export interface PaginationInfo {
  total: number
  page?: number
  limit: number
  offset: number
  hasNext?: boolean
  hasPrev?: boolean
}

interface AppPaginationProps {
  pagination: PaginationInfo
  onPageChange: (page: number) => void
  showSummary?: boolean
  className?: string
}

function computePages(current: number, total: number): Array<number | string> {
  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1)
  }
  const pages: Array<number | string> = [1]
  if (current > 3) {
    pages.push('ellipsis-start')
  }
  const start = Math.max(2, current - 1)
  const end = Math.min(total - 1, current + 1)
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  if (current < total - 2) {
    pages.push('ellipsis-end')
  }
  pages.push(total)
  return pages
}

export function AppPagination({
  pagination,
  onPageChange,
  showSummary = true,
  className = '',
}: AppPaginationProps) {
  const { t } = useTranslation()
  const { total, limit, offset } = pagination
  const totalPages = Math.max(1, Math.ceil(total / limit))
  const currentPage = pagination.page ?? Math.floor(offset / limit) + 1
  const hasPrev = pagination.hasPrev ?? currentPage > 1
  const hasNext = pagination.hasNext ?? currentPage < totalPages

  const pageNumbers = useMemo(
    () => computePages(currentPage, totalPages),
    [currentPage, totalPages],
  )

  if (total <= 0) {
    return null
  }

  return (
    <div className={`mt-4 flex flex-col items-center gap-3 ${className}`}>
      {showSummary && (
        <div className="text-sm text-muted-foreground">
          {t('ui.text.showingResults', {
            from: offset + 1,
            to: Math.min(offset + limit, total),
            total,
          })}
        </div>
      )}
      <Pagination>
        <PaginationContent>
          <PaginationItem>
            <PaginationPrevious
              onClick={() => hasPrev && onPageChange(currentPage - 1)}
              className={
                !hasPrev ? 'pointer-events-none opacity-50' : 'cursor-pointer'
              }
            />
          </PaginationItem>
          {pageNumbers.map((p) => {
            if (typeof p === 'string' && p.startsWith('ellipsis')) {
              return (
                <PaginationItem key={p}>
                  <PaginationEllipsis />
                </PaginationItem>
              )
            }
            const isActive = currentPage === p
            return (
              <PaginationItem key={p}>
                <PaginationLink
                  onClick={() => !isActive && onPageChange(p as number)}
                  isActive={isActive}
                  className={
                    isActive
                      ? 'pointer-events-none opacity-50'
                      : 'cursor-pointer'
                  }
                >
                  {p}
                </PaginationLink>
              </PaginationItem>
            )
          })}
          <PaginationItem>
            <PaginationNext
              onClick={() => hasNext && onPageChange(currentPage + 1)}
              className={
                !hasNext ? 'pointer-events-none opacity-50' : 'cursor-pointer'
              }
            />
          </PaginationItem>
        </PaginationContent>
      </Pagination>
    </div>
  )
}
