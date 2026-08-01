export const bc = new BroadcastChannel('credential-oauth')

export function hideAuthDialog() {
  bc.postMessage('hideAuth')
}
