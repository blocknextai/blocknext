import { Monitor, Moon, Sun } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { ToggleGroup, ToggleGroupItem } from '@/components/ui/toggle-group'
import { useThemeStore } from '@/stores/theme-store'

const modeItems = [
  {
    labelKey: 'ui.text.system',
    value: 'system',
    icon: Monitor,
  },
  {
    labelKey: 'ui.text.dark',
    value: 'dark',
    icon: Moon,
  },
  {
    labelKey: 'ui.text.light',
    value: 'light',
    icon: Sun,
  },
]

const ModeToggle = (props) => {
  const { t } = useTranslation()
  const mode = useThemeStore((s) => s.mode)
  const setTheme = useThemeStore((s) => s.setTheme)
  let collapsedCss = ''
  if (props?.direction === 'vertical') {
    collapsedCss = 'flex-col! gap-2 py-2'
  }

  const handleValueChange = (val: string) => {
    if (val) {
      setTheme(val)
    }
  }

  return (
    <ToggleGroup
      type="single"
      value={mode}
      onValueChange={handleValueChange}
      className={`flex items-center justify-center border-solid border rounded-lg! border-border p-0.5 overflow-hidden ${collapsedCss} ${props.className}`}
    >
      {modeItems.map((item) => (
        <ToggleGroupItem
          key={item.value}
          value={item.value}
          aria-label={t(item.labelKey)}
          className="flex-1 rounded-sm hover:bg-accent! hover:text-foreground! cursor-pointer! flex items-center justify-center gap-2"
        >
          <item.icon />
        </ToggleGroupItem>
      ))}
    </ToggleGroup>
  )
}

ModeToggle.displayName = 'ModeToggle'

export { ModeToggle }
