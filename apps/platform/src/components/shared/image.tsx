import type { ImgHTMLAttributes } from 'react'

type ImageProps = ImgHTMLAttributes<HTMLImageElement>

const Image = ({ src, ...rest }: ImageProps) => {
  if (!src) {
    return null
  }

  return <img loading="lazy" {...rest} src={src} />
}

Image.displayName = 'Image'

export { Image }
