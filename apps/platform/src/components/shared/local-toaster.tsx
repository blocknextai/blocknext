import { Toaster } from 'sonner'
import { useThemeStore } from '@/stores/theme-store'

export function LocalToaster() {
  const mode = useThemeStore((s) => s.mode)

  const appliedTheme =
    mode === 'system'
      ? window.matchMedia('(prefers-color-scheme: dark)').matches
        ? 'dark'
        : 'light'
      : mode
  return (
    <Toaster
      richColors
      theme={appliedTheme}
      position="top-right"
      closeButton
      icons={{
        success: '✅',
        error: '❌',
        info: 'ℹ️',
        warning: '⚠️',
        loading: '⏳',
      }}
    />
  )
}
