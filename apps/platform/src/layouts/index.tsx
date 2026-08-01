import { lazy } from 'react'

import RootLayout from '@/layouts/root-layout'
import AuthLayout from '@/layouts/auth-layout'
import ProtectedLayout from '@/layouts/protected-layout'

import RequireOrganization from '@/layouts/require-organization'

const OrganizationLayout = lazy(() => import('@/layouts/organization-layout'))
const FlowLayout = lazy(() => import('@/layouts/flow-layout'))

export {
  RootLayout,
  AuthLayout,
  ProtectedLayout,
  RequireOrganization,
  OrganizationLayout,
  FlowLayout,
}
