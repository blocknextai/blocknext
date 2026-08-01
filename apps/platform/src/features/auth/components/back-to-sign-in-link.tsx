import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { ChevronLeft } from 'lucide-react'

type Props =
  { to: string; onClick?: never } | { onClick: () => void; to?: never }

const className =
  'inline-flex items-center gap-1 text-sm text-primary hover:underline mb-4 -mt-2 self-start'

const BackToSignInLink = (props: Props) => {
  const { t } = useTranslation()
  const label = (
    <>
      <ChevronLeft className="size-4" />
      {t('ui.text.back_to_sign_in')}
    </>
  )

  if ('to' in props && props.to) {
    return (
      <Link to={props.to} className={className}>
        {label}
      </Link>
    )
  }

  return (
    <button type="button" onClick={props.onClick} className={className}>
      {label}
    </button>
  )
}

BackToSignInLink.displayName = 'BackToSignInLink'

export { BackToSignInLink }
