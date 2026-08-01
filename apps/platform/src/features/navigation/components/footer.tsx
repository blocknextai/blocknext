import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import { landingUrl, SOCIAL_LINKS } from '@/constants'
import { Image } from '@/components/shared/image'

const Footer = (props) => {
  const { t } = useTranslation()
  const startYear = 2025
  const currentYear = new Date().getFullYear()
  const year =
    startYear === currentYear ? startYear : `${startYear}-${currentYear}`

  return (
    <footer className="py-6 border-t bg-background mt-auto w-full">
      <div className="container max-w-7xl mx-auto px-4">
        {props.children && <div>{props.children}</div>}
        <div className="flex flex-col md:flex-row items-center justify-between gap-4">
          {/* Copyright */}
          <p className="text-sm text-muted-foreground text-center md:text-left">
            © {year}{' '}
            <Link
              to={`${landingUrl}`}
              className="hover:underline"
              target="_blank"
              rel="noopener noreferrer"
            >
              BlockNext
            </Link>
            . {t('ui.text.allRightsReserved')}
          </p>

          {/* Legal Links & Social Media */}
          <div className="flex flex-col md:flex-row items-center gap-4 md:gap-6">
            {/* Legal Links */}
            <div className="flex gap-4">
              <a
                href={`${landingUrl}/privacy-policy`}
                target="_blank"
                rel="noopener noreferrer"
                className="text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                {t('ui.text.privacyPolicy')}
              </a>
              <a
                href={`${landingUrl}/terms-and-conditions`}
                target="_blank"
                rel="noopener noreferrer"
                className="text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                {t('ui.text.termsAndConditions')}
              </a>
            </div>

            {/* Social Media */}
            {SOCIAL_LINKS.length > 0 && (
              <div className="flex gap-3">
                {SOCIAL_LINKS.map((link) => (
                  <Link
                    key={link.href}
                    to={link.href}
                    target="_blank"
                    className="p-1"
                    rel="noopener nofollow noreferrer"
                  >
                    <Image
                      src={link.icon}
                      alt={link.label}
                      title={link.label}
                      width={24}
                      height={24}
                      className="transition hover:invert hover:brightness-0 hover:text-zinc-50 cursor-pointer"
                    />
                  </Link>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </footer>
  )
}

Footer.displayName = 'Footer'

export { Footer }
