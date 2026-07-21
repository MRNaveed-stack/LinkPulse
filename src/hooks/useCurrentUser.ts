import { useQuery } from '@tanstack/react-query';
import { getCurrentUser } from '../api/auth';
import { useAuthStore } from '../store/authStore';
import { useEffect } from 'react';

export const useCurrentUser = () => {
  const { accessToken, setUser } = useAuthStore();

  const query = useQuery({
    queryKey: ['currentUser'],
    queryFn: getCurrentUser,
    enabled: !!accessToken,
  });

  const data = query.data;
  useEffect(() => {
    if (data) {
      setUser(data.data);
    }
  }, [data, setUser]);

  return query;
};