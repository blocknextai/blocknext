import { useTranslation } from 'react-i18next'

const TriggerActions = () => {
  const { t } = useTranslation()
  return (
    <div className="flex flex-col gap-2 shrink-0">
      <span className="heading-md">{t('ui.text.triggers')}</span>
      <span className="text-muted-foreground text-sm">
        {t('ui.text.triggersDescription')}
      </span>
    </div>
  )
}

TriggerActions.displayName = 'TriggerActions'

export { TriggerActions }
