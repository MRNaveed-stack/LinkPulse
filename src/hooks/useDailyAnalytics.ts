import {useQuery} from '@tanstack/react-query';
import { getDailyAnalytics } from '../api/analytics';

export const useDailyAnalytics = (days: number = 7) => {
    return useQuery({
        queryKey: ['analytics', 'daily', days],
        queryFn: () => getDailyAnalytics(days),
    });
}

