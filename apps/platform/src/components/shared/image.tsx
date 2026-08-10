import type { ImgHTMLAttributes } from 'react'

type ImageProps = ImgHTMLAttributes<HTMLImageElement> & {
  /** Set on anything visible before the page is scrolled. */
  priority?: boolean
}

const Image = ({ src, priority = false, ...rest }: ImageProps) => {
  if (!src) {
    return null
  }

  return (
    <img
      loading={priority ? 'eager' : 'lazy'}
      fetchPriority={priority ? 'high' : undefined}
      decoding="async"
      {...rest}
      src={src}
    />
  )
}

Image.displayName = 'Image'

export { Image }
