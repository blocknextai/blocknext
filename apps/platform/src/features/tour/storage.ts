const ONBOARDING_TOUR_STORAGE_KEY = 'onboardingTourSeen'

export const hasSeenOnboardingTour = () =>
  localStorage.getItem(ONBOARDING_TOUR_STORAGE_KEY) !== null

export const markOnboardingTourSeen = () =>
  localStorage.setItem(ONBOARDING_TOUR_STORAGE_KEY, 'true')
