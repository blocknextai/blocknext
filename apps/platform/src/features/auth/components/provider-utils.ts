const providerNames = {
  google: 'Google',
  github: 'GitHub',
  x: 'X',
  facebook: 'Facebook',
  email: 'Email',
  password: 'Password',
}

export function getProviderName(provider) {
  return providerNames[provider] || provider?.toUpperCase()
}
