import { DebouncedTextarea } from '@/components/shared/debounced-textarea'
import { useTranslation } from 'react-i18next'

export const default_runner = (props) => {
  const { t } = useTranslation()
  return (
    <DebouncedTextarea
      {...props}
      placeholder={t('ui.text.provideDescription')}
    />
  )
}
