import { Textarea } from '@/components/ui/textarea'
import { useState, useEffect, useRef } from 'react'
import type { ChangeEvent, ComponentProps } from 'react'

type TextareaProps = ComponentProps<typeof Textarea>

interface DebouncedTextareaProps extends Omit<
  TextareaProps,
  'onChange' | 'value'
> {
  value: string
  onChange: (value: string) => void
}

export function DebouncedTextarea({
  value,
  onChange,
  ...rest
}: DebouncedTextareaProps) {
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const [text, setText] = useState(value)

  const handleTextChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    const newText = event.target.value
    setText(newText)
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current)
    }
    timeoutRef.current = setTimeout(() => {
      onChange(newText)
    }, 500)
  }

  useEffect(() => {
    if (value !== text) {
      setText(value)
    }
  }, [value])

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current)
      }
    }
  }, [])

  return <Textarea {...rest} value={text} onChange={handleTextChange} />
}
