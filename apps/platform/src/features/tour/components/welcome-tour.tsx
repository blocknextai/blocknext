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
import {
  hasSeenOnboardingTour,
  markOnboardingTourSeen,
} from '@/features/tour/storage'
import { ONBOARDING_DEMO } from '@/features/tour/demo-flow'

const WelcomeTour = ({
  sidebarOpen,
  setSidebarOpen,
  addNodeById,
  selectGroupForNode,
  connectNodes,
  focusCanvas,
  openNodeSettings,
  getFlowState,
  setChatOpen,
  nextStep,
}) => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const location = useLocation()
  const [isFirstVisit, setIsFirstVisit] = useState(false)
  const [openWelcome, setOpenWelcome] = useState(() => !hasSeenOnboardingTour())
  const [shouldStartTour, setShouldStartTour] = useState(false)
  const organizationId = useOrganizationStore((s) => s.organizationId)
  const { state } = location
  const createPath = `/organizations/${organizationId}/create`
  const demoPath = `${createPath}?demo=${ONBOARDING_DEMO}`

  useEffect(() => {
    if (!hasSeenOnboardingTour() && state?.popup !== false) {
      setIsFirstVisit(true)
    }
  }, [])

  useEffect(() => {
    if (state?.tour) {
      setShouldStartTour(true)
    }
  }, [state?.tour])

  const handleSkip = () => {
    markOnboardingTourSeen()
    setIsFirstVisit(false)
    setOpenWelcome(false)
    setShouldStartTour(false)
  }

  const handleContinue = () => {
    markOnboardingTourSeen()
    setIsFirstVisit(false)
    setOpenWelcome(false)
    navigate(demoPath, {
      replace: true,
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
        addNodeById={addNodeById}
        selectGroupForNode={selectGroupForNode}
        connectNodes={connectNodes}
        focusCanvas={focusCanvas}
        openNodeSettings={openNodeSettings}
        getFlowState={getFlowState}
        setChatOpen={setChatOpen}
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
