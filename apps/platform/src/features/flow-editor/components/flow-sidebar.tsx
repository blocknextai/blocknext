import {
  PanelLeft,
  PanelLeftClose,
  Search,
  ChevronLeft,
  Sparkles,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { GhostInput } from '@/components/shared/ghost-input'
import { getCategoryPrefs } from '@/lib/flow-categories'
import { useTranslation } from 'react-i18next'

const FlowSidebar = ({
  previewMode,
  sidebarOpen,
  setSidebarOpen,
  selectedCategory,
  setSelectedCategory,
  nodeList,
  dragging,
  startDrag,
  annotationDrag,
  handleSearch,
  setChatOpen,
}) => {
  const { t } = useTranslation()

  const renderNodeCard = (item, key) => {
    const name = item.label || item.name
    const prefs = getCategoryPrefs(item.type || item.category)
    const Icon = prefs.icon
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
        style={{ color: prefs.color }}
        className={`border rounded-md p-4 flex flex-1 w-full flex-col transition-colors dark:bg-current/20 bg-current/10 dark:hover:bg-current/10 hover:bg-current/20 items-start justify-start gap-2 ${cursorClass} ${isDisabled ? 'opacity-50' : ''}`}
        onMouseDown={(e) => !isDisabled && startDrag(e, item.id)}
        role="button"
        tabIndex={isDisabled ? -1 : 0}
        aria-label={ariaLabel}
        aria-disabled={isDisabled}
      >
        <div className="flex items-start justify-start gap-3 w-full">
          <div style={{ color: prefs.color }} aria-hidden="true">
            <Icon className="size-6" />
          </div>
          <div className="flex flex-col gap-1 flex-1 min-w-0">
            {item.isComingSoon && (
              <div className="flex gap-2">
                <Badge
                  variant="secondary"
                  className="text-xs bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200"
                >
                  {t('ui.text.comingSoon')}
                </Badge>
              </div>
            )}
            <div>{t(name)}</div>
          </div>
        </div>
        <div className="text-foreground">{t(item.description)}</div>
      </div>
    )
  }

  const renderAnnotationCard = () => {
    const annotationPrefs = getCategoryPrefs('system')
    const AnnotationIcon = annotationPrefs?.icon

    return (
      <div
        key="annotation-node"
        style={{ color: annotationPrefs?.color }}
        className={`border rounded-md p-4 flex flex-1 w-full flex-col transition-colors dark:bg-current/20 bg-current/10  dark:hover:bg-current/10  hover:bg-current/20 items-start justify-start gap-2 ${dragging ? 'cursor-grabbing' : 'cursor-grab'}`}
        onMouseDown={annotationDrag}
        role="button"
        tabIndex={0}
        aria-label={t('ui.text.annotationNode')}
      >
        <div className="flex items-center justify-start gap-3 w-full">
          <div style={{ color: annotationPrefs.color }} aria-hidden="true">
            <AnnotationIcon className="size-6" />
          </div>
          <div>{t('ui.text.annotationNode')}</div>
        </div>
        <div className="text-foreground">
          {t('ui.text.addAnnotationToFlow')}
        </div>
      </div>
    )
  }

  const renderNodeList = () => {
    if (selectedCategory) {
      const isFlat = Array.isArray(selectedCategory)
      const flatItems = isFlat
        ? selectedCategory
        : Object.values(selectedCategory).flat()
      const hasSystem = flatItems.some((cat: any) => cat.type === 'system')

      let body: React.ReactNode

      if (isFlat) {
        body = (
          <>
            {selectedCategory.map((item, j) => renderNodeCard(item, j))}
            {hasSystem && renderAnnotationCard()}
          </>
        )
      } else {
        const subEntries = Object.entries(
          selectedCategory as Record<string, any[]>,
        )
        body = (
          <>
            {subEntries.map(([sub, items]) => (
              <div key={sub || '__default'} className="flex flex-col gap-3">
                {sub && (
                  <div className="text-xs font-semibold uppercase tracking-wider text-muted-foreground px-1">
                    {t(sub)}
                  </div>
                )}
                {items.map((item, j) => renderNodeCard(item, `${sub}-${j}`))}
              </div>
            ))}
            {hasSystem && renderAnnotationCard()}
          </>
        )
      }

      return (
        <li className={`flex flex-col gap-4 p-4`} data-tour="node-list">
          {body}
        </li>
      )
    }

    const categories = Object.keys(nodeList).map((category, i) => {
      const prefs = getCategoryPrefs(category)
      const Icon = prefs.icon
      const categoryLabel = prefs.labelKey ? t(prefs.labelKey) : category

      return (
        <div
          key={i}
          onClick={() => setSelectedCategory(nodeList[category])}
          className={`flex flex-col items-center justify-center rounded-lg transition-colors w-full min-h-33 cursor-pointer bg-current/10 dark:hover:bg-current/50 hover:bg-current/20 gap-3`}
          style={{ color: prefs.color }}
          data-tour={category === 'system' ? 'system-category' : undefined}
          role="button"
          tabIndex={0}
          aria-label={categoryLabel}
        >
          <Icon className="size-8" aria-hidden="true" />
          <div className="capitalize">{categoryLabel}</div>
        </div>
      )
    })

    return (
      <li
        className={`grid grid-cols-2 gap-4 p-4 flex-1`}
        data-tour="category-grid"
      >
        {categories}
      </li>
    )
  }

  if (previewMode) return null

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

          {!previewMode && (
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
        <ul className="min-h-full w-full p-0 rounded-s flex flex-col justify-between">
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
                  <div className="capitalize font-semibold">
                    {(() => {
                      const flat = Array.isArray(selectedCategory)
                        ? selectedCategory
                        : Object.values(
                            selectedCategory as Record<string, any[]>,
                          ).flat()
                      return flat[0]?.type ?? 'Results'
                    })()}
                  </div>
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
