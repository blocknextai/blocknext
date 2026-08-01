import FlowCanvas from '@/features/flow-editor/components/flow-canvas'
import { useEffect, useState } from 'react'
import { useParams } from 'react-router'
import { Loading } from '@/components/shared/loading'
import { workflowsService } from '@/features/workflows'

function FlowsDetailPage() {
  const { organizationId, id } = useParams()
  const [loading, setLoading] = useState(true)
  const [flow, setFlow] = useState()

  const getFlow = async () => {
    const response = await workflowsService.getById(organizationId, id)
    setFlow(response.data)
    setLoading(false)
  }

  useEffect(() => {
    getFlow()
  }, [])

  if (loading) {
    return <Loading />
  }

  return (
    <div className="w-full h-full flex-1 flex">
      <FlowCanvas initialFlow={flow} />
    </div>
  )
}

export default FlowsDetailPage
