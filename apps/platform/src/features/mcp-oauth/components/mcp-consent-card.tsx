import { useTranslation } from 'react-i18next'
import { Check, Plug, ShieldCheck, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { NativeSelect, NativeSelectOption } from '@/components/ui/native-select'
import { Separator } from '@/components/ui/separator'
import type { McpAuthorizationRequest } from '@/features/mcp-oauth/services/mcp-oauth'

type Organization = {
  id: string
  title: string
}

type McpConsentCardProps = {
  authorizationRequest: McpAuthorizationRequest
  organizations: Organization[]
  organizationId: string
  isSubmitting: boolean
  onOrganizationChange: (organizationId: string) => void
  onApprove: () => void
  onDeny: () => void
}

export function McpConsentCard({
  authorizationRequest,
  organizations,
  organizationId,
  isSubmitting,
  onOrganizationChange,
  onApprove,
  onDeny,
}: McpConsentCardProps) {
  const { t } = useTranslation()
  const { client, scopes } = authorizationRequest

  return (
    <Card className="w-full max-w-md">
      <CardHeader>
        <div className="mb-2 flex size-10 items-center justify-center rounded-lg bg-muted">
          {client.logoUri ? (
            <img
              src={client.logoUri}
              alt={client.name}
              className="size-6 rounded"
            />
          ) : (
            <Plug className="size-5 text-muted-foreground" />
          )}
        </div>
        <CardTitle>
          {t('ui.text.mcpConsentTitle', { client: client.name })}
        </CardTitle>
        <CardDescription>{t('ui.text.mcpConsentDescription')}</CardDescription>
      </CardHeader>

      <CardContent className="flex flex-col gap-4">
        <div className="flex flex-col gap-2">
          <Label htmlFor="mcp-consent-organization">
            {t('ui.text.organization')}
          </Label>
          <NativeSelect
            id="mcp-consent-organization"
            className="w-full"
            value={organizationId}
            disabled={isSubmitting}
            onChange={(event) => onOrganizationChange(event.target.value)}
          >
            {organizations.map((organization) => (
              <NativeSelectOption key={organization.id} value={organization.id}>
                {organization.title}
              </NativeSelectOption>
            ))}
          </NativeSelect>
          <p className="text-xs text-muted-foreground">
            {t('ui.text.mcpConsentOrganizationHint')}
          </p>
        </div>

        <Separator />

        <div className="flex flex-col gap-3">
          <div className="flex items-center gap-2 text-sm font-medium">
            <ShieldCheck className="size-4 text-muted-foreground" />
            {t('ui.text.mcpConsentPermissions')}
          </div>
          <ul className="flex flex-col gap-2">
            {scopes.map((scope) => (
              <li key={scope.scope} className="flex items-start gap-2 text-sm">
                <Check className="mt-0.5 size-4 shrink-0 text-muted-foreground" />
                <span>
                  {scope.description || scope.scope}
                  <span className="ml-1 text-xs text-muted-foreground">
                    ({scope.scope})
                  </span>
                </span>
              </li>
            ))}
          </ul>
        </div>

        <p className="text-xs text-muted-foreground break-all">
          {t('ui.text.mcpConsentRedirectHint', {
            redirectUri: authorizationRequest.redirectUri,
          })}
        </p>
      </CardContent>

      <CardFooter className="flex gap-2">
        <Button
          variant="outline"
          className="flex-1"
          disabled={isSubmitting}
          onClick={onDeny}
        >
          <X className="size-4" />
          {t('ui.text.deny')}
        </Button>
        <Button
          className="flex-1"
          disabled={isSubmitting || !organizationId}
          onClick={onApprove}
        >
          <Check className="size-4" />
          {t('ui.text.allow')}
        </Button>
      </CardFooter>
    </Card>
  )
}
