import toast from 'react-hot-toast';
import { ExclamationTriangleIcon, CheckCircleIcon, InformationCircleIcon } from '@heroicons/react/24/outline';

interface ToastOptions {
  duration?: number;
  position?: 'top-left' | 'top-center' | 'top-right' | 'bottom-left' | 'bottom-center' | 'bottom-right';
}

/**
 * Show a success toast notification
 */
export const showSuccess = (message: string, options?: ToastOptions) => {
  return toast.success(message, {
    ...options,
    icon: CheckCircleIcon
  });
};

/**
 * Show an error toast notification
 */
export const showError = (message: string, options?: ToastOptions) => {
  return toast.error(message, {
    ...options,
    duration: options?.duration || 5000, // Error messages stay a bit longer by default
    icon: ExclamationTriangleIcon
  });
};

/**
 * Show an info toast notification
 */
export const showInfo = (message: string, options?: ToastOptions) => {
  return toast(message, {
    ...options,
    icon: InformationCircleIcon,
    style: {
      background: 'var(--toast-info-bg, #EFF6FF)',
      color: 'var(--toast-info-color, #1E40AF)'
    },
  });
};

/**
 * Show a loading toast notification that can be updated later
 */
export const showLoading = (message: string) => {
  return toast.loading(message, {
    style: {
      background: 'var(--toast-loading-bg, #F3F4F6)',
      color: 'var(--toast-loading-color, #1F2937)'
    },
  });
};

/**
 * Dismiss a specific toast by its ID
 */
export const dismissToast = (toastId: string) => {
  toast.dismiss(toastId);
};

/**
 * Dismiss all currently visible toasts
 */
export const dismissAllToasts = () => {
  toast.dismiss();
};

/**
 * Update an existing toast (useful for loading -> success/error transitions)
 */
export const updateToast = (
  toastId: string, 
  message: string, 
  type: 'success' | 'error' | 'loading' | 'info'
) => {
  if (type === 'success') {
    toast.success(message, { id: toastId });
  } else if (type === 'error') {
    toast.error(message, { id: toastId });
  } else if (type === 'loading') {
    toast.loading(message, { id: toastId });
  } else {
    toast(message, { id: toastId });
  }
};

export default {
  success: showSuccess,
  error: showError,
  info: showInfo,
  loading: showLoading,
  dismiss: dismissToast,
  dismissAll: dismissAllToasts,
  update: updateToast
};