import { useRef, useImperativeHandle, forwardRef } from 'react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Send, Loader2, Plus, History } from 'lucide-react'

const ChatInput = forwardRef(
  (
    {
      input,
      streaming,
      initializing,
      onInputChange,
      onSend,
      onStartNewChat,
      onOpenHistory,
    },
    ref,
  ) => {
    const { t } = useTranslation()
    const textareaRef = useRef<HTMLTextAreaElement>(null)

    useImperativeHandle(ref, () => ({
      focus: () => textareaRef.current?.focus(),
      resetHeight: () => {
        if (textareaRef.current) {
          textareaRef.current.style.height = 'auto'
        }
      },
    }))

    const handleKeyDown = (e: React.KeyboardEvent) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault()
        onSend()
      }
    }

    const handleTextareaInput = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
      onInputChange(e.target.value)
      e.target.style.height = 'auto'
      e.target.style.height = Math.min(e.target.scrollHeight, 200) + 'px'
    }

    return (
      <div className="border-t p-3 flex flex-col gap-2">
        <div className="flex gap-2 items-end">
          <Textarea
            ref={textareaRef}
            value={input}
            onChange={handleTextareaInput}
            onKeyDown={handleKeyDown}
            placeholder={t('ui.text.generation.describeWorkflowPlaceholder')}
            disabled={streaming || initializing}
            rows={1}
            className="resize-none min-h-[40px] max-h-[160px] text-sm"
          />
          <Button
            onClick={onSend}
            disabled={!input.trim() || streaming}
            size="icon"
            className="shrink-0 h-[40px] w-[40px]"
          >
            {streaming ? (
              <Loader2 className="w-4 h-4 animate-spin" />
            ) : (
              <Send className="w-4 h-4" />
            )}
          </Button>
        </div>
        <div className="flex items-center gap-1">
          <Button
            variant="ghost"
            size="sm"
            className="h-7 text-xs text-muted-foreground"
            onClick={onStartNewChat}
          >
            <Plus className="w-3.5 h-3.5 mr-1" />
            {t('ui.text.generation.newChat')}
          </Button>
          <Button
            variant="ghost"
            size="sm"
            className="h-7 text-xs text-muted-foreground"
            onClick={onOpenHistory}
          >
            <History className="w-3.5 h-3.5 mr-1" />
            {t('ui.text.generation.history')}
          </Button>
        </div>
      </div>
    )
  },
)

ChatInput.displayName = 'ChatInput'

export default ChatInput
