import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { authService } from '@/features/auth'
import tokenManager from '@/lib/token-manager'

export const ThemeColors = [
  'neon',
  'blue',
  'violet',
  'navy',
  'terminal',
  'green',
  'olive',
  'yellow',
  'orange',
  'red',
  'pastel-red',
  'pink',
  'default',
] as const

type ThemeColor = (typeof ThemeColors)[number]
type ThemeMode = 'light' | 'dark' | 'system'

interface ThemeState {
  color: ThemeColor
  mode: ThemeMode
  getMode: () => ThemeMode
  setColor: (color: ThemeColor) => void
  setTheme: (mode: ThemeMode) => void
  syncFromBackend: (theme: { mode?: string; color?: string }) => void
}

function resolveMode(mode: ThemeMode) {
  if (mode === 'system') {
    return window.matchMedia('(prefers-color-scheme: dark)').matches
      ? 'dark'
      : 'light'
  }
  return mode
}
function applyTheme(mode: ThemeMode, color: ThemeColor) {
  const root = document.documentElement
  const resolvedMode = resolveMode(mode)

  root.classList.remove('light', 'dark')
  root.classList.add(resolvedMode)

  root.removeAttribute('data-theme')
  root.setAttribute('data-theme', `${color}-${resolvedMode}`)
}

export const useThemeStore = create<ThemeState>()(
  persist(
    (set, get) => ({
      color: 'default',
      mode: 'system',
      getMode: () => resolveMode(get().mode),
      setColor: (color) => {
        const { mode } = get()
        applyTheme(mode, color)
        set({ color })
        if (tokenManager.getAccessToken()) {
          authService.updatePreferences({ theme: { color } }).catch(() => {})
        }
      },
      setTheme: (mode) => {
        const { color } = get()
        applyTheme(mode, color)
        set({ mode })
        if (tokenManager.getAccessToken()) {
          authService.updatePreferences({ theme: { mode } }).catch(() => {})
        }
      },
      syncFromBackend: (theme) => {
        const mode = (theme?.mode || 'system') as ThemeMode
        const color = (theme?.color || 'default') as ThemeColor
        applyTheme(mode, color)
        set({ mode, color })
      },
    }),
    {
      name: 'vite-ui-theme',
    },
  ),
)

export function applyStoredTheme() {
  try {
    const raw = localStorage.getItem('vite-ui-theme')
    let mode: ThemeMode = 'system'
    let color: ThemeColor = 'default'

    if (raw) {
      const data = JSON.parse(raw)
      mode = data.state?.mode || 'system'
      color = data.state?.color || 'default'
    }

    applyTheme(mode, color)
  } catch (e) {
    console.error('Failed to apply stored theme:', e)
    applyTheme('system', 'default')
  }
}
