import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import TourOverlay from '@/features/tour/components/tour-overlay'
import type { HighlightPosition } from '@/features/tour/components/tour-overlay'

const TourComponent = ({
  steps = [],
  isOpen = false,
  onClose = () => {},
  onStepChange = () => {},
  onComplete = () => {},
  autoProgressOnClick = false,
  setSidebarOpen,
  onNextStep,
}) => {
  const [currentStep, setCurrentStep] = useState(0)
  const [highlightPosition, setHighlightPosition] = useState<HighlightPosition>(
    {
      top: 0,
      left: 0,
      width: 0,
      height: 0,
    },
  )
  const [originalStyles, setOriginalStyles] = useState(new Map())

  const calculatePosition = (element: Element): HighlightPosition => {
    const rect = element.getBoundingClientRect()
    return {
      top: rect.top,
      left: rect.left,
      width: rect.width,
      height: rect.height,
    }
  }

  const updatePosition = () => {
    if (isOpen && currentStep < steps.length) {
      const targetElement = document.querySelector(steps[currentStep].target)
      if (targetElement) {
        const position = calculatePosition(targetElement)
        setHighlightPosition(position)
      }
    }
  }

  const simulateStepAction = (stepIndex: number) => {
    const step = steps[stepIndex]

    if (step?.target === '[data-tour="sidebar-toggle"]') {
      setSidebarOpen(true)
      return
    }

    if (step?.target === '[data-tour="category-grid"]') {
      setTimeout(() => {
        const firstCategory = document.querySelector(
          '[data-tour="category-grid"] > div:first-child, [data-tour="category-grid"] .category-item:first-child',
        )
        if (firstCategory) {
          ;(firstCategory as HTMLElement).click()
        }
      }, 200)
    }
  }

  const nextStep = () => {
    if (currentStep < steps.length - 1) {
      simulateStepAction(currentStep)

      const newStep = currentStep + 1
      setTimeout(() => {
        setCurrentStep(newStep)
        onStepChange(newStep, steps[newStep])
      }, 200)
    } else {
      handleClose()
      onComplete()
    }
  }

  const prevStep = () => {
    if (currentStep > 0) {
      const newStep = currentStep - 1
      setCurrentStep(newStep)
      onStepChange(newStep, steps[newStep])
    }
  }

  const handleClose = () => {
    restoreAllStyles()
    setCurrentStep(0)
    onClose()
  }

  const restoreOriginalStyle = (element: HTMLElement) => {
    const original = originalStyles.get(element)
    if (original) {
      element.style.position = original.position
      element.style.zIndex = original.zIndex
      element.style.transform = original.transform
    }
  }
  useEffect(() => {
    if (typeof onNextStep === 'function') {
      onNextStep(nextStep)
    }
  }, [onNextStep])
  const restoreAllStyles = () => {
    originalStyles.forEach((_original, element) => {
      restoreOriginalStyle(element)
    })
    setOriginalStyles(new Map())
  }

  useEffect(() => {
    if (!autoProgressOnClick || !isOpen || currentStep >= steps.length) return

    const currentStepData = steps[currentStep]
    let targetElement = document.querySelector(currentStepData.target)

    if (currentStepData.target === '[data-tour="sidebar-toggle"]') {
      const iconElement = targetElement
      const labelElement = iconElement?.closest('label[for="flowSideMenu"]')
      targetElement = labelElement || targetElement
    }

    if (targetElement) {
      const eventType = currentStepData.triggerEvent || 'click'

      const handleTargetInteraction = (_e: Event) => {
        if (currentStepData.autoProgress !== false) {
          setTimeout(() => {
            nextStep()
          }, 200)
        }
      }

      targetElement.addEventListener(eventType, handleTargetInteraction, {
        capture: true,
      })

      return () => {
        targetElement.removeEventListener(eventType, handleTargetInteraction, {
          capture: true,
        })
      }
    }
  }, [autoProgressOnClick, isOpen, currentStep, steps])

  useEffect(() => {
    if (isOpen && currentStep < steps.length) {
      const currentStepData = steps[currentStep]

      if (currentStepData.showOnlyAfterCategorySelected) {
        const nodeItem = document.querySelector('[data-tour="node-item"]')
        if (!nodeItem) {
          setTimeout(() => nextStep(), 100)
          return
        }
      }

      const targetElement = document.querySelector(
        currentStepData.target,
      ) as HTMLElement | null
      if (targetElement) {
        setTimeout(() => {
          const position = calculatePosition(targetElement)
          setHighlightPosition(position)
        }, 100)

        if (targetElement.closest('.react-flow__node')) {
          targetElement.style.zIndex = '10001'
        } else {
          targetElement.style.position = 'relative'
          targetElement.style.zIndex = '10001'
        }

        targetElement.scrollIntoView({
          behavior: 'smooth',
          block: 'center',
        })
      }
    } else if (!isOpen) {
      const allTourElements = document.querySelectorAll(
        '[data-tour], .react-flow__node',
      )
      allTourElements.forEach((el: Element) => {
        ;(el as HTMLElement).style.zIndex = ''
        ;(el as HTMLElement).style.position = ''
      })
    }
  }, [currentStep, isOpen, steps])

  useEffect(() => {
    if (isOpen) {
      const handleUpdate = () => updatePosition()
      window.addEventListener('scroll', handleUpdate)
      window.addEventListener('resize', handleUpdate)
      return () => {
        window.removeEventListener('scroll', handleUpdate)
        window.removeEventListener('resize', handleUpdate)
      }
    }
  }, [isOpen, currentStep])

  if (!isOpen || steps.length === 0) return null

  return (
    <TourOverlay
      highlightPosition={highlightPosition}
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
  onNextStep,
}) => {
  const { t } = useTranslation('ui')
  const [isTourOpen, setIsTourOpen] = useState(false)

  const tourSteps = [
    {
      target: '[data-tour="ai-flow-builder"]',
      title: t('ui.text.generation.aiFlowBuilder'),
      description: t('ui.text.aiFlowBuilderDescription'),
      position: 'bottom' as const,
      autoProgress: false,
    },
    {
      target: '[data-tour="sidebar-toggle"]',
      title: t('ui.text.sidebarMenu'),
      description: t('ui.text.sidebarDescription'),
      position: 'right' as const,
      autoProgress: true,
    },
    {
      target: '[data-tour="category-grid"]',
      title: t('ui.text.nodeCategories'),
      description: t('ui.text.nodeCategoriesDescription'),
      position: 'right' as const,
      autoProgress: true,
    },
    {
      target: '[data-tour="node-item"]:first-child',
      title: t('ui.text.dragAndDropNodes'),
      description: t('ui.text.dragAndDropDescription'),
      position: 'right' as const,
      autoProgress: false,
      showOnlyAfterCategorySelected: true,
    },
    {
      target: '[data-id="0"]',
      title: t('ui.text.clickNodeForSettings'),
      description: t('ui.text.clickNodeForSettingsDescription'),
      position: 'right' as const,
      autoProgress: true,
    },
    {
      target: '[data-tour="flow-toolbar"]',
      title: t('ui.text.flowToolbar'),
      description: t('ui.text.flowToolbarDescription'),
      position: 'top' as const,
      autoProgress: false,
    },
    {
      target: '[data-tour="flow-toolbar"]',
      title: t('ui.text.tourCongratulations'),
      description: t('ui.text.tourCongratulationsDescription'),
      position: 'top' as const,
      autoProgress: false,
    },
  ]

  useEffect(() => {
    if (shouldStartTour) {
      setTimeout(() => {
        setIsTourOpen(true)
      }, 500)
    }
  }, [shouldStartTour])
  const handleTourComplete = () => {
    localStorage.setItem('hasVisited', 'true')
    setIsTourOpen(false)
  }

  const handleTourClose = () => {
    localStorage.setItem('hasVisited', 'true')
    setIsTourOpen(false)
    onTourSkip()
  }

  return (
    <TourComponent
      steps={tourSteps}
      isOpen={isTourOpen}
      autoProgressOnClick={true}
      onClose={handleTourClose}
      onComplete={handleTourComplete}
      sidebarOpen={sidebarOpen}
      setSidebarOpen={setSidebarOpen}
      onNextStep={onNextStep}
    />
  )
}

TourComponent.displayName = 'TourComponent'
TourComp.displayName = 'TourComp'

export default TourComp
