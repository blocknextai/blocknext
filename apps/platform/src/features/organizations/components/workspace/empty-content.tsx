import { GitPullRequestClosed } from 'lucide-react'
import { useTranslation } from 'react-i18next'

const EmptyContent = () => {
  const { t } = useTranslation()

  return (
    <div className="flex justify-center items-center w-full h-1/3 text-muted-foreground gap-2 p-10">
      <GitPullRequestClosed /> {t('ui.text.noFlows')}
    </div>
  )
}
EmptyContent.displayName = 'EmptyContent'

export { EmptyContent }
