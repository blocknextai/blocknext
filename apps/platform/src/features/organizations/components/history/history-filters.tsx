import { useTranslation } from 'react-i18next'

const HistoryFilters = () => {
  const { t } = useTranslation()

  return (
    <div className="flex flex-col gap-2 shrink-0">
      <span className="heading-md">{t('ui.text.runHistory')}</span>
      <span className="text-muted-foreground text-sm">
        {t('ui.text.runHistoryDescription')}
      </span>
    </div>
  )
}

HistoryFilters.displayName = 'HistoryFilters'

export { HistoryFilters }
