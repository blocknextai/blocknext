import { useTranslation } from 'react-i18next'
import i18next from 'i18next'
import {
  FileText,
  Image as ImageIcon,
  Video,
  Music,
  Link,
  CheckCircle2,
  File,
  ExternalLink,
  Copy,
  Table,
} from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import toast from '@/lib/toast'

const isUrl = (value: string) => {
  try {
    const url = new URL(value)
    return url.protocol === 'http:' || url.protocol === 'https:'
  } catch {
    return false
  }
}

const isImageUrl = (value: string) => {
  if (!isUrl(value)) return false
  const ext = value.split('?')[0].toLowerCase()
  return (
    ext.endsWith('.png') ||
    ext.endsWith('.jpg') ||
    ext.endsWith('.jpeg') ||
    ext.endsWith('.gif') ||
    ext.endsWith('.webp') ||
    ext.endsWith('.svg')
  )
}

const isVideoUrl = (value: string) => {
  if (!isUrl(value)) return false
  const ext = value.split('?')[0].toLowerCase()
  return (
    ext.endsWith('.mp4') ||
    ext.endsWith('.webm') ||
    ext.endsWith('.mov') ||
    ext.endsWith('.avi')
  )
}

const isAudioUrl = (value: string) => {
  if (!isUrl(value)) return false
  const ext = value.split('?')[0].toLowerCase()
  return (
    ext.endsWith('.mp3') ||
    ext.endsWith('.wav') ||
    ext.endsWith('.ogg') ||
    ext.endsWith('.m4a')
  )
}

const copyToClipboard = (text: string) => {
  navigator.clipboard.writeText(text)
  toast.success(i18next.t('ui.text.copiedToClipboard'))
}

const TextOutput = ({ value }: { value: string }) => (
  <div className="bg-muted/50 rounded-lg p-3 text-sm">
    <div className="flex items-start gap-2">
      <FileText className="size-4 mt-0.5 shrink-0 text-muted-foreground" />
      <p className="whitespace-pre-wrap break-words flex-1">{value}</p>
    </div>
    <div className="flex justify-end mt-2">
      <Button
        variant="ghost"
        size="sm"
        className="h-6 text-xs"
        onClick={() => copyToClipboard(value)}
      >
        <Copy className="size-3 mr-1" />
        {i18next.t('ui.text.copy')}
      </Button>
    </div>
  </div>
)

const ImageOutput = ({ value }: { value: string }) => (
  <div className="rounded-lg overflow-hidden border">
    <a href={value} target="_blank" rel="noopener noreferrer">
      <img
        src={value}
        alt={i18next.t('ui.text.dataOutput')}
        className="max-w-full max-h-80 object-contain bg-muted/30"
        loading="lazy"
      />
    </a>
    <div className="flex items-center justify-between p-2 bg-muted/30 text-xs">
      <div className="flex items-center gap-1 text-muted-foreground">
        <ImageIcon className="size-3" />
        {i18next.t('ui.text.imageOutput')}
      </div>
      <a
        href={value}
        target="_blank"
        rel="noopener noreferrer"
        className="flex items-center gap-1 text-primary hover:underline"
      >
        {i18next.t('ui.text.openLink')} <ExternalLink className="size-3" />
      </a>
    </div>
  </div>
)

const VideoOutput = ({ value }: { value: string }) => (
  <div className="rounded-lg overflow-hidden border">
    <video controls className="max-w-full max-h-80 bg-muted/30">
      <source src={value} />
    </video>
    <div className="flex items-center justify-between p-2 bg-muted/30 text-xs">
      <div className="flex items-center gap-1 text-muted-foreground">
        <Video className="size-3" />
        {i18next.t('ui.text.videoOutput')}
      </div>
      <a
        href={value}
        target="_blank"
        rel="noopener noreferrer"
        className="flex items-center gap-1 text-primary hover:underline"
      >
        {i18next.t('ui.text.openLink')} <ExternalLink className="size-3" />
      </a>
    </div>
  </div>
)

const AudioOutput = ({ value }: { value: string }) => (
  <div className="rounded-lg border p-3">
    <div className="flex items-center gap-2 mb-2 text-xs text-muted-foreground">
      <Music className="size-3" />
      {i18next.t('ui.text.audioOutput')}
    </div>
    <audio controls className="w-full">
      <source src={value} />
    </audio>
  </div>
)

const LinkOutput = ({ value }: { value: string }) => (
  <div className="flex items-center gap-2 bg-muted/50 rounded-lg p-3 text-sm">
    <Link className="size-4 shrink-0 text-muted-foreground" />
    <a
      href={value}
      target="_blank"
      rel="noopener noreferrer"
      className="text-primary hover:underline truncate flex-1"
    >
      {value}
    </a>
    <ExternalLink className="size-3 shrink-0 text-muted-foreground" />
  </div>
)

const StatusOutput = ({ value }: { value: string }) => (
  <div className="flex items-center gap-2">
    <CheckCircle2 className="size-4 text-green-500" />
    <Badge variant="secondary">{value}</Badge>
  </div>
)

const FileOutput = ({ value }: { value: string }) => (
  <div className="flex items-center gap-2 bg-muted/50 rounded-lg p-3 text-sm">
    <File className="size-4 shrink-0 text-muted-foreground" />
    {isUrl(value) ? (
      <a
        href={value}
        target="_blank"
        rel="noopener noreferrer"
        className="text-primary hover:underline truncate flex-1"
      >
        {value}
      </a>
    ) : (
      <span className="truncate flex-1">{value}</span>
    )}
  </div>
)

const ObjectOutput = ({ value }: { value: any }) => (
  <div className="bg-muted/50 rounded-lg p-3 text-xs">
    <div className="flex items-center gap-2 mb-2 text-muted-foreground">
      <Table className="size-3" />
      {i18next.t('ui.text.objectOutput')}
    </div>
    <pre className="whitespace-pre-wrap break-all overflow-hidden">
      {typeof value === 'string' ? value : JSON.stringify(value, null, 2)}
    </pre>
  </div>
)

const renderValue = (key: string, value: any) => {
  if (value === null || value === undefined) return null

  const strValue = typeof value === 'string' ? value : String(value)

  // Detect by key name hints
  if (key === 'image' || key.includes('image') || key.includes('thumbnail')) {
    if (isUrl(strValue)) return <ImageOutput value={strValue} />
  }

  if (key === 'video' || key.includes('video')) {
    if (isVideoUrl(strValue)) return <VideoOutput value={strValue} />
    if (isUrl(strValue)) return <LinkOutput value={strValue} />
  }

  if (key === 'audio' || key.includes('audio')) {
    if (isAudioUrl(strValue)) return <AudioOutput value={strValue} />
    if (isUrl(strValue)) return <LinkOutput value={strValue} />
  }

  if (key === 'file') {
    return <FileOutput value={strValue} />
  }

  if (key === 'status') {
    return <StatusOutput value={strValue} />
  }

  if (key === 'videoUrl' || key === 'videoId') {
    if (isUrl(strValue)) return <LinkOutput value={strValue} />
    return <TextOutput value={strValue} />
  }

  // Detect by content
  if (typeof value === 'string') {
    if (isImageUrl(value)) return <ImageOutput value={value} />
    if (isVideoUrl(value)) return <VideoOutput value={value} />
    if (isAudioUrl(value)) return <AudioOutput value={value} />
    if (isUrl(value)) return <LinkOutput value={value} />
    return <TextOutput value={value} />
  }

  if (typeof value === 'object') {
    return <ObjectOutput value={value} />
  }

  return <TextOutput value={strValue} />
}

const OutputRenderer = ({ outputs }: { outputs: any }) => {
  const { t } = useTranslation()

  if (!outputs) return null

  // outputs can be an array of objects or a single object
  const items = Array.isArray(outputs) ? outputs : [outputs]

  return (
    <div>
      <h4 className="text-sm font-medium mb-3 flex items-center gap-2">
        <CheckCircle2 className="size-4" />
        {t('ui.text.outputs')}
      </h4>
      <div className="space-y-3">
        {items.map((item, itemIndex) => {
          if (!item || typeof item !== 'object') return null

          return (
            <div key={itemIndex} className="space-y-2">
              {Object.entries(item).map(([key, value]) => {
                if (value === null || value === undefined) return null

                return (
                  <div key={key}>
                    <div className="text-xs text-muted-foreground mb-1 font-mono">
                      {key}
                    </div>
                    {renderValue(key, value)}
                  </div>
                )
              })}
            </div>
          )
        })}
      </div>
    </div>
  )
}

OutputRenderer.displayName = 'OutputRenderer'

export { OutputRenderer }
