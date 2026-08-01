import { Workflow } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'

const EmptyWorkspace = ({ organizationId }) => {
  const { t } = useTranslation()

  return (
    <div className="w-full h-full flex flex-col gap-4 p-6 pt-3 items-center justify-center">
      <div className="flex flex-col items-center gap-4">
        <h2 className="text-2xl font-semibold">
          {t('ui.text.emptyWorkflowTitle')}
        </h2>
        <p
          className="align-center text-center"
          dangerouslySetInnerHTML={{
            __html: t('ui.text.emptyWorkflow'),
          }}
        ></p>
        <div className="flex gap-2">
          <Link to={`/organizations/${organizationId}/create`}>
            <div className="border flex flex-col gap-2 bg-card/30 rounded-md p-4 items-center justify-center align-center cursor-pointer hover:bg-card/60 transition">
              <Workflow className="h-6 w-6" />
              <span>{t('ui.text.createFlow')}</span>
            </div>
          </Link>
        </div>
      </div>
    </div>
  )
}
EmptyWorkspace.displayName = 'EmptyWorkspace'
export { EmptyWorkspace }
