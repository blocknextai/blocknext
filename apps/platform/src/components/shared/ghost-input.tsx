import type { ComponentProps } from 'react'

type GhostInputProps = ComponentProps<'input'>

export function GhostInput({ className = '', ...props }: GhostInputProps) {
  return (
    <input
      {...props}
      className={`border-0! bg-transparent! outline-0! selection:bg-primary placeholder:text-muted-foreground selection:text-primary-foreground ${className}`}
    />
  )
}
