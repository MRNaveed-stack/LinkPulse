import { useMutation } from '@tanstack/react-query';
import { deleteAccount } from '../api/profile';
import { useAuthStore } from '../store/authStore';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-hot-toast';

export const useDeleteAccount = () => {
  const logout = useAuthStore((state) => state.logout);
  const navigate = useNavigate();

  return useMutation({
    mutationFn: deleteAccount,
    onSuccess: (data) => {
      toast.success(data.message || 'Account successfully deleted');
      logout();
      navigate('/login');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Failed to delete account');
    },
  });
};
