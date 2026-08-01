import { useState, useEffect, useCallback } from 'react'
import { useLocation, useNavigate } from 'react-router'

interface UseNavigationGuardOptions {
  hasUnsavedChanges: boolean
  previewMode: boolean
}

export function useNavigationGuard({
  hasUnsavedChanges,
  previewMode,
}: UseNavigationGuardOptions) {
  const location = useLocation()
  const navigate = useNavigate()
  const [showNavigationDialog, setShowNavigationDialog] = useState(false)
  const [pendingNavigation, setPendingNavigation] = useState<string | null>(
    null,
  )

  const handleNavigationConfirm = useCallback(() => {
    setShowNavigationDialog(false)

    if (pendingNavigation === 'back') {
      window.history.back()
    } else if (pendingNavigation) {
      navigate(pendingNavigation)
    }
    setPendingNavigation(null)
  }, [pendingNavigation, navigate])

  const handleNavigationCancel = useCallback(() => {
    setShowNavigationDialog(false)
    setPendingNavigation(null)
  }, [])

  // Link click and popstate interception
  useEffect(() => {
    if (previewMode) return
    if (!hasUnsavedChanges) return

    const handleClick = (e: MouseEvent) => {
      const link = (e.target as HTMLElement).closest('a[href]')
      if (link && hasUnsavedChanges) {
        e.preventDefault()
        e.stopPropagation()

        const href = link.getAttribute('href')
        if (href && !href.startsWith('http') && href !== location.pathname) {
          setPendingNavigation(href)
          setShowNavigationDialog(true)
        }
      }
    }

    const handlePopState = (e: PopStateEvent) => {
      if (hasUnsavedChanges) {
        e.preventDefault()
        window.history.pushState(null, '', location.pathname)
        setShowNavigationDialog(true)
        setPendingNavigation('back')
      }
    }

    document.addEventListener('click', handleClick, true)
    window.addEventListener('popstate', handlePopState)

    window.history.pushState(null, '', location.pathname)

    return () => {
      document.removeEventListener('click', handleClick, true)
      window.removeEventListener('popstate', handlePopState)
    }
  }, [hasUnsavedChanges, location.pathname, previewMode])

  // beforeunload
  useEffect(() => {
    if (previewMode) return
    const handleBeforeUnload = (e: BeforeUnloadEvent) => {
      if (hasUnsavedChanges) {
        e.preventDefault()
        e.returnValue = ''
        return ''
      }
    }

    window.addEventListener('beforeunload', handleBeforeUnload)

    return () => {
      window.removeEventListener('beforeunload', handleBeforeUnload)
    }
  }, [hasUnsavedChanges, previewMode])

  return {
    showNavigationDialog,
    setShowNavigationDialog,
    pendingNavigation,
    handleNavigationConfirm,
    handleNavigationCancel,
  }
}
