import useSWRImmutable from 'swr/immutable'
import platformService from '@/features/platform/services/platform'
import { unwrap } from '@/lib/swr'

interface PlatformFeatures {
  functionCalling?: boolean
  workflowsGeneration?: boolean
}

export function usePlatformFeatures() {
  const { data, error, isLoading } = useSWRImmutable('platform-features', () =>
    unwrap<PlatformFeatures>(platformService.getFeatures()),
  )
  const features = data ?? {}

  return {
    functionCallingEnabled: features.functionCalling === true,
    workflowsGenerationEnabled: features.workflowsGeneration === true,
    isLoading,
    error,
  }
}
