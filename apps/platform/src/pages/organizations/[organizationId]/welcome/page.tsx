import FlowCanvas from '@/features/flow-editor/components/flow-canvas'

const flow = {
  nodes: [
    {
      id: '0',
      position: { x: 0, y: 0 },
      data: {
        id: '0',
        title: 'ui.text.starterNode',
        description: 'ui.text.starterNodeDescription',
        hide_input: true,
        hide_left: true,
        tags: ['starter'],
        category: 'system',
      },
      type: 'starter',
      nodeId: 'starter',
      deletable: false,
    },
  ],
  edges: [],
}

const WelcomePage = () => {
  return (
    <div className="w-full h-full flex-1 flex">
      <FlowCanvas initialFlow={flow} defaultSidebarOpen={false} />
    </div>
  )
}

export default WelcomePage
