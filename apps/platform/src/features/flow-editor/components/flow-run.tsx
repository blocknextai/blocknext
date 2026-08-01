import { useTranslation } from 'react-i18next'
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion'
import * as Runner from '@/features/flow-editor/components/runner'

const FlowRun = ({ flow, setFlow }) => {
  const { t } = useTranslation()
  const updateNodeValue = (index, value) => {
    const nFlow = { ...flow }
    nFlow.nodes[index].runtimeInstruction = value
    setFlow(nFlow)
  }
  const renderSteps = () => {
    const steps = []
    for (let i = 0; i < flow.nodes.length; i++) {
      const node = flow.nodes[i]
      if (node.type === 'starter') continue
      const runnerKey = node.nodeId?.replaceAll('.', '_')
      if (runnerKey && Runner[runnerKey]) {
        node.runner = Runner[runnerKey]
      } else {
        node.runner = Runner['default_runner']
      }
      steps.push(
        <AccordionItem
          key={i}
          value={`item-${i}`}
          className="bg-input/30 pb-2 border-b-0! rounded-lg!"
        >
          <AccordionTrigger>
            <div className="flex items-start justify-between w-full space-y-1.5 px-3">
              <div className="flex flex-col gap-2">
                <span className="heading-md font-semibold">
                  {t(node.nodeId)}
                </span>
              </div>
            </div>
          </AccordionTrigger>
          <AccordionContent>
            <div className="flex flex-col gap-4 items-start justify-center">
              <div className="w-full flex-1/2 px-2">
                <node.runner
                  className="resize-none"
                  rows={100}
                  defaultValue={node.runtimeInstruction}
                  onChange={(value) => updateNodeValue(i, value)}
                />
              </div>
            </div>
          </AccordionContent>
        </AccordionItem>,
      )
    }
    return steps
  }

  return (
    <div className="w-full h-fit">
      <div className="flex flex-col gap-10 px-0">
        <Accordion
          type="single"
          defaultValue="item-1"
          collapsible
          className="flex flex-col w-full gap-4"
        >
          {renderSteps()}
        </Accordion>
      </div>
    </div>
  )
}

FlowRun.displayName = 'FlowRun'

export default FlowRun
