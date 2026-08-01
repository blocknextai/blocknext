import { createRoot } from 'react-dom/client'
import { SWRConfig } from 'swr'
import Routes from '@/routes'
import { LocalToaster } from '@/components/shared/local-toaster'
import { applyStoredTheme } from '@/stores/theme-store'
import { swrConfig } from '@/lib/swr'
import '@/lib/i18n'
import '@/styles/index.css'

applyStoredTheme()

createRoot(document.getElementById('root')).render(
  <SWRConfig value={swrConfig}>
    <LocalToaster />
    <Routes />
  </SWRConfig>,
)
