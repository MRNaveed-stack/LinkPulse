import { useMutation, useQueryClient } from '@tanstack/react-query';
import { updateProfile, UpdateProfileRequest } from '../api/profile';
import { toast } from 'react-hot-toast';

export const useUpdateProfile = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: UpdateProfileRequest) => updateProfile(data),
    onSuccess: (data) => {
      toast.success(data.message || 'Profile updated successfully');
      queryClient.invalidateQueries({ queryKey: ['currentUser'] });
      queryClient.invalidateQueries({ queryKey: ['profile'] });
      queryClient.invalidateQueries({ queryKey: ['publicProfile'] });
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Failed to update profile');
    },
  });
};
