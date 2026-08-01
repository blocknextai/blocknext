import { useTranslation } from 'react-i18next'
import { KeyRound } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { CopyField } from '@/components/shared/copy-field'
import { FlowIcon } from '@/components/shared/custom-icons'
import { McpToolSchemaSection } from '@/features/mcp/components/mcp-tool-schema'
import type { McpServer } from '@/features/mcp/services/mcp'

type McpServerDetailDialogProps = {
  server: McpServer | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

const McpServerDetailDialog = ({
  server,
  open,
  onOpenChange,
}: McpServerDetailDialogProps) => {
  const { t } = useTranslation()

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl max-h-[80vh] overflow-y-auto">
        {server && (
          <>
            <DialogHeader>
              <div className="flex flex-wrap items-center gap-2">
                <DialogTitle>{server.name}</DialogTitle>
                <Badge variant="outline" className="font-mono text-[10px]">
                  v{server.version}
                </Badge>
              </div>
              {server.description && (
                <DialogDescription>{server.description}</DialogDescription>
              )}
            </DialogHeader>

            {server.url && (
              <div className="flex flex-col gap-2">
                <h3 className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                  {t('ui.text.mcpConnectionUrl')}
                </h3>
                <CopyField value={server.url} inputClassName="bg-muted/40" />
              </div>
            )}

            <div className="flex flex-col gap-2">
              <h3 className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                {t('ui.text.tools')} ({server.tools?.length ?? 0})
              </h3>

              {!server.tools?.length ? (
                <p className="text-sm text-muted-foreground">—</p>
              ) : (
                <div className="flex flex-col gap-2">
                  {server.tools.map((tool) => (
                    <div
                      key={tool.id}
                      className="flex items-start gap-2.5 rounded-lg border bg-card p-3 ring-1 ring-foreground/5 transition-colors hover:border-foreground/20"
                    >
                      <div
                        className="size-8 shrink-0 overflow-hidden rounded-md ring-1 ring-foreground/10"
                        style={{
                          backgroundImage: FlowIcon(tool.id),
                          backgroundSize: 'cover',
                          backgroundPosition: 'center',
                        }}
                      />
                      <div className="flex min-w-0 flex-1 flex-col gap-1">
                        <div className="flex flex-wrap items-center gap-1.5">
                          <span className="truncate font-mono text-sm font-medium">
                            {tool.name}
                          </span>
                          <Badge
                            variant="outline"
                            className="font-mono text-[10px]"
                          >
                            v{tool.version}
                          </Badge>
                        </div>
                        {tool.description && (
                          <p className="line-clamp-2 text-xs text-muted-foreground">
                            {tool.description}
                          </p>
                        )}
                        <div className="flex flex-wrap items-center gap-1.5 pt-0.5">
                          {tool.supportedCredentials?.map((cred) => (
                            <Badge
                              key={cred}
                              variant="outline"
                              className="gap-1 text-[10px] font-normal"
                            >
                              <KeyRound className="size-3 text-muted-foreground" />
                              {cred}
                            </Badge>
                          ))}
                        </div>

                        <div className="flex flex-col gap-1.5 pt-1">
                          <McpToolSchemaSection
                            label={t('ui.text.inputs')}
                            schema={tool.inputSchema}
                          />
                          <McpToolSchemaSection
                            label={t('ui.text.outputs')}
                            schema={tool.outputSchema}
                          />
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </>
        )}
      </DialogContent>
    </Dialog>
  )
}

McpServerDetailDialog.displayName = 'McpServerDetailDialog'

export { McpServerDetailDialog }
