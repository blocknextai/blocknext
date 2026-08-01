import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import TimeAgoI18n from '@/components/shared/timeagoi18'
import { Key, RefreshCw, Trash2 } from 'lucide-react'

const ApiKeyDetailSheet = ({
  apiKey,
  onOpenChange,
  onRegenerate,
  onDelete,
}) => {
  const { t } = useTranslation()

  return (
    <Sheet open={!!apiKey} onOpenChange={onOpenChange}>
      <SheetContent className="sm:max-w-lg overflow-y-auto gap-0">
        <SheetHeader>
          <SheetTitle className="flex items-center gap-2">
            <Key className="size-4 text-icon" />
            {apiKey?.name}
          </SheetTitle>
          <SheetDescription>
            {t(
              'ui.text.apiKeyDetailsDescription',
              'Review and manage this API key.',
            )}
          </SheetDescription>
        </SheetHeader>

        {apiKey && (
          <div className="flex flex-col gap-6 px-4 pb-6">
            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-1.5">
                <Label className="text-muted-foreground">
                  {t('ui.text.lastUsed')}
                </Label>
                <span className="text-sm">
                  {apiKey.lastUsedAt ? (
                    <TimeAgoI18n date={apiKey.lastUsedAt} />
                  ) : (
                    <span className="text-muted-foreground">
                      {t('ui.text.neverUsed')}
                    </span>
                  )}
                </span>
              </div>
              <div className="grid gap-1.5">
                <Label className="text-muted-foreground">
                  {t('ui.text.created')}
                </Label>
                <span className="text-sm">
                  <TimeAgoI18n date={apiKey.createdAt} />
                </span>
              </div>
            </div>

            {!!apiKey.scopes?.length && (
              <div className="grid gap-1.5">
                <Label className="text-muted-foreground">
                  {t('ui.text.scopes', 'Scopes')}
                </Label>
                <div className="flex flex-wrap gap-1.5">
                  {apiKey.scopes.map((scope: string) => (
                    <Badge
                      key={scope}
                      variant="secondary"
                      className="font-mono"
                    >
                      {scope}
                    </Badge>
                  ))}
                </div>
              </div>
            )}

            <div className="flex flex-col gap-2 pt-4 border-t">
              <Button
                variant="outline"
                className="justify-start gap-2"
                onClick={() => onRegenerate(apiKey)}
              >
                <RefreshCw className="size-4" />
                {t('ui.text.revokeAndRegenerate', 'Revoke & Regenerate')}
              </Button>
              <p className="text-xs text-muted-foreground px-1">
                {t(
                  'ui.text.revokeAndRegenerateHint',
                  'Invalidates the current key and creates a new one with the same settings.',
                )}
              </p>

              <Button
                variant="ghost"
                className="justify-start gap-2 text-destructive hover:text-destructive hover:bg-destructive/10 mt-2"
                onClick={() => onDelete(apiKey)}
              >
                <Trash2 className="size-4" />
                {t('ui.text.delete')}
              </Button>
            </div>
          </div>
        )}
      </SheetContent>
    </Sheet>
  )
}

ApiKeyDetailSheet.displayName = 'ApiKeyDetailSheet'

export { ApiKeyDetailSheet }
