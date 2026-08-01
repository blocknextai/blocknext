import { createApiClient } from '@/lib/api-client'
import { config } from '@/lib/config'

export default createApiClient(config.mcpApiUrl)
