import React, { useState, useRef, useEffect } from 'react'
import { Textarea } from '@/components/ui/textarea'

export const AdvancedDebouncedTextarea = React.memo(
  function AdvancedDebouncedTextarea(props: {
    value?: string
    onChange?: (value: string) => void
    onFocus?: React.FocusEventHandler
    onBlur?: React.FocusEventHandler
    mapping?: Array<{ value: string; label: string }>
    placeholder?: string
    className?: string
  }) {
    const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)
    const [text, setText] = useState(props.value || '')

    useEffect(() => {
      if (props.value !== undefined && props.value !== text) {
        setText(props.value)
      }
    }, [props.value])

    useEffect(() => {
      return () => {
        if (timeoutRef.current) clearTimeout(timeoutRef.current)
      }
    }, [])

    const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
      const newText = e.target.value
      setText(newText)
      if (timeoutRef.current) clearTimeout(timeoutRef.current)
      timeoutRef.current = setTimeout(() => {
        props.onChange?.(newText)
      }, 500)
    }

    return (
      <Textarea
        value={text}
        onChange={handleChange}
        onFocus={props.onFocus}
        onBlur={props.onBlur}
        placeholder={props.placeholder}
        className={props.className}
      />
    )
  },
  (prevProps, nextProps) => {
    return (
      prevProps.value === nextProps.value &&
      prevProps.onChange === nextProps.onChange
    )
  },
)
