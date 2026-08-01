import {
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from '@/components/ui/tooltip'
import { Info } from 'lucide-react'
import { useTranslation } from 'react-i18next'

const MemberRoleInfo = () => {
  const { t } = useTranslation()

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <Info />
      </TooltipTrigger>
      <TooltipContent side="bottom" className="flex flex-col">
        <p className="text-sm mb-1">{t('ui.text.rolesInfo')}</p>
        <ul className="pl-4 space-y-1 list-disc text-sm">
          <li>
            <b>{t('ui.text.roleOwner')}</b>: {t('ui.text.roleOwnerDesc')}
          </li>
          <li>
            <b>{t('ui.text.roleAdmin')}</b>: {t('ui.text.roleAdminDesc')}
          </li>
          <li>
            <b>{t('ui.text.roleEditor')}</b>: {t('ui.text.roleEditorDesc')}
          </li>
          <li>
            <b>{t('ui.text.roleViewer')}</b>: {t('ui.text.roleViewerDesc')}
          </li>
        </ul>
      </TooltipContent>
    </Tooltip>
  )
}

MemberRoleInfo.displayName = 'MemberRoleInfo'

export { MemberRoleInfo }
