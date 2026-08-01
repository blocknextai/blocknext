import { heroui } from '@heroui/theme'

export default {
  darkMode: 'class',
  content: [
    './index.html',
    './src/**/*.{js,jsx,ts,tsx}',
    './node_modules/@heroui/theme/dist/components/(toast|spinner).js',
  ],
  safelist: ['dark', 'ProseMirror'],
  plugins: [heroui()],
}
