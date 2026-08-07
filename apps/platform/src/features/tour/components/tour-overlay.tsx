import { useTranslation } from 'react-i18next'

const getTooltipPosition = (steps, currentStep, highlightPosition) => {
  if (!steps[currentStep]) {
    return { position: 'fixed' as const, top: 0, left: 0 }
  }

  const { position = 'bottom' } = steps[currentStep]
  const { top, left, width, height } = highlightPosition

  const tooltipStyle: Record<string, any> = { position: 'fixed', zIndex: 10002 }

  switch (position) {
    case 'bottom':
      tooltipStyle.top = top + height + 15
      tooltipStyle.left = left + width / 2
      tooltipStyle.transform = 'translateX(-50%)'
      break
    case 'top':
      tooltipStyle.top = top - 15
      tooltipStyle.left = left + width / 2
      tooltipStyle.transform = 'translate(-50%, -100%)'
      break
    case 'right':
      tooltipStyle.top = top + height / 2
      tooltipStyle.left = left + width + 15
      tooltipStyle.transform = 'translateY(-50%)'
      break
    case 'left':
      tooltipStyle.top = top + height / 2
      tooltipStyle.left = left - 15
      tooltipStyle.transform = 'translate(-100%, -50%)'
      break
    default:
      tooltipStyle.top = top + height + 15
      tooltipStyle.left = left + width / 2
      tooltipStyle.transform = 'translateX(-50%)'
  }

  return tooltipStyle
}

const TourOverlay = ({
  highlightPosition,
  currentStep,
  steps,
  onNext,
  onPrev,
  onClose,
}) => {
  const { t } = useTranslation()
  const currentStepData = steps[currentStep]

  return (
    <div className="fixed inset-0 pointer-events-none" style={{ zIndex: 9999 }}>
      <div
        className="absolute transition-all duration-300 animate-pulse"
        style={{
          top: highlightPosition.top - 8,
          left: highlightPosition.left - 8,
          width: highlightPosition.width + 16,
          height: highlightPosition.height + 16,
          border: '3px solid var(--primary)',
          borderRadius: '12px',
          background: 'transparent',
        }}
      />

      <div
        className="absolute transition-all duration-300 animate-pulse"
        style={{
          top: highlightPosition.top - 12,
          left: highlightPosition.left - 12,
          width: highlightPosition.width + 24,
          height: highlightPosition.height + 24,
          borderRadius: '16px',
        }}
      />

      <div
        className="bg-card text-white rounded-lg shadow-xl p-4 max-w-sm border border-primary/30 pointer-events-auto"
        style={{
          ...getTooltipPosition(steps, currentStep, highlightPosition),
          zIndex: 10002,
        }}
      >
        {currentStepData?.title && (
          <h4 className="font-semibold text-white mb-2">
            {currentStepData.title}
          </h4>
        )}
        {currentStepData?.description && (
          <p className="text-gray-200 text-sm mb-4">
            {currentStepData.description}
          </p>
        )}

        <div className="flex justify-between items-center">
          <span className="text-xs text-gray-400">
            {currentStep + 1} / {steps.length}
          </span>
          <div className="flex gap-2">
            {currentStep > 0 && (
              <button
                onClick={onPrev}
                className="text-xs text-gray-400 hover:text-gray-200 transition-colors px-2 py-1 rounded cursor-pointer"
              >
                {t('ui.text.back')}
              </button>
            )}
            <button
              onClick={onClose}
              className="text-xs text-gray-400 hover:text-gray-200 transition-colors px-2 py-1 rounded cursor-pointer"
            >
              {t('ui.text.skip')}
            </button>
            <button
              onClick={onNext}
              className="bg-primary hover:bg-primary/80 text-primary-foreground text-xs px-3 py-1 rounded transition-colors cursor-pointer"
            >
              {currentStep === steps.length - 1
                ? t('ui.text.finish')
                : t('ui.text.next')}
            </button>
          </div>
        </div>

        <div
          className="absolute w-3 h-3 bg-gray-900 border-l border-t border-gray-700 transform rotate-45"
          style={{
            ...(currentStepData?.position === 'bottom' && {
              top: '-6px',
              left: '50%',
              transform: 'translateX(-50%) rotate(45deg)',
            }),
            ...(currentStepData?.position === 'top' && {
              bottom: '-6px',
              left: '50%',
              transform: 'translateX(-50%) rotate(-135deg)',
            }),
            ...(currentStepData?.position === 'right' && {
              left: '-6px',
              top: '50%',
              transform: 'translateY(-50%) rotate(135deg)',
            }),
            ...(currentStepData?.position === 'left' && {
              right: '-6px',
              top: '50%',
              transform: 'translateY(-50%) rotate(-45deg)',
            }),
          }}
        />
      </div>
    </div>
  )
}

TourOverlay.displayName = 'TourOverlay'

export default TourOverlay
