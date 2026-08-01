import { useEffect, useState } from 'react'
import { useLocation, useNavigate } from 'react-router'
import { useTranslation } from 'react-i18next'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { useOrganizationStore } from '@/stores/organization'
import TourComponent from '@/features/tour/components/tour-component'

const WelcomeTour = ({ sidebarOpen, setSidebarOpen, nextStep }) => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const location = useLocation()
  const [isFirstVisit, setIsFirstVisit] = useState(false)
  const [openWelcome, setOpenWelcome] = useState(
    () => !localStorage.getItem('hasVisited'),
  )
  const [shouldStartTour, setShouldStartTour] = useState(false)
  const organizationId = useOrganizationStore((s) => s.organizationId)
  const { state } = location
  useEffect(() => {
    const hasVisited = localStorage.getItem('hasVisited')
    if (!hasVisited && state?.popup !== false) {
      setIsFirstVisit(true)
    }
    if (state?.tour) {
      setShouldStartTour(true)
    }
  }, [])

  const handleSkip = () => {
    localStorage.setItem('hasVisited', 'true')
    setIsFirstVisit(false)
    setOpenWelcome(false)
    setShouldStartTour(false)
    if (location.pathname === `/organizations/${organizationId}/welcome`) {
      navigate(`/organizations/${organizationId}/create`)
    }
  }

  const handleContinue = async () => {
    setIsFirstVisit(false)
    setOpenWelcome(false)
    navigate(`/organizations/${organizationId}/welcome`, {
      state: { popup: false, tour: true },
    })
  }

  const handleTourSkip = () => {
    setShouldStartTour(false)
  }

  return (
    <>
      <TourComponent
        shouldStartTour={shouldStartTour}
        onTourSkip={handleTourSkip}
        sidebarOpen={sidebarOpen}
        setSidebarOpen={setSidebarOpen}
        onNextStep={nextStep}
      />
      {isFirstVisit && (
        <AlertDialog open={openWelcome} onOpenChange={setOpenWelcome}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>
                {t('ui.text.welcomeToBlockNext')}
              </AlertDialogTitle>
              <AlertDialogDescription>
                {t('ui.text.welcomeDescription')}
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel onClick={handleSkip}>
                {t('ui.text.skip')}
              </AlertDialogCancel>
              <AlertDialogAction onClick={handleContinue}>
                {t('ui.text.continue')}
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      )}
    </>
  )
}

WelcomeTour.displayName = 'WelcomeTour'

export default WelcomeTour
