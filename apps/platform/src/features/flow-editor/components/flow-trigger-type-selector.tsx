import { useTranslation } from 'react-i18next'
import { MousePointerClick, ClockFading, Webhook } from 'lucide-react'

const triggerTypes = [
  {
    id: 'manual',
    icon: MousePointerClick,
    titleKey: 'ui.text.manual',
    descriptionKey: 'ui.text.manualTriggerDescription',
  },
  {
    id: 'schedule',
    icon: ClockFading,
    titleKey: 'ui.text.schedule',
    descriptionKey: 'ui.text.scheduleTriggerDescription',
  },
  {
    id: 'webhook',
    icon: Webhook,
    titleKey: 'ui.text.webhook',
    descriptionKey: 'ui.text.webhookTriggerDescription',
  },
]

interface FlowTriggerTypeSelectorProps {
  onSelect: (triggerType: string) => void
  selected?: string
}

const FlowTriggerTypeSelector = ({
  onSelect,
  selected,
}: FlowTriggerTypeSelectorProps) => {
  const { t } = useTranslation()

  return (
    <div className="flex flex-col gap-4 w-full">
      <div className="text-center">
        <h3 className="text-lg font-semibold">
          {t('ui.text.selectTriggerType')}
        </h3>
        <p className="text-sm text-muted-foreground mt-1">
          {t('ui.text.howShouldFlowStart')}
        </p>
      </div>
      <div className="grid grid-cols-3 gap-4">
        {triggerTypes.map((trigger) => {
          const Icon = trigger.icon
          const isSelected = selected === trigger.id
          return (
            <button
              key={trigger.id}
              onClick={() => onSelect(trigger.id)}
              className={`flex flex-col items-center gap-3 p-6 rounded-xl border-2 transition-all duration-200 cursor-pointer hover:shadow-md ${
                isSelected
                  ? 'border-primary bg-primary/5 shadow-md'
                  : 'border-border hover:border-primary/50'
              }`}
            >
              <div
                className={`size-12 rounded-lg flex items-center justify-center ${
                  isSelected
                    ? 'bg-primary text-primary-foreground'
                    : 'bg-muted text-muted-foreground'
                }`}
              >
                <Icon className="size-6" />
              </div>
              <div className="text-center">
                <div className="font-medium text-sm">{t(trigger.titleKey)}</div>
                <div className="text-xs text-muted-foreground mt-1">
                  {t(trigger.descriptionKey)}
                </div>
              </div>
            </button>
          )
        })}
      </div>
    </div>
  )
}

FlowTriggerTypeSelector.displayName = 'FlowTriggerTypeSelector'

export default FlowTriggerTypeSelector
