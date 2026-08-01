import { makeIntlFormatter } from 'react-timeago/defaultFormatter'

export const intlFormatter = () => {
  const lang = localStorage.getItem('i18nextLng') || 'en'

  return makeIntlFormatter({
    locale: lang,
    localeMatcher: 'best fit',
    numberingSystem: 'latn',
    style: 'long',
    numeric: 'auto',
  })
}
