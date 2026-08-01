import { getCategoryPrefs } from '@/lib/flow-categories'

const FlowDragPreview = ({ dragging, previewData, previewPos }) => {
  if (!dragging || !previewData) return null

  const pref = getCategoryPrefs(previewData.category)
  const Icon = pref.icon

  return (
    <div
      className="size-10 rounded-lg bg-background fixed pointer-events-none z-100"
      style={{
        left: previewPos.x,
        top: previewPos.y,
        color: pref.color,
      }}
    >
      <div className="bg-current/10 w-full h-full flex items-center justify-center rounded-lg p-1">
        <Icon className="size-4" />
      </div>
    </div>
  )
}
FlowDragPreview.displayName = 'FlowDragPreview'

export { FlowDragPreview }
