import { z } from 'zod'

export const walletAddressSchema = z
  .string()
  .trim()
  .regex(/^0x[a-fA-F0-9]{40}$/, 'ui.text.invalidWallet')

export const emailSchema = z.string().trim().email('ui.text.invalidEmail')

export const githubUsernameSchema = z
  .string()
  .trim()
  .regex(/^[a-zA-Z0-9-]+$/, 'ui.text.invalidGithub')
  .min(1, 'ui.text.invalidGithub')

export const uuidSchema = z.string().uuid()

export const nonEmptyStringSchema = (message = 'ui.text.required') =>
  z.string().trim().min(1, message)

export type ValidationResult<T> =
  { success: true; data: T } | { success: false; message: string }

export function firstErrorMessage(error: z.ZodError): string {
  const issue = error.issues[0]
  return issue?.message ?? 'ui.text.validationFailed'
}
