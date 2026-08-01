import { useTranslation } from 'react-i18next'
import { ChevronDown } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { ThemeColors } from '@/stores/theme-store'

const AppearanceSection = ({
  mode,
  setTheme,
  getMode,
  color,
  setColor,
  changeLanguage,
}) => {
  const { t } = useTranslation()

  const renderColors = () => {
    const colorArray = []
    for (let i = 0; i < ThemeColors.length; i++) {
      const c = ThemeColors[i]
      colorArray.push(
        <div
          title={c}
          key={i}
          className={`w-10 h-10 rounded-full p-0.5 cursor-pointer border-foreground ${color === c ? 'border-2' : 'border-0'}`}
          onClick={() => setColor(c)}
        >
          <div
            data-theme={`${c}-${getMode()}`}
            className={`w-full h-full rounded-full bg-primary`}
          ></div>
        </div>,
      )
    }
    return colorArray
  }

  return (
    <div className="flex flex-col gap-4 w-full min-w-0 p-2 sm:p-4">
      <div className="flex justify-between items-center">
        <div className="grid gap-3">
          <Label htmlFor={`theme`}>{t('ui.text.theme')}</Label>
          <span id="theme" className="text-muted-foreground text-sm!">
            {t('ui.text.chooseTheme')}
          </span>
        </div>
        <div>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline">
                {mode.charAt(0).toUpperCase() + mode.slice(1)}
                <ChevronDown />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem onClick={() => setTheme('system')}>
                {t('ui.text.system')}
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => setTheme('dark')}>
                {t('ui.text.dark')}
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => setTheme('light')}>
                {t('ui.text.light')}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
      <Separator className="my-2 p-0" decorative={true} />

      <div className="flex justify-between items-center">
        <div className="grid gap-3">
          <Label htmlFor={`language`}>{t('ui.text.language')}</Label>
          <span id="language" className="text-muted-foreground text-sm!">
            {t('ui.text.chooseLanguage')}
          </span>
        </div>
        <div>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline">
                {t('ui.text.english')}
                <ChevronDown />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem onClick={() => changeLanguage('en')}>
                {t('ui.text.english')}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
      <Separator className="my-2 p-0" decorative={true} />

      <div className="flex justify-between gap-2 items-center">
        <div className="grid gap-3">
          <Label htmlFor={`color`}>{t('ui.text.color')}</Label>
          <span id="color" className="text-muted-foreground text-sm!">
            {t('ui.text.chooseColor')}
          </span>
        </div>
        <div className="flex flex-wrap gap-2">{renderColors()}</div>
      </div>
    </div>
  )
}

AppearanceSection.displayName = 'AppearanceSection'

export { AppearanceSection }
