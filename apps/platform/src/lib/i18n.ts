import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import LanguageDetector from 'i18next-browser-languagedetector'
import HttpBackend from 'i18next-http-backend'

// The app uses a single `ui` namespace, loaded eagerly at app boot.
const criticalNamespaces = ['ui']

i18n
  .use(HttpBackend)
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    fallbackLng: 'en',
    load: 'languageOnly',
    debug: false,
    keySeparator: false,
    interpolation: {
      escapeValue: false,
    },
    backend: {
      loadPath: '/locales/{{lng}}/{{ns}}.json',
    },
    ns: criticalNamespaces,
    defaultNS: 'ui',
    fallbackNS: 'ui',
    partialBundledLanguages: true,
  })

export default i18n
