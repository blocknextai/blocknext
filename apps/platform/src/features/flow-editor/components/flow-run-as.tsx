import { useEffect, useState, useRef } from 'react'
import { useTranslation } from 'react-i18next'
import { DebouncedTextarea } from '@/components/shared/debounced-textarea'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Bot, ClockFading, SendHorizontal, Settings2 } from 'lucide-react'
import FlowRun from '@/features/flow-editor/components/flow-run'

const FlowRunAs = ({
  onClick,
  flowAvatar,
  defaultValue,
  cronString,
  setCronString,
  flow,
  credentials,
  setFlow,
  triggerType,
}) => {
  const { t } = useTranslation()
  const [runAs, setRunAs] = useState('member')
  const [prompt, setPrompt] = useState(defaultValue)
  const [open, setOpen] = useState(false)
  const [cron, setCron] = useState(cronString)
  const [toggleAdvanced, setToggleAdvanced] = useState(false)

  const wrapperRef = useRef(null)
  const contentRef = useRef(null)

  useEffect(() => {
    if (toggleAdvanced && wrapperRef.current && contentRef.current) {
      const contentHeight = contentRef.current.offsetHeight
      wrapperRef.current.style.paddingBottom = `${contentHeight}px`
    } else if (wrapperRef.current) {
      wrapperRef.current.style.paddingBottom = '0px'
    }
  }, [toggleAdvanced])
  return (
    <div ref={wrapperRef} className="relative transition-all duration-200">
      <div className="flex bg-input/30 rounded-xl flex-col items-start gap-2 py-3 px-4 ">
        {triggerType !== 'webhook' && (
          <DebouncedTextarea
            defaultValue={prompt}
            onChange={setPrompt}
            placeholder={t('ui.text.writeYourInspirationHere')}
            className="bg-transparent
              w-full text-foreground placeholder:text-foreground/50
              resize-none outline-0
              focus:ring-0
              focus:border-0
              focus-visible:ring-0
              max:h-50
              "
          />
        )}
        <div className="flex w-full gap-2 items-center justify-between">
          <div className="flex gap-3 w-full items-center">
            <div
              className="size-7 rounded-sm"
              style={{ backgroundImage: `${flowAvatar}` }}
            ></div>
            <div className="bg-transparent text-muted-foreground p-0 flex items-center">
              <Select onValueChange={setRunAs}>
                <SelectTrigger className="bg-transparent! rounded-lg!">
                  <SelectValue placeholder={t('ui.text.runAs')} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="user">{t('ui.text.user')}</SelectItem>
                  <SelectItem value="member">
                    {t('ui.text.organization')}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
            {triggerType === 'schedule' && (
              <Button
                className="rounded-lg"
                size={'sm'}
                variant={'outline'}
                onClick={() => setOpen(true)}
              >
                <ClockFading />
              </Button>
            )}
          </div>
          <div className="flex w-full gap-2 justify-end items-center">
            <Button
              className="rounded-lg"
              size={'sm'}
              variant={'ghost'}
              onClick={() => setToggleAdvanced(!toggleAdvanced)}
            >
              {toggleAdvanced ? (
                <>
                  <Bot /> {t('ui.text.basic')}
                </>
              ) : (
                <>
                  <Settings2 /> {t('ui.text.advanced')}
                </>
              )}
            </Button>
            <Button
              className="rounded-lg"
              size={'icon'}
              onClick={() => onClick(runAs, prompt)}
            >
              <SendHorizontal />
            </Button>
          </div>
        </div>
      </div>
      {toggleAdvanced && (
        <div ref={contentRef} className="absolute w-full mt-4">
          <FlowRun
            flow={flow}
            credentials={credentials}
            runType={runAs}
            setFlow={setFlow}
          />
        </div>
      )}
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>{t('ui.text.setSchedule')}</DialogTitle>
            <DialogDescription>
              {t('ui.text.setTimerForAgent')}
            </DialogDescription>
            <div className="flex gap-2 items-start w-full flex-col mt-4">
              <Select onValueChange={setCron} className="w-full!">
                <SelectTrigger className="bg-transparent! rounded-lg! w-full">
                  <SelectValue placeholder={t('ui.text.selectTime')} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="*/15 * * * *">
                    {t('ui.text.every15Minutes')}
                  </SelectItem>
                  <SelectItem value="*/30 * * * *">
                    {t('ui.text.every30Minutes')}
                  </SelectItem>
                  <SelectItem value="0 * * * *">
                    {t('ui.text.everyHour')}
                  </SelectItem>
                  <SelectItem value="0 0 * * *">
                    {t('ui.text.everyDay')}
                  </SelectItem>
                </SelectContent>
              </Select>
              <span className="text-center w-full">{t('ui.text.or')}</span>
              <Input
                defaultValue={cron}
                onChange={(e) => setCron(e.target.value)}
                placeholder={t('ui.text.cronFormat')}
              />
            </div>
          </DialogHeader>
          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline" size="sm">
                {t('ui.text.cancel')}
              </Button>
            </DialogClose>
            <DialogClose asChild>
              <Button
                size="sm"
                onClick={() => {
                  setCronString(cron)
                }}
              >
                {t('ui.text.save')}
              </Button>
            </DialogClose>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

FlowRunAs.displayName = 'FlowRunAs'

export default FlowRunAs
