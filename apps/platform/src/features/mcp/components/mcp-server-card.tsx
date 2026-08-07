import { memo, useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { Check, Copy, KeyRound, Server, Wrench } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter } from '@/components/ui/card'
import { useIconResolver } from '@/features/flow-editor/icons'
import { useCopy } from '@/hooks/use-copy'
import { cn } from '@/lib/utils'
import type { McpServer } from '@/features/mcp/services/mcp'

type McpServerCardProps = {
  server: McpServer
  onClick?: (server: McpServer) => void
}

const McpServerCard = memo(({ server, onClick }: McpServerCardProps) => {
  const { copied, copy } = useCopy()
  const resolveIcon = useIconResolver()
  const BrandIcon = resolveIcon(server.icon)
  const { t } = useTranslation()

  const toolCount = server.tools?.length ?? 0

  const credentialCount = useMemo(() => {
    const credentials = new Set<string>()
    for (const tool of server.tools ?? []) {
      for (const cred of tool.supportedCredentials ?? []) {
        credentials.add(cred)
      }
    }
    return credentials.size
  }, [server.tools])

  return (
    <Card
      onClick={() => onClick?.(server)}
      className={cn(
        'group/mcp h-full justify-between transition-all duration-300 ease-in-out hover:shadow-lg hover:ring-foreground/20',
        onClick && 'cursor-pointer hover:-translate-y-0.5',
      )}
    >
      <CardContent className="flex flex-col gap-3">
        <div className="flex items-start gap-3">
          <div className="bg-muted ring-foreground/10 relative flex size-11 shrink-0 items-center justify-center overflow-hidden rounded-lg ring-1">
            {BrandIcon ? (
              <BrandIcon className="size-7" />
            ) : (
              <Server className="text-muted-foreground size-5" />
            )}
          </div>

          <div className="flex min-w-0 flex-1 flex-col gap-1">
            <div className="flex flex-wrap items-center gap-2">
              <span className="truncate font-heading text-base font-medium leading-snug">
                {server.name}
              </span>
              <Badge variant="outline" className="font-mono text-[10px]">
                v{server.version}
              </Badge>
            </div>
            {server.description && (
              <p className="text-muted-foreground text-sm">
                {server.description}
              </p>
            )}
          </div>
        </div>
      </CardContent>

      <CardFooter className="flex-wrap items-center gap-2">
        <Badge variant="secondary" className="gap-1">
          <Wrench className="size-3" />
          {t('ui.text.mcpToolCount', { count: toolCount })}
        </Badge>
        {credentialCount > 0 && (
          <Badge variant="secondary" className="gap-1">
            <KeyRound className="size-3" />
            {credentialCount}
          </Badge>
        )}
        {server.url && (
          <Button
            variant="ghost"
            size="sm"
            className="ml-auto gap-1.5"
            onClick={(event) => {
              event.stopPropagation()
              copy(server.url!)
            }}
          >
            {copied ? (
              <Check className="size-3.5" />
            ) : (
              <Copy className="size-3.5" />
            )}
            {t('ui.text.mcpCopyUrl')}
          </Button>
        )}
      </CardFooter>
    </Card>
  )
})

McpServerCard.displayName = 'McpServerCard'

export { McpServerCard }
