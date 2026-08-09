import {
  ChevronLeft,
  PanelLeft,
  PanelLeftClose,
  Search,
  SearchX,
  Sparkles,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { GhostInput } from '@/components/shared/ghost-input'
import { getCategoryPrefs } from '@/lib/flow-categories'
import { useIconResolver } from '@/features/flow-editor/icons'
import { usePlatformFeatures } from '@/features/platform'
import { useTranslation } from 'react-i18next'

const FlowSidebar = ({
  previewMode,
  sidebarOpen,
  setSidebarOpen,
  selectedCategory,
  setSelectedCategory,
  providerList,
  dragging,
  startDrag,
  handleSearch,
  setChatOpen,
}) => {
  const { t } = useTranslation()
  const { workflowsGenerationEnabled } = usePlatformFeatures()
  const resolveIcon = useIconResolver()

  const renderNodeCard = (item, key) => {
    const name = item.label || item.name
    const prefs = getCategoryPrefs(item.type || item.category)
    const Icon = resolveIcon(item.icon) ?? prefs.icon
    const isDisabled = item.isComingSoon
    const cursorClass = isDisabled
      ? 'cursor-not-allowed'
      : dragging
        ? 'cursor-grabbing'
        : 'cursor-grab'
    const nodeLabel = t(name)
    const ariaLabel = isDisabled
      ? `${nodeLabel} (${t('ui.text.comingSoon')})`
      : nodeLabel

    return (
      <div
        key={key}
        data-tour="node-item"
        style={{ borderColor: prefs.color }}
        className={`bg-card hover:bg-accent flex w-full items-center gap-3 rounded-md border px-3 py-2.5 transition-colors ${cursorClass} ${isDisabled ? 'opacity-50' : ''}`}
        onMouseDown={(e) => !isDisabled && startDrag(e, item.id)}
        role="button"
        tabIndex={isDisabled ? -1 : 0}
        aria-label={ariaLabel}
        aria-disabled={isDisabled}
      >
        <Icon className="size-5 shrink-0" aria-hidden="true" />
        <div className="flex min-w-0 flex-1 flex-col">
          <div className="flex items-center gap-2">
            <span className="text-sm">{t(name)}</span>
            {item.isComingSoon && (
              <Badge
                variant="secondary"
                className="shrink-0 text-xs bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200"
              >
                {t('ui.text.comingSoon')}
              </Badge>
            )}
          </div>
          <span className="text-xs text-muted-foreground">
            {t(item.description)}
          </span>
        </div>
      </div>
    )
  }

  const renderNodeList = () => {
    if (selectedCategory) {
      const items = selectedCategory.nodes ?? []

      if (items.length === 0) {
        return (
          <li
            className="flex flex-1 flex-col items-center justify-center gap-2 p-8 text-center"
            data-tour="node-list"
          >
            <SearchX
              className="size-8 text-muted-foreground/70"
              aria-hidden="true"
            />
            <div className="text-sm font-medium">
              {t('ui.text.noNodesFound')}
            </div>
            <div className="text-xs text-muted-foreground">
              {t('ui.text.noNodesFoundDescription')}
            </div>
          </li>
        )
      }

      return (
        <li className="flex flex-1 flex-col gap-2 p-4" data-tour="node-list">
          {items.map((item: any, j: number) => renderNodeCard(item, j))}
        </li>
      )
    }

    const providers = providerList.map((provider: any) => {
      const prefs = getCategoryPrefs(provider.nodes[0]?.type ?? '')
      const Icon = resolveIcon(provider.icon) ?? prefs.icon

      return (
        <div
          key={provider.key}
          onClick={() =>
            setSelectedCategory({
              label: provider.label,
              nodes: provider.nodes,
            })
          }
          style={{ borderColor: prefs.color }}
          className="bg-card hover:bg-accent flex w-full cursor-pointer items-center gap-3 rounded-md border px-3 py-2.5 transition-colors"
          data-tour={provider.key === 'system' ? 'system-category' : undefined}
          role="button"
          tabIndex={0}
          aria-label={provider.label}
        >
          <Icon className="size-5 shrink-0" aria-hidden="true" />
          <span className="flex-1 truncate text-sm">{provider.label}</span>
          <span className="shrink-0 text-xs text-muted-foreground">
            {t('ui.text.toolCount', { count: provider.nodes.length })}
          </span>
        </div>
      )
    })

    return (
      <li className="flex flex-1 flex-col gap-2 p-4" data-tour="category-grid">
        {providers}
      </li>
    )
  }

  if (previewMode) {
    return null
  }

  return (
    <div className="drawer w-full h-full">
      <input
        id="flowSideMenu"
        type="checkbox"
        className="drawer-toggle"
        checked={sidebarOpen}
        onChange={() => setSidebarOpen(!sidebarOpen)}
      />
      <div className="drawer-content">
        <div className="z-49 absolute top-5 left-5 flex items-center justify-center gap-2">
          <label
            htmlFor="flowSideMenu"
            aria-label={t('ui.text.openSidebar')}
            className="drawer-button cursor-pointer"
            data-tour="sidebar-toggle"
          >
            <div className="z-30 p-3 bg-primary text-primary-foreground rounded-lg">
              <PanelLeft className="size-6" aria-hidden="true" />
            </div>
          </label>

          {!previewMode && workflowsGenerationEnabled && (
            <Button
              onClick={() => setChatOpen((prev) => !prev)}
              data-tour="ai-flow-builder"
              className="rounded-xl shadow-xs !py-6 px-4 cursor-pointer flex items-center gap-2"
            >
              <Sparkles className="size-5" />{' '}
              {t('ui.text.generation.aiFlowBuilder')}
            </Button>
          )}
        </div>
      </div>
      <div className="drawer-side w-sm! p-0! rounded-xl bg-background! relative z-50 top-4 left-4 h-[calc(100%-2rem)] shadow-lg">
        <ul className="min-h-full w-full p-0 rounded-s flex flex-col">
          <li
            className="flex cursor-default! w-full flex-row!
                        justify-between items-center sticky
                        right-0 z-55 top-0 border-b bg-background dark:border-zinc-800 border-zinc-200 px-4 py-2"
          >
            <div className="px-2 py-2 text-sm flex justify-start items-center gap-6 flex-1 cursor-default!">
              {selectedCategory ? (
                <div
                  className="flex items-center gap-2 cursor-pointer"
                  onClick={() => setSelectedCategory(null)}
                  role="button"
                  tabIndex={0}
                  aria-label={t('ui.text.back')}
                >
                  <div className="cursor-pointer transition-colors rounded-sm">
                    {' '}
                    <ChevronLeft className="size-4" aria-hidden="true" />
                  </div>
                  <div className="font-semibold">{selectedCategory.label}</div>
                </div>
              ) : (
                <></>
              )}
              <div
                className="flex items-center gap-2 flex-1"
                data-tour="search-input"
              >
                <Search className="text-foreground size-4" aria-hidden="true" />
                <GhostInput
                  placeholder={t('ui.text.searchNodes')}
                  aria-label={t('ui.text.searchNodes')}
                  className="w-full"
                  onChange={handleSearch}
                />
              </div>
            </div>

            <label
              htmlFor="flowSideMenu"
              aria-label={t('ui.text.closeSidebar')}
              className="size-8 flex items-center transition-colors rounded-sm! justify-center drawer-button cursor-pointer! hover:bg-accent! hover:text-foreground!"
            >
              <PanelLeftClose className="size-5" aria-hidden="true" />
            </label>
          </li>
          {renderNodeList()}
        </ul>
      </div>
    </div>
  )
}
FlowSidebar.displayName = 'FlowSidebar'

export { FlowSidebar }
