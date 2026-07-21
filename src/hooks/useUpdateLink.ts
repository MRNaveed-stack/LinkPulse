import { useMutation, useQueryClient } from '@tanstack/react-query';
import { updateLink } from '../api/links';
import { toast } from 'react-hot-toast';

export const useUpdateLink = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => updateLink(id, data),
    onSuccess: (res) => {
      toast.success(res.message || 'Link updated');
      queryClient.invalidateQueries({ queryKey: ['links'] });
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Update failed');
    },
  });
};