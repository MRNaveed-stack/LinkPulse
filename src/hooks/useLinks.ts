import { useQuery } from '@tanstack/react-query';
import { getLinks } from '../api/links';

export const useLinks = () => {
  return useQuery({
    queryKey: ['links'],
    queryFn: getLinks,
  });
};

