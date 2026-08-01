import FlowCanvas from '@/features/flow-editor/components/flow-canvas'
import { useSearchParams } from 'react-router'

function FlowsCreatePage() {
  const [searchParams] = useSearchParams()
  const aiMode = searchParams.get('ai') === 'true'

  return (
    <div className="w-full h-full flex-1 flex relative">
      <FlowCanvas defaultChatOpen={aiMode} />
    </div>
  )
}

export default FlowsCreatePage
