import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createLink } from '../api/links';
import { toast } from 'react-hot-toast';

export const useCreateLink = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createLink,
    onSuccess: (data) => {
      toast.success(data.message || 'Link created');
      queryClient.invalidateQueries({ queryKey: ['links'] });
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Failed to create link');
    },
  });
};