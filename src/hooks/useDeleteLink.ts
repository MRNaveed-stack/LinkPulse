import { useMutation, useQueryClient } from '@tanstack/react-query';
import { deleteLink } from '../api/links';
import { toast } from 'react-hot-toast';

export const useDeleteLink = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deleteLink,
    onSuccess: (data) => {
      toast.success(data.message || 'Link deleted');
      queryClient.invalidateQueries({ queryKey: ['links'] });
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Delete failed');
    },
  });
};