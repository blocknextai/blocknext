import FlowCanvas from '@/features/flow-editor/components/flow-canvas'
import { useSearchParams } from 'react-router'
import { useTranslation } from 'react-i18next'
import { ONBOARDING_DEMO, onboardingDemoFlow } from '@/features/tour/demo-flow'

function FlowsCreatePage() {
  const { t } = useTranslation()
  const [searchParams] = useSearchParams()
  const aiMode = searchParams.get('ai') === 'true'
  const demoFlow =
    searchParams.get('demo') === ONBOARDING_DEMO
      ? { ...onboardingDemoFlow, title: t('ui.text.demoFlowTitle') }
      : null

  return (
    <div className="w-full h-full flex-1 flex relative">
      <FlowCanvas initialFlow={demoFlow} defaultChatOpen={aiMode} />
    </div>
  )
}

export default FlowsCreatePage
