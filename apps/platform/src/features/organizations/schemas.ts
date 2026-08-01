import { z } from 'zod'
import {
  emailSchema,
  githubUsernameSchema,
  walletAddressSchema,
} from '@/lib/validation'

export const inviteMemberSchema = z.discriminatedUnion('provider', [
  z.object({
    provider: z.literal('wallet'),
    identifier: walletAddressSchema,
  }),
  z.object({
    provider: z.literal('metamask'),
    identifier: walletAddressSchema,
  }),
  z.object({
    provider: z.literal('email'),
    identifier: emailSchema,
  }),
  z.object({
    provider: z.literal('google'),
    identifier: emailSchema,
  }),
  z.object({
    provider: z.literal('github'),
    identifier: githubUsernameSchema,
  }),
])

export type InviteMemberInput = z.infer<typeof inviteMemberSchema>

export const organizationBasicInfoSchema = z.object({
  title: z.string().trim().min(1, 'ui.text.required').max(100),
  description: z.string().trim().max(500).optional(),
})

export type OrganizationBasicInfo = z.infer<typeof organizationBasicInfoSchema>
