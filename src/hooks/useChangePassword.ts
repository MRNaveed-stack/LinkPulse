import { useMutation } from '@tanstack/react-query';
import { changePassword, ChangePasswordRequest } from '../api/profile';
import { toast } from 'react-hot-toast';

export const useChangePassword = () => {
  return useMutation({
    mutationFn: (data: ChangePasswordRequest) => changePassword(data),
    onSuccess: (data) => {
      toast.success(data.message || 'Password changed successfully');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Failed to change password');
    },
  });
};
