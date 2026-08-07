import FlowCredentials from '@/features/flow-editor/components/flow-credentials'
import FlowRunAs from '@/features/flow-editor/components/flow-run-as'
import FlowTriggerTypeSelector from '@/features/flow-editor/components/flow-trigger-type-selector'
import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router'
import { useTranslation } from 'react-i18next'
import { Stepper } from '@/components/ui/stepper'
import { Loading } from '@/components/shared/loading'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Copy,
  Check,
  MousePointerClick,
  ClockFading,
  Webhook,
} from 'lucide-react'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  credentialsService,
  credentialOAuthService,
} from '@/features/credentials'
import { organizationsService } from '@/features/organizations'
import { taskRunnerService, workflowsService } from '@/features/workflows'
import { FlowIcon } from '@/components/shared/custom-icons'

const triggerTypeConfig = {
  manual: { icon: MousePointerClick, label: 'ui.text.manual' },
  schedule: { icon: ClockFading, label: 'ui.text.schedule' },
  webhook: { icon: Webhook, label: 'ui.text.webhook' },
}

const steps = [
  { title: '1', prefs: [undefined, undefined] },
  { title: '2', prefs: [undefined, undefined] },
  {
    title: '3',
    prefs: ['ui.text.credentials', 'ui.text.selectCredentials'],
  },
]

function OrganizationRunPage() {
  const { organizationId, id } = useParams()
  const { t } = useTranslation()
  const [loading, setLoading] = useState(true)
  const [schemas, setSchemas] = useState([])
  const [userSecrets, setUserSecrets] = useState([])
  const [orgSecrets, setOrgSecrets] = useState([])
  const [flow, setFlow] = useState()
  const [currentStep, setCurrentStep] = useState(0)
  const [triggerType, setTriggerType] = useState('manual')
  const [runType, setRunType] = useState()
  const [credentials, setCredentials] = useState()
  const [prompt, setPrompt] = useState()
  const [cronString, setCronString] = useState('')
  const [webhookTokenDialog, setWebhookTokenDialog] = useState(null)
  const [copied, setCopied] = useState(false)
  const navigate = useNavigate()

  const getCredentials = async () => {
    const flowsResponse = await workflowsService.getRun(organizationId, id)
    const flows = flowsResponse.data
    const creds = flows.credentialSchemas
    const nodeIds = (flows.nodes ?? []).map((node) => node.nodeId)
    const userSecretsResponse =
      await credentialsService.getUserCredentialsByNodes(nodeIds)
    const userSecrets = userSecretsResponse.data
    const orgSecretsResponse =
      await organizationsService.getOrganizationCredentialsByNodes(
        organizationId,
        nodeIds,
      )
    const organizationSecrets = orgSecretsResponse.data

    setSchemas(creds)
    setUserSecrets(userSecrets)
    setOrgSecrets(organizationSecrets)
    setFlow(flows)
    setLoading(false)
  }

  const saveFlowCredential = async (payload) => {
    const isUserEndpoint = runType !== 'member'
    const resResponse = isUserEndpoint
      ? await credentialsService.createUserCredential(payload)
      : await organizationsService.createOrganizationCredential(
          organizationId,
          payload,
        )
    return resResponse.data
  }

  const authorizeFlowOAuth2 = async (payload) => {
    const isUserEndpoint = runType !== 'member'
    const resResponse = isUserEndpoint
      ? await credentialsService.createUserCredential(payload)
      : await organizationsService.createOrganizationCredential(
          organizationId,
          payload,
        )
    const res = resResponse.data
    const authUrlResponse = isUserEndpoint
      ? await credentialOAuthService.authorizeUser({ credentialId: res.id })
      : await credentialOAuthService.authorizeOrganization(organizationId, {
          credentialId: res.id,
        })
    const authUrl = authUrlResponse.data
    const popupWindow = window.open('', '_blank')
    if (popupWindow) {
      popupWindow.location.href = authUrl.url
    }
  }

  const handleTriggerTypeSelect = (type) => {
    setTriggerType(type)
    setCurrentStep(1)
  }

  const handleRunAsSubmit = (selectedRunType, promptData) => {
    setRunType(selectedRunType)
    setPrompt(promptData)
    setCurrentStep(2)
  }

  const formatNodes = () => {
    const nodeArray = []
    for (let i = 0; i < flow.nodes.length; i++) {
      const node = flow.nodes[i]
      let credentialSchemaId = null
      if (flow.credentialSchemas) {
        for (let j = 0; j < flow.credentialSchemas.length; j++) {
          const schema = flow.credentialSchemas[j]
          if (
            schema.supportedNodes &&
            schema.supportedNodes.includes(node.nodeId)
          ) {
            credentialSchemaId = schema.id
            break
          }
        }
      }

      if (credentials && credentialSchemaId) {
        const secret = credentials[credentialSchemaId]
        if (secret) {
          node.apiKey = secret.uiKey
        }
      }

      nodeArray.push(node)
    }
    return nodeArray
  }

  const buildRunners = (nodes) => {
    const runners = []
    for (let i = 0; i < nodes.length; i++) {
      const node = nodes[i]
      let credentialSchemaId = null
      if (flow.credentialSchemas) {
        for (let j = 0; j < flow.credentialSchemas.length; j++) {
          const schema = flow.credentialSchemas[j]
          if (
            schema.supportedNodes &&
            schema.supportedNodes.includes(node.nodeId)
          ) {
            credentialSchemaId = schema.id
            break
          }
        }
      }
      const nodeObj = {
        id: node.id,
        nodeId: node.nodeId,
        runtimeInstruction: node.runtimeInstruction || null,
      }

      if (node.apiKey && credentialSchemaId) {
        nodeObj.credentials = { [credentialSchemaId]: node.apiKey }
      }

      runners.push(nodeObj)
    }
    return runners
  }

  const handleCopy = (text) => {
    navigator.clipboard.writeText(text)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  const runTask = async () => {
    const nodes = formatNodes()
    const runners = buildRunners(nodes)

    const runData = {
      executionContext: 'workflow',
      contextItemId: flow.id,
      triggerType,
      cronPattern: cronString,
      runtimePrompt: prompt,
      nodes: runners,
    }

    const res = await taskRunnerService.trigger(organizationId, runData)
    if (res.isSuccess) {
      if (triggerType === 'webhook' && res.data?.webhookToken) {
        setWebhookTokenDialog(res.data.webhookToken)
        return
      }
      if (triggerType === 'schedule') {
        navigate(`/organizations/${organizationId}/triggers`)
      } else {
        navigate(`/organizations/${organizationId}/history`)
      }
    }
  }

  const renderTriggerBadge = () => {
    const config = triggerTypeConfig[triggerType]
    if (!config) {
      return null
    }
    const Icon = config.icon
    return (
      <Badge variant="outline" className="gap-1.5 font-normal">
        <Icon className="size-3" />
        {t(config.label)}
      </Badge>
    )
  }

  const renderFlowCard = () => {
    return (
      <div className="flex flex-col items-center gap-3 w-full">
        <div className="heading-lg">
          {steps[currentStep].prefs[0]
            ? t(steps[currentStep].prefs[0])
            : flow.title}
        </div>
        <div className="text-muted-foreground">
          {steps[currentStep].prefs[1]
            ? t(steps[currentStep].prefs[1])
            : flow.description}
        </div>
        {currentStep > 0 && renderTriggerBadge()}
      </div>
    )
  }

  const renderStep = () => {
    switch (currentStep) {
      case 0:
        return (
          <FlowTriggerTypeSelector
            onSelect={handleTriggerTypeSelect}
            selected={triggerType}
          />
        )
      case 1:
        return (
          <FlowRunAs
            onClick={handleRunAsSubmit}
            cronString={cronString}
            setCronString={setCronString}
            defaultValue={prompt}
            flowAvatar={FlowIcon(flow.id)}
            flow={flow}
            credentials={credentials}
            setFlow={setFlow}
            triggerType={triggerType}
          />
        )
      case 2:
        return (
          <FlowCredentials
            schemas={schemas}
            getCredentials={getCredentials}
            usedCredentials={{
              cred: credentials,
              use: setCredentials,
            }}
            secrets={runType === 'member' ? orgSecrets : userSecrets}
            onSaveCredential={saveFlowCredential}
            onAuthorizeOAuth2={authorizeFlowOAuth2}
          />
        )
      default:
        return null
    }
  }

  const changeStep = (step) => {
    if (step < 0) {
      step = 0
    }
    if (step === 3) {
      runTask()
      return
    }
    if (step > steps.length - 1) {
      step = steps.length - 1
    }
    setCurrentStep(step)
  }

  useEffect(() => {
    getCredentials()
  }, [])

  if (loading) {
    return <Loading />
  }

  return (
    <div className="w-full h-full flex relative flex-col gap-20 justify-start">
      <div className="w-full">&nbsp;</div>
      <div className="flex w-1/2 m-auto flex-col items-start justify-start gap-6">
        {renderFlowCard()}
        <div className="w-full gap-2 pb-12 relative">{renderStep()}</div>
      </div>
      <div className="sticky bottom-0 w-full p-4 pt-6 bg-linear-to-t from-background from-70% to-transparent to-90%">
        <Stepper
          className="items-center"
          steps={steps}
          currentStep={currentStep}
          onStepChange={changeStep}
        />
      </div>

      <Dialog
        open={!!webhookTokenDialog}
        onOpenChange={() => {
          setWebhookTokenDialog(null)
          navigate(`/organizations/${organizationId}/triggers`)
        }}
      >
        <DialogContent className="sm:max-w-[550px]">
          <DialogHeader>
            <DialogTitle>{t('ui.text.webhookCreated')}</DialogTitle>
            <DialogDescription>
              {t('ui.text.webhookCreatedDescription')}
            </DialogDescription>
          </DialogHeader>
          {webhookTokenDialog && (
            <div className="flex flex-col gap-3">
              <div className="flex gap-2">
                <Input
                  value={webhookTokenDialog}
                  readOnly
                  className="font-mono text-xs"
                />
                <Button
                  variant="outline"
                  size="icon"
                  onClick={() => handleCopy(webhookTokenDialog)}
                >
                  {copied ? (
                    <Check className="size-4" />
                  ) : (
                    <Copy className="size-4" />
                  )}
                </Button>
              </div>
              <p className="text-xs text-muted-foreground">
                {t('ui.text.webhookTokenWarning')}
              </p>
            </div>
          )}
          <DialogFooter>
            <Button
              size="sm"
              onClick={() => {
                setWebhookTokenDialog(null)
                navigate(`/organizations/${organizationId}/triggers`)
              }}
            >
              {t('ui.text.done')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

export default OrganizationRunPage
