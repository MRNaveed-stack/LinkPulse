import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toggleStatus } from '../api/links';
import { toast } from 'react-hot-toast';

export const useToggleStatus = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, is_active }: { id: string; is_active: boolean }) =>
      toggleStatus(id, { is_active }),
    onSuccess: (res) => {
      toast.success(res.message || `Link ${res.is_active ? 'enabled' : 'disabled'}`);
      queryClient.invalidateQueries({ queryKey: ['links'] });
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Status update failed');
    },
  });
};