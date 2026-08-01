import { Component } from 'react'
import type { ErrorInfo, ReactNode } from 'react'
import i18next from 'i18next'
import { Button } from '@/components/ui/button'

interface Props {
  children: ReactNode
  fallback?: ReactNode | ((reset: () => void) => ReactNode)
  onReset?: () => void
  scope?: 'root' | 'section'
}

interface State {
  hasError: boolean
}

class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = { hasError: false }
  }

  static getDerivedStateFromError(): State {
    return { hasError: true }
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('ErrorBoundary caught an error:', error, errorInfo)
  }

  handleReset = () => {
    this.setState({ hasError: false })
    if (this.props.onReset) {
      this.props.onReset()
    } else if (this.props.scope !== 'section') {
      window.location.href = '/'
    }
  }

  renderSectionFallback() {
    return (
      <div className="flex flex-col items-center justify-center p-8 text-center gap-4 border border-destructive/20 bg-destructive/5 rounded-lg m-4">
        <div className="space-y-1">
          <h3 className="text-lg font-semibold text-foreground">
            {i18next.t('ui.text.errorTitle')}
          </h3>
          <p className="text-sm text-muted-foreground">
            {i18next.t('ui.text.sectionFailedToLoadDescription')}
          </p>
        </div>
        <Button onClick={this.handleReset}>
          {i18next.t('ui.text.tryAgain')}
        </Button>
      </div>
    )
  }

  renderRootFallback() {
    return (
      <div className="min-h-screen bg-linear-to-br from-background via-background to-muted/20 flex items-center justify-center p-4">
        <div className="w-full max-w-2xl mx-auto text-center space-y-8">
          <div className="relative">
            <h1 className="text-8xl md:text-9xl font-black bg-linear-to-r from-destructive via-destructive/80 to-destructive/60 bg-clip-text text-transparent">
              {i18next.t('ui.text.oops')}
            </h1>
            <div className="absolute inset-0 bg-linear-to-r from-destructive/20 via-transparent to-destructive/20 blur-3xl"></div>
          </div>

          <div className="space-y-4">
            <h2 className="text-3xl md:text-4xl font-bold text-foreground">
              {i18next.t('ui.text.errorTitle')}
            </h2>
            <p className="text-lg text-muted-foreground max-w-md mx-auto leading-relaxed">
              {i18next.t('ui.text.somethingWentWrong')}
            </p>
          </div>

          <div className="flex justify-center items-center pt-6">
            <Button size="lg" onClick={this.handleReset}>
              {i18next.t('ui.text.goToHome')}
            </Button>
          </div>
        </div>
      </div>
    )
  }

  render() {
    if (!this.state.hasError) {
      return this.props.children
    }

    if (typeof this.props.fallback === 'function') {
      return this.props.fallback(this.handleReset)
    }
    if (this.props.fallback !== undefined) {
      return this.props.fallback
    }
    if (this.props.scope === 'section') {
      return this.renderSectionFallback()
    }
    return this.renderRootFallback()
  }
}

export { ErrorBoundary }
