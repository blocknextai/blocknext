import { forwardRef } from 'react'
import { cn } from '@/lib/utils'

const AnnotationNode = forwardRef(({ className, children, ...props }, ref) => {
  return (
    <div
      ref={ref}
      {...props}
      className={cn(
        'relative flex max-w-[180px] items-start p-2 text-sm text-secondary-foreground',
        className,
      )}
    >
      {children}
    </div>
  )
})

AnnotationNode.displayName = 'AnnotationNode'

const AnnotationNodeNumber = forwardRef(
  ({ className, children, ...props }, ref) => {
    return (
      <div ref={ref} {...props} className={cn('mr-1 leading-snug', className)}>
        {children}
      </div>
    )
  },
)

AnnotationNodeNumber.displayName = 'AnnotationNodeNumber'

const AnnotationNodeContent = forwardRef(
  ({ className, children, ...props }, ref) => {
    return (
      <div ref={ref} {...props} className={cn('leading-snug', className)}>
        {children}
      </div>
    )
  },
)

AnnotationNodeContent.displayName = 'AnnotationNodeContent'

const AnnotationNodeIcon = forwardRef(
  ({ className, children, ...props }, ref) => {
    return (
      <div
        ref={ref}
        {...props}
        className={cn('absolute bottom-0 right-2 text-2xl', className)}
      >
        {children}
      </div>
    )
  },
)

AnnotationNodeIcon.displayName = 'AnnotationNodeIcon'

export {
  AnnotationNode,
  AnnotationNodeNumber,
  AnnotationNodeContent,
  AnnotationNodeIcon,
}
