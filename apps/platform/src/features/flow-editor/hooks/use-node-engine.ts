import useSWR from 'swr'
import nodeEngineService from '@/features/flow-editor/services/node-engine'
import { unwrap } from '@/lib/swr'

export function useNodeEngineCredentials() {
  const { data, error, isLoading, mutate } = useSWR(
    'node-engine-credentials',
    () => unwrap(nodeEngineService.getCredentials()),
  )
  return {
    credentials: (data ?? []) as any[],
    isLoading,
    error,
    mutate,
  }
}

export function useNodeEngineNodes() {
  const { data, error, isLoading } = useSWR('node-engine-nodes', () =>
    unwrap(nodeEngineService.getNodes()),
  )
  return {
    nodes: (data ?? []) as any[],
    isLoading,
    error,
  }
}

export function useNodeEngineWebhookSources() {
  const { data, error, isLoading } = useSWR('node-engine-webhook-sources', () =>
    unwrap(nodeEngineService.getWebhookSources()),
  )
  return {
    sources: (data ?? []) as any[],
    isLoading,
    error,
  }
}

export function useNodeEngineTriggerVariables() {
  const { data, error, isLoading } = useSWR(
    'node-engine-trigger-variables',
    () => unwrap(nodeEngineService.getTriggerVariables()),
  )
  return {
    variables: (data ?? []) as any[],
    isLoading,
    error,
  }
}
