import { useState, useRef, useEffect, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import {
  AnnotationNode,
  AnnotationNodeContent,
} from '@/features/flow-editor/nodes/annotation-node'
import { useFlowSetNodes } from '@/features/flow-editor/contexts/flow-nodes-context'

const Annotation = ({ id, data }) => {
  const { t } = useTranslation()
  const setNodes = useFlowSetNodes()
  const [isEditing, setIsEditing] = useState(false)
  const inputRef = useRef<HTMLTextAreaElement>(null)

  const note = data?.note ?? ''

  useEffect(() => {
    if (isEditing && inputRef.current) {
      inputRef.current.focus()
    }
  }, [isEditing])

  const updateNote = useCallback(
    (value: string) => {
      setNodes((nodes) =>
        nodes.map((n) =>
          n.id === id ? { ...n, data: { ...n.data, note: value } } : n,
        ),
      )
    },
    [id, setNodes],
  )

  const handleInputKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Escape') {
      setIsEditing(false)
    }
  }

  return (
    <AnnotationNode>
      <AnnotationNodeContent
        onClick={() => setIsEditing(true)}
        className="w-full cursor-text"
      >
        {isEditing ? (
          <textarea
            ref={inputRef}
            value={note}
            onChange={(e) => updateNote(e.target.value)}
            onBlur={() => setIsEditing(false)}
            onKeyDown={handleInputKeyDown}
            className="bg-card border-border text-secondary-foreground min-h-20 w-full resize-none rounded-md border px-2 py-1.5 text-sm outline-none"
          />
        ) : (
          <span className={note ? undefined : 'text-muted-foreground'}>
            {note || t('ui.text.annotationPlaceholder')}
          </span>
        )}
      </AnnotationNodeContent>
    </AnnotationNode>
  )
}
Annotation.displayName = 'Annotation'

export { Annotation }
