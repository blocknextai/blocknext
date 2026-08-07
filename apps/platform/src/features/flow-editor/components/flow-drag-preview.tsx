import { getCategoryPrefs } from '@/lib/flow-categories'

const FlowDragPreview = ({ dragging, previewData, previewPos }) => {
  if (!dragging || !previewData) {
    return null
  }

  const pref = getCategoryPrefs(previewData.category)
  const Icon = pref.icon

  return (
    <div
      className="bg-card border-border pointer-events-none fixed z-100 size-10 rounded-lg border"
      style={{
        left: previewPos.x,
        top: previewPos.y,
      }}
    >
      <div className="flex h-full w-full items-center justify-center rounded-lg p-1">
        <Icon className="size-4" />
      </div>
    </div>
  )
}
FlowDragPreview.displayName = 'FlowDragPreview'

export { FlowDragPreview }
