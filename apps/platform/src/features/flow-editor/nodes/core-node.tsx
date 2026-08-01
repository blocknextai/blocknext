import { memo, useState, useMemo, useCallback, useEffect } from 'react'
import { Handle, Position, useUpdateNodeInternals } from '@xyflow/react'
import {
  useFlowSetNodes,
  useFlowSetEdges,
} from '@/features/flow-editor/contexts/flow-nodes-context'
import { Card, CardContent } from '@/components/ui/card'
import {
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from '@/components/ui/tooltip'
import { useTranslation } from 'react-i18next'
import { flowCategoryPreferences } from '@/lib/flow-categories'
import { directionIcons } from '@/features/flow-editor/icons'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
  ContextMenuSub,
  ContextMenuSubTrigger,
  ContextMenuSubContent,
} from '@/components/ui/context-menu'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

const LAYOUT = {
  'l-r': [Position.Left, Position.Right],
  'r-l': [Position.Right, Position.Left],
  't-b': [Position.Top, Position.Bottom],
  'b-t': [Position.Bottom, Position.Top],
  'l-t': [Position.Left, Position.Top],
  'r-t': [Position.Right, Position.Top],
  't-l': [Position.Top, Position.Left],
  't-r': [Position.Top, Position.Right],
  'l-b': [Position.Left, Position.Bottom],
  'r-b': [Position.Right, Position.Bottom],
  'b-l': [Position.Bottom, Position.Left],
  'b-r': [Position.Bottom, Position.Right],
}
const CoreNode = memo(({ selected, data }) => {
  const updateNodeInternals = useUpdateNodeInternals()
  const setNodes = useFlowSetNodes()
  const setEdges = useFlowSetEdges()
  const [nameOpen, setNameOpen] = useState(false)
  const [name, setName] = useState(data?.title ?? '')
  const [position, setPosition] = useState(data?.position || 'l-r')
  const { t } = useTranslation()
  const { prefs, Icon } = useMemo(() => {
    if (!data?.category || !flowCategoryPreferences[data.category]) {
      return { prefs: { color: '#ffffff' }, Icon: () => null }
    }
    const p = flowCategoryPreferences[data.category]
    return { prefs: p, Icon: p.icon }
  }, [data?.category])

  const updateField = useCallback(
    (field, value) => {
      if (!data?.id) return
      setNodes((nodes) =>
        nodes.map((n) => {
          if (n.id === data.id) {
            return {
              ...n,
              data: {
                ...n.data,
                [field]: value,
              },
            }
          }
          return n
        }),
      )
    },
    [data?.id, setNodes],
  )

  const deleteNode = useCallback(() => {
    if (!data?.id) return
    setNodes((nodes) => nodes.filter((n) => n.id !== data.id))
    setEdges((edges) =>
      edges.filter((e) => e.source !== data.id && e.target !== data.id),
    )
  }, [data?.id, setNodes, setEdges])

  useEffect(() => {
    if (!data) return
    if (!data.position || (data.position && data.position !== position)) {
      updateField('position', position)
    }
  }, [position])
  useEffect(() => {
    if (!data?.id) return
    updateNodeInternals(data.id)
  }, [position, updateNodeInternals, data?.id])

  useEffect(() => {
    if (data?.title !== undefined && data.title !== name) {
      setName(data.title)
    }
  }, [data?.title])
  const renderPositionItems = useMemo(() => {
    return (
      <>
        <ContextMenuSub>
          <ContextMenuSubTrigger>{t('ui.text.leftTo')}</ContextMenuSubTrigger>
          <ContextMenuSubContent className="w-44">
            <ContextMenuItem
              className="gap-2"
              onClick={() => setPosition('l-r')}
            >
              {' '}
              <directionIcons.lr className="size-4" /> {t('ui.text.right')}{' '}
            </ContextMenuItem>
            <ContextMenuItem
              className="gap-2"
              onClick={() => setPosition('l-t')}
            >
              {' '}
              <directionIcons.lt className="size-4" /> {t('ui.text.top')}{' '}
            </ContextMenuItem>
            <ContextMenuItem
              className="gap-2"
              onClick={() => setPosition('l-b')}
            >
              {' '}
              <directionIcons.lb className="size-4" />{' '}
              {t('ui.text.bottom')}{' '}
            </ContextMenuItem>
          </ContextMenuSubContent>
        </ContextMenuSub>
        <ContextMenuSub>
          <ContextMenuSubTrigger>{t('ui.text.topTo')}</ContextMenuSubTrigger>
          <ContextMenuSubContent className="w-44">
            <ContextMenuItem
              className="gap-2"
              onClick={() => setPosition('t-l')}
            >
              {' '}
              <directionIcons.tl className="size-4" /> {t('ui.text.left')}{' '}
            </ContextMenuItem>
            <ContextMenuItem
              className="gap-2"
              onClick={() => setPosition('t-b')}
            >
              {' '}
              <directionIcons.tb className="size-4" />{' '}
              {t('ui.text.bottom')}{' '}
            </ContextMenuItem>
            <ContextMenuItem
              className="gap-2"
              onClick={() => setPosition('t-r')}
            >
              {' '}
              <directionIcons.tr className="size-4" /> {t('ui.text.right')}{' '}
            </ContextMenuItem>
          </ContextMenuSubContent>
        </ContextMenuSub>
        <ContextMenuSub>
          <ContextMenuSubTrigger>{t('ui.text.bottomTo')}</ContextMenuSubTrigger>
          <ContextMenuSubContent className="w-44">
            <ContextMenuItem
              className="gap-2"
              onClick={() => setPosition('b-l')}
            >
              {' '}
              <directionIcons.bl className="size-4" /> {t('ui.text.left')}{' '}
            </ContextMenuItem>
            <ContextMenuItem
              className="gap-2"
              onClick={() => setPosition('b-t')}
            >
              {' '}
              <directionIcons.bt className="size-4" /> {t('ui.text.top')}{' '}
            </ContextMenuItem>
            <ContextMenuItem
              className="gap-2"
              onClick={() => setPosition('b-r')}
            >
              {' '}
              <directionIcons.br className="size-4" /> {t('ui.text.right')}{' '}
            </ContextMenuItem>
          </ContextMenuSubContent>
        </ContextMenuSub>
        <ContextMenuSub>
          <ContextMenuSubTrigger>{t('ui.text.rightTo')}</ContextMenuSubTrigger>
          <ContextMenuSubContent className="w-44">
            <ContextMenuItem
              className="gap-2"
              onClick={() => setPosition('r-l')}
            >
              {' '}
              <directionIcons.rl className="size-4" /> {t('ui.text.left')}{' '}
            </ContextMenuItem>
            <ContextMenuItem
              className="gap-2"
              onClick={() => setPosition('r-t')}
            >
              {' '}
              <directionIcons.rt className="size-4" /> {t('ui.text.top')}{' '}
            </ContextMenuItem>
            <ContextMenuItem
              className="gap-2"
              onClick={() => setPosition('r-b')}
            >
              {' '}
              <directionIcons.rb className="size-4" />{' '}
              {t('ui.text.bottom')}{' '}
            </ContextMenuItem>
          </ContextMenuSubContent>
        </ContextMenuSub>
      </>
    )
  }, [position, updateField, t])
  return (
    <>
      <ContextMenu>
        <ContextMenuTrigger>
          <div
            className={`rounded-xl relative bg-background transition-all duration-200 hover:shadow-md border-2 border-current spread-transition
                             ${selected ? 'spread-container' : ''}`}
            style={{ color: prefs?.color }}
          >
            <div className="bg-current/20 rounded-lg">
              {/* Giriş Noktası */}
              {!data?.hide_left && (
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Handle
                      type="target"
                      position={LAYOUT[position][0]}
                      className="border-ring hover:shadow-md"
                    />
                  </TooltipTrigger>
                  <TooltipContent side="left">
                    {t('ui.text.dataEntry')}
                  </TooltipContent>
                </Tooltip>
              )}

              <Card className="items-start bg-transparent py-2! shadow-none border-none">
                <CardContent className="pl-3 pr-4">
                  <div className="flex items-center gap-4 cursor-pointer">
                    <div
                      style={{ color: prefs?.color }}
                      className="size-10 rounded-lg flex items-center justify-center bg-current shrink-0"
                    >
                      {Icon && <Icon className="size-6 text-zinc-50" />}
                    </div>
                    <div className="flex flex-1 flex-col items-start justify-start gap-0">
                      <div
                        className="capitalize text-sm"
                        style={{ color: prefs?.color }}
                      >
                        {prefs?.labelKey ? t(prefs.labelKey) : data?.category}
                      </div>
                      <div className="text-sm">{name ? t(name) : null}</div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Çıkış Noktası */}
              <Tooltip>
                <TooltipTrigger asChild>
                  <Handle
                    type="source"
                    position={LAYOUT[position][1]}
                    className="border-ring hover:shadow-md"
                  />
                </TooltipTrigger>
                <TooltipContent side="right">
                  {t('ui.text.dataOutput')}
                </TooltipContent>
              </Tooltip>
            </div>
          </div>
        </ContextMenuTrigger>
        <ContextMenuContent className="max-h-110 pb-2">
          <ContextMenuItem onClick={() => setNameOpen(true)}>
            {t('ui.text.rename')}
          </ContextMenuItem>
          {!data?.hide_input && (
            <>
              <ContextMenuItem
                variant="destructive"
                onClick={() => deleteNode()}
              >
                <span className="text-destructive">{t('ui.text.delete')}</span>
              </ContextMenuItem>
              <ContextMenuSub>
                <ContextMenuSubTrigger>
                  {t('ui.text.connectionPoints')}
                </ContextMenuSubTrigger>
                <ContextMenuSubContent className="w-44">
                  {renderPositionItems}
                </ContextMenuSubContent>
              </ContextMenuSub>
            </>
          )}
        </ContextMenuContent>
      </ContextMenu>
      <Dialog open={nameOpen} onOpenChange={setNameOpen}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>{t('ui.text.renameNode')}</DialogTitle>
            <DialogDescription>{t('ui.text.addCustomName')}</DialogDescription>
          </DialogHeader>
          <div className="grid gap-4">
            <div className="grid gap-3">
              <Label htmlFor={`f-title`}>{t('ui.text.title')}</Label>
              <Input
                id={`f-title`}
                name={`f-title`}
                defaultValue={t(name)}
                onChange={(e) => setName(e.target.value)}
              />
            </div>
          </div>
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
                  updateField('title', name)
                }}
              >
                {t('ui.text.save')}
              </Button>
            </DialogClose>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  )
})
CoreNode.displayName = 'CoreNode'

export default CoreNode
