import { useTranslation } from 'react-i18next'
import { KeyRound, Wrench } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { CopyField } from '@/components/shared/copy-field'
import { CodeBlock } from '@/components/shared/code-block'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useIconResolver } from '@/features/flow-editor/icons'
import { McpToolSchemaSection } from '@/features/mcp/components/mcp-tool-schema'
import type { IconSource } from '@/features/flow-editor/types'
import type { McpServer } from '@/features/mcp/services/mcp'

const API_KEY_PLACEHOLDER = 'YOUR_API_KEY'

const configSnippet = (name: string, url: string) =>
  JSON.stringify(
    {
      mcpServers: {
        [`blocknext-${name}`]: {
          type: 'http',
          url,
          headers: { 'X-API-Key': API_KEY_PLACEHOLDER },
        },
      },
    },
    null,
    2,
  )

const curlSnippet = (url: string) =>
  `curl --location '${url}' \\\n` +
  `--header 'X-API-Key: ${API_KEY_PLACEHOLDER}' \\\n` +
  `--header 'Content-Type: application/json' \\\n` +
  `--header 'Accept: application/json, text/event-stream' \\\n` +
  `--data '{"jsonrpc":"2.0","id":1,"method":"tools/list"}'`

type McpServerDetailDialogProps = {
  server: McpServer | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

const McpToolIcon = ({ icon }: { icon?: IconSource }) => {
  const resolveIcon = useIconResolver()
  const Icon = resolveIcon(icon)

  return (
    <div className="bg-muted ring-foreground/10 flex size-12 shrink-0 items-center justify-center overflow-hidden rounded-lg ring-1">
      {Icon ? (
        <Icon className="size-9" />
      ) : (
        <Wrench className="text-muted-foreground size-5" />
      )}
    </div>
  )
}

const McpServerDetailDialog = ({
  server,
  open,
  onOpenChange,
}: McpServerDetailDialogProps) => {
  const { t } = useTranslation()

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="flex max-h-[85vh] flex-col overflow-y-auto sm:max-w-5xl md:overflow-hidden">
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

            <div className="flex flex-col gap-4 md:min-h-0 md:flex-1 md:flex-row">
              <div className="flex flex-col gap-2 md:min-h-0 md:flex-1">
                <h3 className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                  {t('ui.text.tools')} ({server.tools?.length ?? 0})
                </h3>

                {!server.tools?.length ? (
                  <p className="text-sm text-muted-foreground">—</p>
                ) : (
                  <div className="flex flex-col gap-2 md:overflow-y-auto md:pr-1">
                    {server.tools.map((tool) => (
                      <div
                        key={tool.id}
                        className="flex items-start gap-2.5 rounded-lg border bg-card p-3 ring-1 ring-foreground/5 transition-colors hover:border-foreground/20"
                      >
                        <McpToolIcon icon={tool.icon} />
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

              {server.url && (
                <div className="flex shrink-0 flex-col gap-2 md:sticky md:top-0 md:w-96 md:self-start">
                  <Tabs defaultValue="config" className="gap-3">
                    <TabsList>
                      <TabsTrigger value="config">
                        {t('ui.text.mcpClientConfig')}
                      </TabsTrigger>
                      <TabsTrigger value="curl">cURL</TabsTrigger>
                    </TabsList>
                    <TabsContent value="config">
                      <CodeBlock code={configSnippet(server.id, server.url)} />
                    </TabsContent>
                    <TabsContent value="curl">
                      <CodeBlock code={curlSnippet(server.url)} />
                    </TabsContent>
                  </Tabs>
                  <p className="text-xs text-muted-foreground">
                    {t('ui.text.mcpApiKeyHint')}
                  </p>
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
