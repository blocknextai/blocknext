import { useCallback, useEffect, useRef, useState } from 'react'

/**
 * Shared copy-to-clipboard logic with a transient "copied" flag, so any
 * surface (button, row, field) can trigger a copy and show feedback.
 */
export function useCopy(resetMs = 2000) {
  const [copied, setCopied] = useState(false)
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current)
      }
    }
  }, [])

  const copy = useCallback(
    (value: string) => {
      if (!value) {
        return
      }
      navigator.clipboard.writeText(value)
      setCopied(true)
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current)
      }
      timeoutRef.current = setTimeout(() => setCopied(false), resetMs)
    },
    [resetMs],
  )

  return { copied, copy }
}
