import { Outlet } from 'react-router'
import { TooltipProvider } from '@/components/ui/tooltip'

const RootLayout = () => {
  return (
    <>
      <TooltipProvider>
        <Outlet />
      </TooltipProvider>
    </>
  )
}

export default RootLayout
