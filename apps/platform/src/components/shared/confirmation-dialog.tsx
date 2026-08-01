import type { ReactNode } from 'react'
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
import { Button } from '@/components/ui/button'

type ButtonVariant =
  'default' | 'outline' | 'secondary' | 'ghost' | 'destructive' | 'link'

interface ConfirmationDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  title?: string
  description?: ReactNode
  confirmText?: string
  cancelText?: string
  onConfirm?: () => void | Promise<void>
  variant?: ButtonVariant
  confirmVariant?: ButtonVariant
}

const ConfirmationDialog = ({
  open,
  onOpenChange,
  title = 'Are you absolutely sure?',
  description,
  confirmText = 'Confirm',
  cancelText = 'Cancel',
  onConfirm,
  variant = 'default',
  confirmVariant,
}: ConfirmationDialogProps) => {
  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{title}</AlertDialogTitle>
          <AlertDialogDescription>{description}</AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>{cancelText}</AlertDialogCancel>
          <AlertDialogAction asChild>
            <Button variant={confirmVariant || variant} onClick={onConfirm}>
              {confirmText}
            </Button>
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}

export default ConfirmationDialog
