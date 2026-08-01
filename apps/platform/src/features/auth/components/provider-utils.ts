const providerNames = {
  google: 'Google',
  github: 'GitHub',
  x: 'X',
  facebook: 'Facebook',
  metamask: 'MetaMask',
  email: 'Email',
  password: 'Password',
}

export function getProviderName(provider) {
  return providerNames[provider] || provider?.toUpperCase()
}
