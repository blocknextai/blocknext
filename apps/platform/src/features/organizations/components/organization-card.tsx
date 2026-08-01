import { memo } from 'react'
import { Link } from 'react-router'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useTranslation } from 'react-i18next'
import { Eye, Calendar, ArrowRight } from 'lucide-react'
import { CompanyIcon } from '@/components/shared/custom-icons'
import { VerifiedBadge } from '@/components/shared/verified-tooltip'
import TimeAgoI18n from '@/components/shared/timeagoi18'

const OrganizationCard = memo(({ organization }) => {
  const { t } = useTranslation()

  return (
    <Card
      key={organization.id}
      className="group relative overflow-hidden hover:shadow-xl transition-all duration-300 hover:scale-[1.02] border-0 bg-linear-to-br from-card to-card/50 backdrop-blur-sm"
    >
      <CardHeader className="relative z-10 pb-3">
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-3">
            <div className="relative rounded-xl overflow-hidden">
              <Link to={`/organizations/${organization.id}`}>
                <div className="flex h-14 w-14 items-center justify-center rounded-xl bg-linear-to-br from-primary to-primary/80 text-primary-foreground shadow-lg group-hover:shadow-xl transition-all duration-300 group-hover:scale-110">
                  <CompanyIcon
                    seed={organization.id}
                    className="size-7 shrink-0"
                    size={40}
                    avatarColor="fafafa"
                  />
                </div>
              </Link>
            </div>
            <div className="flex-1 min-w-0">
              <CardTitle className="text-lg font-semibold text-foreground group-hover:text-primary transition-colors duration-300 flex flex-row gap-1 items-center">
                <Link
                  to={`/organizations/${organization.id}`}
                  className="hover:underline"
                >
                  {organization.title.length > 20
                    ? `${organization.title.substring(0, 20)}...`
                    : organization.title}
                </Link>
                {organization.isVerified && (
                  <VerifiedBadge className="dark:text-zinc-900 dark:fill-amber-400 fill-amber-700 text-zinc-100" />
                )}
              </CardTitle>
            </div>
          </div>
        </div>
      </CardHeader>
      <CardContent className="relative z-10 pt-0">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2 text-xs text-muted-foreground">
            <Calendar className="h-3 w-3" />
            <span>
              {t('ui.text.updated')}{' '}
              <TimeAgoI18n date={organization.updatedAt} />
            </span>
          </div>
          <Link to={`/organizations/${organization.id}`} className="block">
            <Button
              className="group/btn relative overflow-hidden bg-linear-to-r from-primary to-primary/90 hover:from-primary/90 hover:to-primary text-primary-foreground transition-all duration-300 hover:scale-105"
              size="sm"
            >
              <span className="relative z-10 flex items-center gap-2">
                <Eye className="h-4 w-4" />
                {t('ui.text.viewFlows')}
              </span>
              <ArrowRight className="h-4 w-4 ml-1 group-hover/btn:translate-x-1 transition-transform duration-300" />
            </Button>
          </Link>
        </div>
      </CardContent>
    </Card>
  )
})

OrganizationCard.displayName = 'OrganizationCard'

export { OrganizationCard }
