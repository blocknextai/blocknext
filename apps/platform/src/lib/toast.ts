import { toast as sonner } from 'sonner'

const toast = {
  info: (message: string) => sonner(message),
  success: (message: string) => sonner.success(message),
  error: (message: string) => sonner.error(message),
}

export default toast
