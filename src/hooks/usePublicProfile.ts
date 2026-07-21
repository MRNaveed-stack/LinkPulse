import { useQuery } from '@tanstack/react-query';
import { getPublicProfile } from '../api/public';

export const usePublicProfile = (username: string) => {
  return useQuery({
    queryKey: ['publicProfile', username],
    queryFn: () => getPublicProfile(username),
    enabled: !!username,
  });
};