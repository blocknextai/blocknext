import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { Check, Copy } from 'lucide-react'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { useCopy } from '@/hooks/use-copy'
import { config } from '@/lib/config'

const API_KEY_PLACEHOLDER = 'YOUR_API_KEY'
const CREDENTIAL_PLACEHOLDER = '<credential-ui-key>'
const RUNTIME_PROMPT_PLACEHOLDER = 'YOUR_RUNTIME_PROMPT'
const RUNTIME_INSTRUCTION_PLACEHOLDER = 'YOUR_RUNTIME_INSTRUCTION'

function CodeBlock({ code }: { code: string }) {
  const { t } = useTranslation()
  const { copied, copy } = useCopy()
  const label = t('ui.text.copy', 'Copy')

  return (
    <div className="overflow-hidden rounded-lg border bg-muted/50">
      <div className="flex items-center justify-end border-b bg-muted/50 px-2 py-1">
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={() => copy(code)}
          aria-label={label}
          title={label}
          className="h-6 gap-1.5 px-2 text-xs text-muted-foreground"
        >
          {copied ? (
            <Check className="size-3.5 text-green-500" />
          ) : (
            <Copy className="size-3.5" />
          )}
          {label}
        </Button>
      </div>
      <pre className="max-h-[55vh] overflow-auto p-3 font-mono text-xs leading-relaxed">
        <code>{code}</code>
      </pre>
    </div>
  )
}

function buildRequestBody(nodes: any[], apiNodes: any[]) {
  const credentialKeyFor = (nodeId: string) =>
    apiNodes.find((n) => n.id === nodeId)?.supportedCredentials?.[0] ?? null

  const runners = (nodes ?? []).map((node) => {
    const runner: Record<string, any> = {
      id: node.id,
      nodeId: node.nodeId,
      runtimeInstruction: node.instruction || RUNTIME_INSTRUCTION_PLACEHOLDER,
    }
    const credentialKey = credentialKeyFor(node.nodeId)
    if (credentialKey) {
      runner.credentials = { [credentialKey]: CREDENTIAL_PLACEHOLDER }
    }
    return runner
  })

  return { runtimePrompt: RUNTIME_PROMPT_PLACEHOLDER, nodes: runners }
}

export function FlowApiSheet({ open, onOpenChange, flowId, nodes, apiNodes }) {
  const { t } = useTranslation()

  const endpoint = `${config.platformApiUrl}/task-runner/trigger/${flowId}`

  const bodyJson = useMemo(
    () => JSON.stringify(buildRequestBody(nodes, apiNodes), null, 4),
    [nodes, apiNodes],
  )

  const curlSnippet = useMemo(
    () =>
      `curl --location '${endpoint}' \\\n` +
      `--header 'x-api-key: ${API_KEY_PLACEHOLDER}' \\\n` +
      `--header 'Content-Type: application/json' \\\n` +
      `--data '${bodyJson}'`,
    [endpoint, bodyJson],
  )

  const jsFetchSnippet = useMemo(
    () =>
      `const response = await fetch('${endpoint}', {\n` +
      `  method: 'POST',\n` +
      `  headers: {\n` +
      `    'x-api-key': '${API_KEY_PLACEHOLDER}',\n` +
      `    'Content-Type': 'application/json',\n` +
      `  },\n` +
      `  body: JSON.stringify(${bodyJson.replace(/\n/g, '\n  ')}),\n` +
      `})\n\n` +
      `const data = await response.json()`,
    [endpoint, bodyJson],
  )

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="sm:max-w-3xl">
        <SheetHeader>
          <SheetTitle>{t('ui.text.apiTrigger', 'API Trigger')}</SheetTitle>
          <SheetDescription>
            {t(
              'ui.text.apiTriggerDescription',
              'Trigger this flow from your own code using the API.',
            )}
          </SheetDescription>
        </SheetHeader>

        <div className="flex flex-1 flex-col gap-5 overflow-y-auto px-4 pb-4">
          <Tabs defaultValue="curl" className="gap-3">
            <TabsList>
              <TabsTrigger value="curl">cURL</TabsTrigger>
              <TabsTrigger value="js-fetch">JavaScript - Fetch</TabsTrigger>
            </TabsList>
            <TabsContent value="curl">
              <CodeBlock code={curlSnippet} />
            </TabsContent>
            <TabsContent value="js-fetch">
              <CodeBlock code={jsFetchSnippet} />
            </TabsContent>
          </Tabs>
        </div>
      </SheetContent>
    </Sheet>
  )
}

FlowApiSheet.displayName = 'FlowApiSheet'
