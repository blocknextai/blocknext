import { useLayoutEffect, useRef, useState } from 'react'
import { createPortal } from 'react-dom'
import { useTranslation } from 'react-i18next'

export interface HighlightPosition {
  top: number
  left: number
  width: number
  height: number
}

type Placement = 'top' | 'bottom' | 'left' | 'right'

interface Size {
  width: number
  height: number
}

interface Bounds {
  top: number
  left: number
  right: number
  bottom: number
}

interface Rect {
  top: number
  left: number
  right: number
  bottom: number
  centerX: number
  centerY: number
}

const VIEWPORT_MARGIN = 16
const HIGHLIGHT_PADDING = 8
const TOOLTIP_GAP = 14
const ARROW_INSET = 18

const OPPOSITE: Record<Placement, Placement> = {
  top: 'bottom',
  bottom: 'top',
  left: 'right',
  right: 'left',
}

const clamp = (value: number, min: number, max: number) =>
  Math.min(Math.max(value, min), Math.max(min, max))

const toHighlightRect = (position: HighlightPosition): Rect => ({
  top: position.top - HIGHLIGHT_PADDING,
  left: position.left - HIGHLIGHT_PADDING,
  right: position.left + position.width + HIGHLIGHT_PADDING,
  bottom: position.top + position.height + HIGHLIGHT_PADDING,
  centerX: position.left + position.width / 2,
  centerY: position.top + position.height / 2,
})

const availableSpace = (placement: Placement, rect: Rect, bounds: Bounds) => {
  switch (placement) {
    case 'top':
      return rect.top - bounds.top - VIEWPORT_MARGIN
    case 'bottom':
      return bounds.bottom - rect.bottom - VIEWPORT_MARGIN
    case 'left':
      return rect.left - bounds.left - VIEWPORT_MARGIN
    case 'right':
      return bounds.right - rect.right - VIEWPORT_MARGIN
  }
}

const requiredSpace = (placement: Placement, size: Size) =>
  placement === 'top' || placement === 'bottom'
    ? size.height + TOOLTIP_GAP
    : size.width + TOOLTIP_GAP

const resolvePlacement = (
  preferred: Placement,
  rect: Rect,
  size: Size,
  bounds: Bounds,
): Placement => {
  const candidates: Placement[] = [
    preferred,
    OPPOSITE[preferred],
    'bottom',
    'top',
    'right',
    'left',
  ]

  const surplus = (placement: Placement) =>
    availableSpace(placement, rect, bounds) - requiredSpace(placement, size)

  return (
    candidates.find((placement) => surplus(placement) >= 0) ??
    candidates.reduce((best, placement) =>
      surplus(placement) > surplus(best) ? placement : best,
    )
  )
}

const layoutFor = (
  placement: Placement,
  rect: Rect,
  size: Size,
  bounds: Bounds,
) => {
  const minLeft = bounds.left + VIEWPORT_MARGIN
  const minTop = bounds.top + VIEWPORT_MARGIN
  const maxLeft = bounds.right - size.width - VIEWPORT_MARGIN
  const maxTop = bounds.bottom - size.height - VIEWPORT_MARGIN

  const alignedLeft = clamp(rect.centerX - size.width / 2, minLeft, maxLeft)
  const alignedTop = clamp(rect.centerY - size.height / 2, minTop, maxTop)

  const anchors: Record<Placement, { top: number; left: number }> = {
    top: { top: rect.top - TOOLTIP_GAP - size.height, left: alignedLeft },
    bottom: { top: rect.bottom + TOOLTIP_GAP, left: alignedLeft },
    left: { top: alignedTop, left: rect.left - TOOLTIP_GAP - size.width },
    right: { top: alignedTop, left: rect.right + TOOLTIP_GAP },
  }

  const top = clamp(anchors[placement].top, minTop, maxTop)
  const left = clamp(anchors[placement].left, minLeft, maxLeft)

  const arrowOffset =
    placement === 'top' || placement === 'bottom'
      ? clamp(rect.centerX - left, ARROW_INSET, size.width - ARROW_INSET)
      : clamp(rect.centerY - top, ARROW_INSET, size.height - ARROW_INSET)

  return { top, left, placement, arrowOffset }
}

const arrowStyleFor = (placement: Placement, offset: number) => {
  switch (placement) {
    case 'bottom':
      return {
        top: -6,
        left: offset,
        transform: 'translateX(-50%) rotate(45deg)',
      }
    case 'top':
      return {
        bottom: -6,
        left: offset,
        transform: 'translateX(-50%) rotate(-135deg)',
      }
    case 'right':
      return {
        left: -6,
        top: offset,
        transform: 'translateY(-50%) rotate(135deg)',
      }
    case 'left':
      return {
        right: -6,
        top: offset,
        transform: 'translateY(-50%) rotate(-45deg)',
      }
  }
}

const TourOverlay = ({
  highlightPosition,
  bounds,
  currentStep,
  steps,
  onNext,
  onPrev,
  onClose,
}) => {
  const { t } = useTranslation()
  const currentStepData = steps[currentStep]
  const tooltipRef = useRef<HTMLDivElement>(null)
  const [size, setSize] = useState<Size | null>(null)

  useLayoutEffect(() => {
    const element = tooltipRef.current
    if (!element) {
      return
    }

    const { width, height } = element.getBoundingClientRect()
    setSize((previous) =>
      previous && previous.width === width && previous.height === height
        ? previous
        : { width, height },
    )
  })

  const area: Bounds = bounds ?? {
    top: 0,
    left: 0,
    right: window.innerWidth,
    bottom: window.innerHeight,
  }
  const rect = toHighlightRect(highlightPosition)
  const layout = size
    ? layoutFor(
        resolvePlacement(
          currentStepData?.position ?? 'bottom',
          rect,
          size,
          area,
        ),
        rect,
        size,
        area,
      )
    : null

  const ring = {
    top: clamp(rect.top, area.top, area.bottom),
    left: clamp(rect.left, area.left, area.right),
    right: clamp(rect.right, area.left, area.right),
    bottom: clamp(rect.bottom, area.top, area.bottom),
  }

  return createPortal(
    <div className="fixed inset-0 pointer-events-none" style={{ zIndex: 9999 }}>
      <div
        className="absolute transition-all duration-300 animate-pulse"
        style={{
          top: ring.top,
          left: ring.left,
          width: ring.right - ring.left,
          height: ring.bottom - ring.top,
          border: '3px solid var(--primary)',
          borderRadius: '12px',
          background: 'transparent',
        }}
      />

      <div
        ref={tooltipRef}
        data-tour-tooltip=""
        className="bg-card text-card-foreground rounded-lg shadow-xl p-4 w-80 max-w-[calc(100vw-2rem)] border border-primary/30 pointer-events-auto"
        style={{
          position: 'fixed',
          zIndex: 10002,
          top: layout?.top ?? 0,
          left: layout?.left ?? 0,
          visibility: layout ? 'visible' : 'hidden',
        }}
      >
        {currentStepData?.title && (
          <h4 className="font-semibold mb-2">{currentStepData.title}</h4>
        )}
        {currentStepData?.description && (
          <p className="text-muted-foreground text-sm mb-4">
            {currentStepData.description}
          </p>
        )}

        <div className="flex justify-between items-center">
          <span className="text-xs text-muted-foreground">
            {currentStep + 1} / {steps.length}
          </span>
          <div className="flex gap-2">
            {currentStep > 0 && (
              <button
                onClick={onPrev}
                className="text-xs text-muted-foreground hover:text-foreground transition-colors px-2 py-1 rounded cursor-pointer"
              >
                {t('ui.text.back')}
              </button>
            )}
            <button
              onClick={onClose}
              className="text-xs text-muted-foreground hover:text-foreground transition-colors px-2 py-1 rounded cursor-pointer"
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

        {layout && (
          <div
            className="absolute w-3 h-3 bg-card border-l border-t border-primary/30"
            style={arrowStyleFor(layout.placement, layout.arrowOffset)}
          />
        )}
      </div>
    </div>,
    document.body,
  )
}

TourOverlay.displayName = 'TourOverlay'

export default TourOverlay
