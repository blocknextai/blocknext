import * as React from 'react'
import { useTranslation } from 'react-i18next'
import { Check } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'

interface StepProps {
  title: string
  description?: string
  isCompleted?: boolean
  isActive?: boolean
}

const Step: React.FC<StepProps> = ({ title, isCompleted, isActive }) => {
  return (
    <div className="flex items-center">
      <div className="relative flex items-center justify-center">
        <div
          className={cn(
            'w-8 h-8 rounded-full flex items-center justify-center',
            isCompleted
              ? 'border-primary bg-primary text-primary-foreground'
              : isActive
                ? 'bg-primary text-primary-foreground'
                : 'bg-muted text-foreground',
          )}
        >
          {isCompleted ? (
            <Check className="w-4 h-4" />
          ) : (
            <span className="text-sm font-medium">{title[0]}</span>
          )}
        </div>
      </div>
    </div>
  )
}

interface StepperProps {
  steps: Array<{ title: string; description?: string }>
  className?: string
  currentStep: number
  onStepChange: (step: number) => void
}

export function Stepper({
  steps,
  currentStep,
  onStepChange,
  className,
}: StepperProps) {
  const { t } = useTranslation()
  return (
    <div className={`w-full flex justify-between gap-6 ${className}`}>
      {currentStep === 0 ? (
        <span className="w-24 py-2 px-4">&nbsp;</span>
      ) : (
        <Button
          variant={'secondary'}
          onClick={() => onStepChange(currentStep - 1)}
          className={`w-24`}
        >
          {t('ui.text.previous')}
        </Button>
      )}

      <div className="flex flex-col md:flex-row justify-center items-center md:items-center w-full max-w-xs m-auto">
        {steps.map((step, index) => (
          <React.Fragment key={step.title}>
            <Step
              title={step.title}
              description={step.description}
              isCompleted={index < currentStep}
              isActive={index === currentStep}
            />
            {index < steps.length - 1 && (
              <Separator
                className="w-16! bg-muted/60! py-0.5!  m-0!"
                decorative={true}
              />
            )}
          </React.Fragment>
        ))}
      </div>
      {currentStep === 0 ? (
        <span className="w-24 py-2 px-4">&nbsp;</span>
      ) : (
        <Button className="w-24" onClick={() => onStepChange(currentStep + 1)}>
          {currentStep === steps.length - 1 ? 'Run' : 'Next'}
        </Button>
      )}
    </div>
  )
}
