import { useCallback } from 'react'
import type { ReactNode } from 'react'

import type { IconSource } from '@/features/flow-editor/types'
import { Image } from '@/components/shared/image'
import { cn } from '@/lib/utils'
import { useThemeStore } from '@/stores/theme-store'

const brandSource = (brand: string, mode: 'light' | 'dark') =>
  `/assets/icons/brands/${brand}/${mode}.svg`
const glyphSource = (glyph: string) => `/assets/icons/glyphs/${glyph}.svg`

type IconProps = {
  className?: string
}

export type ProviderIconComponent = (props: IconProps) => ReactNode

const builtIcons = new Map<string, ProviderIconComponent>()

const artwork = (source: string, className?: string) => (
  <Image
    src={source}
    alt=""
    aria-hidden
    className={cn('size-4 object-contain', className)}
  />
)

const badge = (glyph: string) => (
  <span
    aria-hidden
    className="bg-background ring-border absolute right-0 bottom-0 box-content flex size-[58%] items-center justify-center rounded-full border-[0.1em] border-transparent ring-1 ring-inset"
    style={{ backgroundClip: 'padding-box' }}
  >
    <span
      className="bg-foreground size-[68%]"
      style={{
        maskImage: `url(${glyphSource(glyph)})`,
        WebkitMaskImage: `url(${glyphSource(glyph)})`,
        maskSize: 'contain',
        WebkitMaskSize: 'contain',
        maskRepeat: 'no-repeat',
        WebkitMaskRepeat: 'no-repeat',
        maskPosition: 'center',
        WebkitMaskPosition: 'center',
      }}
    />
  </span>
)

const buildIcon = (
  icon: { brand?: string; glyph?: string },
  mode: 'light' | 'dark',
): ProviderIconComponent | undefined => {
  if (icon.brand && icon.glyph) {
    return ({ className }: IconProps) => (
      <span className={cn('relative inline-block size-4 shrink-0', className)}>
        {artwork(brandSource(icon.brand!, mode), 'size-[70%]')}
        {badge(icon.glyph!)}
      </span>
    )
  }

  if (icon.brand) {
    return ({ className }: IconProps) =>
      artwork(
        brandSource(icon.brand!, mode),
        cn('inline-block shrink-0', className),
      )
  }

  if (icon.glyph) {
    return ({ className }: IconProps) => (
      <span
        aria-hidden
        className={cn('bg-foreground inline-block size-4 shrink-0', className)}
        style={{
          maskImage: `url(${glyphSource(icon.glyph!)})`,
          WebkitMaskImage: `url(${glyphSource(icon.glyph!)})`,
          maskSize: 'contain',
          WebkitMaskSize: 'contain',
          maskRepeat: 'no-repeat',
          WebkitMaskRepeat: 'no-repeat',
          maskPosition: 'center',
          WebkitMaskPosition: 'center',
        }}
      />
    )
  }

  return undefined
}

export const useIconResolver = () => {
  const mode = useThemeStore((state) => state.getMode())

  return useCallback(
    (icon: IconSource) => {
      if (!icon) {
        return undefined
      }

      const key = `${mode}:${icon.brand ?? ''}:${icon.glyph ?? ''}`
      const existing = builtIcons.get(key)
      if (existing) {
        return existing
      }

      const built = buildIcon(icon, mode)
      if (built) {
        builtIcons.set(key, built)
      }

      return built
    },
    [mode],
  )
}
