import { PlatformLogo } from '@/features/navigation/components/logo'
import { Home } from 'lucide-react'
import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'

const NotFoundPage = () => {
  const { t } = useTranslation()

  return (
    <div className="min-h-screen bg-linear-to-br from-background via-background to-muted/20 flex items-center justify-center p-4">
      <div className="w-full max-w-2xl mx-auto text-center space-y-8">
        <div className="flex justify-center mb-8">
          <Link to="/">
            <PlatformLogo width={180} height={36} />
          </Link>
        </div>

        <div className="relative">
          <h1 className="text-8xl md:text-9xl font-black bg-linear-to-r from-primary via-primary/80 to-primary/60 bg-clip-text text-transparent">
            404
          </h1>
          <div className="absolute inset-0 bg-linear-to-r from-primary/20 via-transparent to-primary/20 blur-3xl"></div>
        </div>

        <div className="space-y-4">
          <h2 className="text-3xl md:text-4xl font-bold text-foreground">
            {t('ui.text.pageNotFound')}
          </h2>
          <p className="text-lg text-muted-foreground max-w-md mx-auto leading-relaxed">
            {t('ui.text.pageNotFoundDescription')}
          </p>
        </div>

        <div className="flex justify-center items-center pt-6">
          <Link
            to="/"
            className="inline-flex items-center gap-2 px-6 py-3 bg-primary hover:bg-primary/90 text-primary-foreground rounded-md font-medium transition-colors"
          >
            <Home className="w-4 h-4" />
            {t('ui.text.goToHome')}
          </Link>
        </div>
      </div>
    </div>
  )
}

export default NotFoundPage
