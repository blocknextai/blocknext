const cookie = {
  get: (key: string): string | null => {
    const name = key + '='
    const decodedCookie = decodeURIComponent(document.cookie)
    const cookies = decodedCookie.split(';')

    for (let i = 0; i < cookies.length; i++) {
      let cookie = cookies[i]
      while (cookie.charAt(0) === ' ') {
        cookie = cookie.substring(1)
      }
      if (cookie.indexOf(name) === 0) {
        return cookie.substring(name.length, cookie.length)
      }
    }
    return null
  },

  set: (key: string, value: string, days: number = 7): void => {
    const date = new Date()
    date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000)
    const expires = 'expires=' + date.toUTCString()
    document.cookie = `${key}=${value};${expires};path=/;SameSite=Lax`
  },

  remove: (key: string): void => {
    document.cookie = `${key}=;expires=Thu, 01 Jan 1970 00:00:00 UTC;path=/;`
  },
}

export default cookie
