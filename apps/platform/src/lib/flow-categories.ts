import {
  Mail,
  Cloud,
  Sparkles,
  Image,
  MonitorUp,
  Wrench,
  AudioLines,
  Video,
  Blocks,
  Database,
  TrendingUp,
  ShoppingCart,
  Users,
  Calendar,
  BarChart3,
  Lightbulb,
  Box,
} from 'lucide-react'

const DEFAULT_PREFS = {
  color: 'rgba(120, 120, 120, 1)',
  icon: Box,
  labelKey: '',
}

export const getCategoryPrefs = (key?: string) => {
  if (!key) return DEFAULT_PREFS
  return (
    (flowCategoryPreferences as Record<string, typeof DEFAULT_PREFS>)[key] ??
    DEFAULT_PREFS
  )
}

export const flowCategoryPreferences = {
  mailing: {
    color: 'rgba(33, 221, 102, 1)',
    icon: Mail,
    labelKey: 'ui.text.category.mail',
  },
  genai: {
    color: 'rgba(255, 187, 0, 1)',
    icon: Sparkles,
    labelKey: 'ui.text.category.ai',
  },
  ai: {
    color: 'rgba(255, 187, 0, 1)',
    icon: Sparkles,
    labelKey: 'ui.text.category.ai',
  },
  image: {
    color: 'rgba(69, 221, 222, 1)',
    icon: Image,
    labelKey: 'ui.text.category.image',
  },
  googleworkspace: {
    color: 'rgba(11, 187, 255, 1)',
    icon: Cloud,
    labelKey: 'ui.text.category.googleWorkspace',
  },
  mediapublishing: {
    color: 'rgba(182, 86, 255, 1)',
    icon: MonitorUp,
    labelKey: 'ui.text.category.publishing',
  },
  system: {
    color: 'rgba(118, 162, 195, 1)',
    icon: Wrench,
    labelKey: 'ui.text.category.system',
  },
  audio: {
    color: 'rgba(7, 119, 255, 1)',
    icon: AudioLines,
    labelKey: 'ui.text.category.audio',
  },
  video: {
    color: 'rgba(255, 86, 119, 1)',
    icon: Video,
    labelKey: 'ui.text.category.video',
  },
  blockchain: {
    color: 'rgba(252, 165, 3, 1)',
    icon: Blocks,
    labelKey: 'ui.text.category.blockchain',
  },
  database: {
    color: 'rgba(19, 184, 166, 1)',
    icon: Database,
    labelKey: 'ui.text.category.database',
  },
  marketing: {
    color: 'rgba(239, 68, 68, 1)',
    icon: TrendingUp,
    labelKey: 'ui.text.category.marketing',
  },
  ecommerce: {
    color: 'rgba(219, 39, 119, 1)',
    icon: ShoppingCart,
    labelKey: 'ui.text.category.ecommerce',
  },
  crm: {
    color: 'rgba(79, 70, 229, 1)',
    icon: Users,
    labelKey: 'ui.text.category.crm',
  },
  productivity: {
    color: 'rgba(132, 204, 22, 1)',
    icon: Calendar,
    labelKey: 'ui.text.category.productivity',
  },
  analytics: {
    color: 'rgba(217, 119, 6, 1)',
    icon: BarChart3,
    labelKey: 'ui.text.category.analytics',
  },
  utility: {
    color: 'rgba(100, 116, 139, 1)',
    icon: Lightbulb,
    labelKey: 'ui.text.category.utility',
  },
}
