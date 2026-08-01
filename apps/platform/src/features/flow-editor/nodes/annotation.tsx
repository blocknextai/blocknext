import { useState, useRef, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import {
  AnnotationNode,
  AnnotationNodeContent,
} from '@/features/flow-editor/nodes/annotation-node'

const Annotation = () => {
  const { t } = useTranslation()
  const [content, setContent] = useState(t('ui.text.annotationPlaceholder'))
  const [isEditing, setIsEditing] = useState(false)
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    if (isEditing && inputRef.current) {
      inputRef.current.focus()
    }
  }, [isEditing])

  const handleContentClick = () => {
    setIsEditing(true)
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setContent(e.target.value)
  }

  const handleInputBlur = () => {
    setIsEditing(false)
  }

  const handleInputKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      setIsEditing(false)
    }
  }

  return (
    <AnnotationNode>
      <AnnotationNodeContent
        onClick={handleContentClick}
        style={{ cursor: 'text' }}
      >
        {isEditing ? (
          <textarea
            ref={inputRef}
            value={content}
            onChange={handleInputChange}
            onBlur={handleInputBlur}
            onKeyDown={handleInputKeyDown}
            style={{
              fontSize: 'inherit',
              fontFamily: 'inherit',
              border: '1px solid #ccc',
              borderRadius: 4,
              padding: '6px 8px',
              width: '100%',
              minHeight: '80px',
              boxSizing: 'border-box',
            }}
          />
        ) : (
          content
        )}
      </AnnotationNodeContent>
    </AnnotationNode>
  )
}
Annotation.displayName = 'Annotation'

export { Annotation }
