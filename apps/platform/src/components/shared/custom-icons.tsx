import { generate } from 'geopattern'
import { Avatar, Style } from '@dicebear/core'
import bottts from '@dicebear/styles/bottts.json'
import identicon from '@dicebear/styles/identicon.json'

const botttsStyle = new Style(bottts)
const identiconStyle = new Style(identicon)

interface UserIconProps {
  seed: string
  className?: string
  size?: number
}

export const UserIcon2 = ({ seed, className, size = 20 }: UserIconProps) => {
  let svg = new Avatar(botttsStyle, {
    seed: seed,
    size: size,
    borderRadius: 0,
    scale: 1.6,
  }).toString()
  svg = svg.replaceAll('viewboxMask', `viewboxMask-${seed}`)
  return (
    <div
      className={className}
      dangerouslySetInnerHTML={{ __html: svg }}
      style={{ width: size, height: size }}
    />
  )
}

const FLOW_ICON_COLORS = [
  '#64748b',
  '#14b8a6',
  '#3b3b5f',
  '#4f46e5',
  '#8b5cf6',
  '#8b5cf6',
  '#e879f9',
  '#f59e0b',
]

const flowIconCache = new Map<string, string>()

export const FlowIcon = (seed: string) => {
  const cached = flowIconCache.get(seed)
  if (cached) return cached

  let hash = 0
  for (let i = 0; i < seed.length; i++) hash += seed.charCodeAt(i)
  const baseColor = FLOW_ICON_COLORS[hash % FLOW_ICON_COLORS.length]
  const dataUrl = generate(seed, { baseColor }).toDataUrl()

  flowIconCache.set(seed, dataUrl)
  return dataUrl
}

interface CompanyIconProps {
  seed: string
  className?: string
  size?: number
  avatarColor?: string
}

export const CompanyIcon = ({
  seed,
  className,
  size,
  avatarColor,
}: CompanyIconProps) => {
  let svg = new Avatar(identiconStyle, {
    seed: seed,
    size: size,
    borderRadius: 16,
    rowColor: avatarColor ? [avatarColor] : undefined,
  }).toString()
  svg = svg.replaceAll('viewboxMask', `viewboxMask-${seed}`)
  return (
    <div
      className={className}
      dangerouslySetInnerHTML={{ __html: svg }}
      style={{ width: size, height: size }}
    />
  )
}
