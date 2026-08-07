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
  color: 'var(--muted-foreground)',
  icon: Box,
  labelKey: '',
}

export const getCategoryPrefs = (key?: string) => {
  if (!key) {
    return DEFAULT_PREFS
  }
  return (
    (flowCategoryPreferences as Record<string, typeof DEFAULT_PREFS>)[key] ??
    DEFAULT_PREFS
  )
}

export const flowCategoryPreferences = {
  mailing: {
    color: 'var(--category-mailing)',
    icon: Mail,
    labelKey: 'ui.text.category.mail',
  },
  genai: {
    color: 'var(--category-genai)',
    icon: Sparkles,
    labelKey: 'ui.text.category.ai',
  },
  ai: {
    color: 'var(--category-ai)',
    icon: Sparkles,
    labelKey: 'ui.text.category.ai',
  },
  image: {
    color: 'var(--category-image)',
    icon: Image,
    labelKey: 'ui.text.category.image',
  },
  googleworkspace: {
    color: 'var(--category-googleworkspace)',
    icon: Cloud,
    labelKey: 'ui.text.category.googleWorkspace',
  },
  publishing: {
    color: 'var(--category-publishing)',
    icon: MonitorUp,
    labelKey: 'ui.text.category.publishing',
  },
  system: {
    color: 'var(--category-system)',
    icon: Wrench,
    labelKey: 'ui.text.category.system',
  },
  audio: {
    color: 'var(--category-audio)',
    icon: AudioLines,
    labelKey: 'ui.text.category.audio',
  },
  video: {
    color: 'var(--category-video)',
    icon: Video,
    labelKey: 'ui.text.category.video',
  },
  blockchain: {
    color: 'var(--category-blockchain)',
    icon: Blocks,
    labelKey: 'ui.text.category.blockchain',
  },
  database: {
    color: 'var(--category-database)',
    icon: Database,
    labelKey: 'ui.text.category.database',
  },
  marketing: {
    color: 'var(--category-marketing)',
    icon: TrendingUp,
    labelKey: 'ui.text.category.marketing',
  },
  ecommerce: {
    color: 'var(--category-ecommerce)',
    icon: ShoppingCart,
    labelKey: 'ui.text.category.ecommerce',
  },
  crm: {
    color: 'var(--category-crm)',
    icon: Users,
    labelKey: 'ui.text.category.crm',
  },
  productivity: {
    color: 'var(--category-productivity)',
    icon: Calendar,
    labelKey: 'ui.text.category.productivity',
  },
  analytics: {
    color: 'var(--category-analytics)',
    icon: BarChart3,
    labelKey: 'ui.text.category.analytics',
  },
  utility: {
    color: 'var(--category-utility)',
    icon: Lightbulb,
    labelKey: 'ui.text.category.utility',
  },
}
