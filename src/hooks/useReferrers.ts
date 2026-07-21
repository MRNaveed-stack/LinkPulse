import {useQuery} from '@tanstack/react-query';
import {getReferrerAnalytics} from '../api/analytics';

export const useReferrers = () => {
    return useQuery({
        queryKey: ['analytics', 'referrers'],
        queryFn: () => getReferrerAnalytics(),
    });
}

