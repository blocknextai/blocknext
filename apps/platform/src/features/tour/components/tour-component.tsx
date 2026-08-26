import { useState, useEffect, useMemo, useRef } from 'react'
import { useTranslation } from 'react-i18next'
import TourOverlay from '@/features/tour/components/tour-overlay'
import type { HighlightPosition } from '@/features/tour/components/tour-overlay'
import { markOnboardingTourSeen } from '@/features/tour/storage'
import {
  DEMO_VEO_NODE_ID,
  DEMO_VIDEO_REFERENCE,
  MISSING_DEMO_NODE_CANVAS_ID,
  MISSING_DEMO_NODE_ID,
  MISSING_DEMO_NODE_POSITION,
} from '@/features/tour/demo-flow'
import { usePlatformFeatures } from '@/features/platform/hooks/use-platform-features'

const TARGET_POLL_INTERVAL = 60
const TARGET_WAIT_TIMEOUT = 900
const OPTIONAL_TARGET_WAIT_TIMEOUT = 240
const COMPLETION_POLL_INTERVAL = 350
const LOST_TARGET_TIMEOUT = 3000
const CANVAS_TARGET = /^\[data-id="([^"]+)"\]$/
const STEP_TRANSITION_DELAY = 200

const EMPTY_HIGHLIGHT: HighlightPosition = {
  top: 0,
  left: 0,
  width: 0,
  height: 0,
}

const canvasBounds = () => {
  const pane = document.querySelector('.react-flow')
  if (!pane) {
    return null
  }

  const rect = pane.getBoundingClientRect()
  return {
    top: rect.top,
    left: rect.left,
    right: rect.right,
    bottom: rect.bottom,
  }
}

const activate = (element: HTMLElement) => {
  const options = { bubbles: true, cancelable: true, button: 0 }
  element.dispatchEvent(new MouseEvent('mousedown', options))
  element.dispatchEvent(new MouseEvent('mouseup', options))
  element.click()
}

const writeControlledValue = (element: HTMLElement, value: string) => {
  const prototype =
    element instanceof HTMLTextAreaElement
      ? HTMLTextAreaElement.prototype
      : HTMLInputElement.prototype
  const setter = Object.getOwnPropertyDescriptor(prototype, 'value')?.set
  setter?.call(element, value)
  element.dispatchEvent(new Event('input', { bubbles: true }))
}

const measure = (element: Element): HighlightPosition => {
  const rect = element.getBoundingClientRect()
  return {
    top: rect.top,
    left: rect.left,
    width: rect.width,
    height: rect.height,
  }
}

const TourComponent = ({
  steps = [],
  isOpen = false,
  onClose = () => {},
  onStepChange = () => {},
  onComplete = () => {},
  setSidebarOpen,
  addNodeById,
  selectGroupForNode,
  connectNodes,
  focusCanvas,
  openNodeSettings,
  getFlowState,
  setChatOpen,
  onNextStep,
}) => {
  const [currentStep, setCurrentStep] = useState(0)
  const [highlightPosition, setHighlightPosition] =
    useState<HighlightPosition>(EMPTY_HIGHLIGHT)
  const [targetElement, setTargetElement] = useState<HTMLElement | null>(null)
  const advancingRef = useRef(false)
  const selfActingRef = useRef(false)
  const alreadyCompleteRef = useRef(false)

  const goToStep = (index: number) => {
    setCurrentStep(index)
    onStepChange(index, steps[index])
  }

  const handleClose = () => {
    setCurrentStep(0)
    onClose()
  }

  const finish = () => {
    handleClose()
    onComplete()
  }

  const applyStepSideEffects = (step: any) => {
    if (typeof step?.chat === 'boolean' && typeof setChatOpen === 'function') {
      setChatOpen(step.chat)
    }

    if (step?.opensSidebar) {
      setSidebarOpen(true)
    }

    if (step?.opensNodeGroup && typeof selectGroupForNode === 'function') {
      selectGroupForNode(step.opensNodeGroup)
    }

    if (step?.clicksTarget && targetElement) {
      activate(targetElement)
    }

    if (step?.clicks) {
      const element = document.querySelector(step.clicks) as HTMLElement | null
      if (element) {
        activate(element)
      }
    }

    if (step?.opensNodeSettings && typeof openNodeSettings === 'function') {
      openNodeSettings(step.opensNodeSettings)
    }

    if (step?.connectsFromNodeId && typeof connectNodes === 'function') {
      const nodes =
        typeof getFlowState === 'function' ? getFlowState().nodes : []
      const targetId =
        nodes.find((node: any) => node.nodeId === step.connectsToNodeType)
          ?.id ?? step.connectsToNodeId
      connectNodes(step.connectsFromNodeId, targetId)
    }

    if (step?.fillsField) {
      const input = document.querySelector(
        step.fillsField.selector,
      ) as HTMLElement | null
      if (input) {
        writeControlledValue(input, step.fillsField.value)
      }
    }

    if (step?.addsSampleNode && typeof addNodeById === 'function') {
      const nodeId =
        targetElement?.getAttribute('data-node-id') ?? step.fallbackNodeId
      if (nodeId) {
        addNodeById(nodeId, step.sampleNodePosition)
      }
    }
  }

  const advance = () => {
    if (advancingRef.current) {
      return
    }
    advancingRef.current = true

    if (currentStep >= steps.length - 1) {
      finish()
      return
    }

    setTimeout(() => goToStep(currentStep + 1), STEP_TRANSITION_DELAY)
  }

  const nextStep = () => {
    if (advancingRef.current) {
      return
    }
    selfActingRef.current = true
    applyStepSideEffects(steps[currentStep])
    setTimeout(() => {
      selfActingRef.current = false
    }, STEP_TRANSITION_DELAY * 2)
    advance()
  }

  const prevStep = () => {
    if (currentStep > 0) {
      goToStep(currentStep - 1)
    }
  }

  useEffect(() => {
    if (typeof onNextStep === 'function') {
      onNextStep(advance)
    }
  })

  useEffect(() => {
    if (!isOpen) {
      return
    }

    const handleClick = (event: MouseEvent) => {
      if (selfActingRef.current || advancingRef.current) {
        return
      }

      const clicked = event.target as Element | null
      if (!clicked?.closest || clicked.closest('[data-tour-tooltip]')) {
        return
      }

      const index = steps.findIndex(
        (candidate, position) =>
          position !== currentStep &&
          candidate.target &&
          clicked.closest(candidate.target),
      )

      if (index >= 0) {
        advancingRef.current = true
        goToStep(index)
      }
    }

    document.addEventListener('click', handleClick, true)
    return () => document.removeEventListener('click', handleClick, true)
  }, [isOpen, currentStep, steps])

  useEffect(() => {
    if (isOpen && typeof focusCanvas === 'function') {
      focusCanvas()
    }
  }, [isOpen, focusCanvas])

  useEffect(() => {
    const isComplete = steps[currentStep]?.isComplete
    if (!isOpen || !targetElement || typeof isComplete !== 'function') {
      return
    }

    const flowState = () =>
      typeof getFlowState === 'function'
        ? getFlowState()
        : { nodes: [], edges: [] }

    if (alreadyCompleteRef.current) {
      return
    }

    const timer = window.setInterval(() => {
      if (isComplete(flowState())) {
        window.clearInterval(timer)
        advance()
      }
    }, COMPLETION_POLL_INTERVAL)

    return () => window.clearInterval(timer)
  }, [isOpen, currentStep, targetElement, getFlowState, steps])

  useEffect(() => {
    if (!isOpen || currentStep >= steps.length) {
      return
    }

    const step = steps[currentStep]
    const timeout = step.optional
      ? OPTIONAL_TARGET_WAIT_TIMEOUT
      : TARGET_WAIT_TIMEOUT

    let cancelled = false
    let waited = 0
    let lost = 0
    let frame = 0
    let tracked: HTMLElement | null = null

    setTargetElement(null)
    advancingRef.current = false
    alreadyCompleteRef.current =
      typeof step.isComplete === 'function' &&
      step.isComplete(
        typeof getFlowState === 'function'
          ? getFlowState()
          : { nodes: [], edges: [] },
      )

    const prepare = () => {
      if (!step.prepare) {
        return
      }

      if (step.prepare.sidebar) {
        setSidebarOpen(true)
      }
      if (step.prepare.nodeGroup && typeof selectGroupForNode === 'function') {
        selectGroupForNode(step.prepare.nodeGroup)
      }
      if (step.prepare.nodeSettings && typeof openNodeSettings === 'function') {
        openNodeSettings(step.prepare.nodeSettings)
      }
      if (step.prepare.activates) {
        setTimeout(() => {
          const element = document.querySelector(
            step.prepare.activates,
          ) as HTMLElement | null
          if (element) {
            activate(element)
          }
        }, TARGET_POLL_INTERVAL)
      }
    }

    const track = () => {
      if (cancelled) {
        return
      }

      if (!tracked || !tracked.isConnected) {
        tracked = null
        lost = 0
        setTargetElement(null)
        reacquire()
        return
      }

      const next = measure(tracked)
      setHighlightPosition((previous) =>
        previous.top === next.top &&
        previous.left === next.left &&
        previous.width === next.width &&
        previous.height === next.height
          ? previous
          : next,
      )
      frame = requestAnimationFrame(track)
    }

    const attach = (element: HTMLElement) => {
      tracked = element
      setHighlightPosition(measure(element))
      setTargetElement(element)
      frame = requestAnimationFrame(track)
    }

    const reacquire = () => {
      if (cancelled) {
        return
      }

      const element = document.querySelector(step.target) as HTMLElement | null
      if (!element) {
        lost += TARGET_POLL_INTERVAL
        if (lost >= LOST_TARGET_TIMEOUT) {
          if (currentStep >= steps.length - 1) {
            finish()
          } else {
            goToStep(currentStep + 1)
          }
          return
        }
        setTimeout(reacquire, TARGET_POLL_INTERVAL)
        return
      }

      attach(element)
    }

    const acquire = () => {
      if (cancelled) {
        return
      }

      const element = document.querySelector(step.target) as HTMLElement | null

      if (!element) {
        waited += TARGET_POLL_INTERVAL
        if (waited >= timeout) {
          if (currentStep >= steps.length - 1) {
            finish()
          } else {
            goToStep(currentStep + 1)
          }
          return
        }
        setTimeout(acquire, TARGET_POLL_INTERVAL)
        return
      }

      const canvasNode = step.target.match(CANVAS_TARGET)
      if (canvasNode && typeof focusCanvas === 'function') {
        focusCanvas(canvasNode[1])
      } else {
        element.scrollIntoView({ behavior: 'smooth', block: 'center' })
      }

      attach(element)
    }

    prepare()
    acquire()

    return () => {
      cancelled = true
      cancelAnimationFrame(frame)
    }
  }, [currentStep, isOpen, steps])

  if (!isOpen || steps.length === 0 || !targetElement) {
    return null
  }

  return (
    <TourOverlay
      highlightPosition={highlightPosition}
      bounds={canvasBounds()}
      currentStep={currentStep}
      steps={steps}
      onNext={nextStep}
      onPrev={prevStep}
      onClose={handleClose}
    />
  )
}

const TourComp = ({
  shouldStartTour = false,
  onTourSkip = () => {},
  sidebarOpen,
  setSidebarOpen,
  addNodeById,
  selectGroupForNode,
  connectNodes,
  focusCanvas,
  openNodeSettings,
  getFlowState,
  setChatOpen,
  onNextStep,
}) => {
  const { t } = useTranslation()
  const { workflowsGenerationEnabled, functionCallingEnabled } =
    usePlatformFeatures()
  const [isTourOpen, setIsTourOpen] = useState(false)

  const tourSteps = useMemo(
    () =>
      [
        {
          target: '[data-tour="ai-flow-builder"]',
          title: t('ui.text.generation.aiFlowBuilder'),
          description: t('ui.text.aiFlowBuilderDescription'),
          position: 'bottom' as const,
          requires: 'generation',
          chat: true,
          isComplete: () =>
            document.querySelector('[data-tour="ai-chat"]') !== null,
        },
        {
          target: '[data-tour="ai-chat"]',
          title: t('ui.text.aiChatPanel'),
          description: t('ui.text.aiChatPanelDescription'),
          position: 'left' as const,
          requires: 'generation',
          chat: false,
        },
        {
          target: '[data-tour="sidebar-toggle"]',
          title: t('ui.text.sidebarMenu'),
          description: t('ui.text.sidebarDescription'),
          position: 'right' as const,
          opensSidebar: true,
          isComplete: () =>
            (document.querySelector('#flowSideMenu') as HTMLInputElement | null)
              ?.checked === true,
        },
        {
          target: '[data-tour="category-grid"]',
          title: t('ui.text.nodeGroups'),
          description: t('ui.text.nodeGroupsDescription'),
          position: 'right' as const,
          opensNodeGroup: MISSING_DEMO_NODE_ID,
          isComplete: () => document.querySelector('[data-node-id]') !== null,
          prepare: { sidebar: true },
        },
        {
          target: `[data-node-id="${MISSING_DEMO_NODE_ID}"]`,
          title: t('ui.text.dragAndDropNodes'),
          description: t('ui.text.dragMissingNodeDescription'),
          position: 'right' as const,
          addsSampleNode: true,
          fallbackNodeId: MISSING_DEMO_NODE_ID,
          sampleNodePosition: MISSING_DEMO_NODE_POSITION,
          isComplete: ({ nodes }) =>
            nodes.some((node: any) => node.nodeId === MISSING_DEMO_NODE_ID),
          prepare: { sidebar: true, nodeGroup: MISSING_DEMO_NODE_ID },
        },
        {
          target: `[data-id="${MISSING_DEMO_NODE_CANVAS_ID}"]`,
          title: t('ui.text.connectMissingNode'),
          description: t('ui.text.connectMissingNodeDescription'),
          position: 'right' as const,
          connectsFromNodeId: DEMO_VEO_NODE_ID,
          connectsToNodeId: MISSING_DEMO_NODE_CANVAS_ID,
          connectsToNodeType: MISSING_DEMO_NODE_ID,
          isComplete: ({ nodes, edges }) => {
            const youtube = nodes.find(
              (node: any) => node.nodeId === MISSING_DEMO_NODE_ID,
            )
            return (
              !!youtube &&
              edges.some(
                (edge: any) =>
                  edge.source === DEMO_VEO_NODE_ID &&
                  edge.target === youtube.id,
              )
            )
          },
        },
        {
          target: `[data-id="${MISSING_DEMO_NODE_CANVAS_ID}"]`,
          title: t('ui.text.configureMissingNode'),
          description: t('ui.text.configureMissingNodeDescription'),
          position: 'right' as const,
          opensNodeSettings: MISSING_DEMO_NODE_CANVAS_ID,
          isComplete: () =>
            document.querySelector('[data-tour="node-settings-panel"]') !==
            null,
        },
        {
          target: '[data-tour="node-settings-basic"]',
          title: t('ui.text.manageWithNaturalLanguage'),
          description: t('ui.text.naturalLanguageTourDescription'),
          position: 'left' as const,
          requires: 'functionCalling',
          clicks: '[data-tour="node-settings-advanced"]',
          prepare: {
            nodeSettings: MISSING_DEMO_NODE_CANVAS_ID,
            activates: '[data-tour="node-settings-basic"]',
          },
          isComplete: () =>
            document.querySelector('[data-field-key="videoUrl"]') !== null,
        },
        {
          target: '[data-tour="node-settings-advanced"]',
          title: t('ui.text.advanced'),
          description: t('ui.text.advancedTabTourDescription'),
          position: 'left' as const,
          requires: 'functionCalling',
          clicksTarget: true,
          prepare: {
            nodeSettings: MISSING_DEMO_NODE_CANVAS_ID,
            activates: '[data-tour="node-settings-advanced"]',
          },
          isComplete: () =>
            document.querySelector('[data-field-key="videoUrl"]') !== null,
        },
        {
          target: '[data-field-key="videoUrl"] [data-tour="data-source"]',
          title: t('ui.text.setVideoUrl'),
          description: t('ui.text.setVideoUrlDescription'),
          position: 'left' as const,
          isComplete: () =>
            (
              document.querySelector(
                '[data-field-key="videoUrl"] input, [data-field-key="videoUrl"] textarea',
              ) as HTMLInputElement | null
            )?.value?.includes(`$veo_${DEMO_VEO_NODE_ID}.`) === true,
          prepare: {
            nodeSettings: MISSING_DEMO_NODE_CANVAS_ID,
            activates: '[data-tour="node-settings-advanced"]',
          },
          fillsField: {
            selector:
              '[data-field-key="videoUrl"] input, [data-field-key="videoUrl"] textarea',
            value: DEMO_VIDEO_REFERENCE,
          },
        },
        {
          target: '[data-tour="node-settings-apply"]',
          title: t('ui.text.applyNodeSettings'),
          description: t('ui.text.applyNodeSettingsDescription'),
          position: 'left' as const,
          clicksTarget: true,
          prepare: {
            nodeSettings: MISSING_DEMO_NODE_CANVAS_ID,
            activates: '[data-tour="node-settings-advanced"]',
          },
          isComplete: ({ nodes }) =>
            nodes
              .find((node: any) => node.nodeId === MISSING_DEMO_NODE_ID)
              ?.parameters?.videoUrl?.includes(`$veo_${DEMO_VEO_NODE_ID}.`) ===
            true,
        },
        {
          target: '[data-tour="templates-marketplace"]',
          title: t('ui.text.templatesAndMarketplace'),
          description: t('ui.text.templatesDescription'),
          position: 'right' as const,
          optional: true,
        },
        {
          target: '[data-tour="flow-view-controls"]',
          title: t('ui.text.flowViewControls'),
          description: t('ui.text.flowViewControlsDescription'),
          position: 'top' as const,
        },
        {
          target: '[data-tour="flow-toolbar"]',
          title: t('ui.text.flowToolbar'),
          description: t('ui.text.flowToolbarDescription'),
          position: 'top' as const,
        },
        {
          target: '[data-id="0"]',
          title: t('ui.text.tourCongratulations'),
          description: t('ui.text.tourCongratulationsDescription'),
          position: 'top' as const,
        },
      ].filter((step) => {
        if (step.requires === 'generation') {
          return workflowsGenerationEnabled
        }
        if (step.requires === 'functionCalling') {
          return functionCallingEnabled
        }
        return true
      }),
    [t, workflowsGenerationEnabled, functionCallingEnabled],
  )

  useEffect(() => {
    if (shouldStartTour) {
      setTimeout(() => setIsTourOpen(true), 500)
    }
  }, [shouldStartTour])

  const handleTourComplete = () => {
    markOnboardingTourSeen()
    setIsTourOpen(false)
  }

  const handleTourClose = () => {
    markOnboardingTourSeen()
    setIsTourOpen(false)
    onTourSkip()
  }

  return (
    <TourComponent
      steps={tourSteps}
      isOpen={isTourOpen}
      onClose={handleTourClose}
      onComplete={handleTourComplete}
      sidebarOpen={sidebarOpen}
      setSidebarOpen={setSidebarOpen}
      addNodeById={addNodeById}
      selectGroupForNode={selectGroupForNode}
      connectNodes={connectNodes}
      focusCanvas={focusCanvas}
      openNodeSettings={openNodeSettings}
      getFlowState={getFlowState}
      setChatOpen={setChatOpen}
      onNextStep={onNextStep}
    />
  )
}

TourComponent.displayName = 'TourComponent'
TourComp.displayName = 'TourComp'

export default TourComp
