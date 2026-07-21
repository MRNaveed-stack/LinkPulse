import { useQuery } from '@tanstack/react-query';
import apiClient from '../api/client';

interface OverviewData {
  totalClicks: number;
  totalLinks: number;
  activeLinks: number;
  clicksToday: number;
  topPerformingLink?: {
    id: string;
    title: string;
    shortUrl: string;
    clicks: number;
  };
}

export const useAnalyticsOverview = () => {
  return useQuery<OverviewData>({
    queryKey: ['analyticsOverview'],
    queryFn: () =>
      apiClient.get('/analytics/overview').then((res) => {
        const data = res.data;
        return {
          totalClicks: data.total_clicks,
          totalLinks: data.total_links,
          activeLinks: data.active_links,
          clicksToday: data.clicks_today,
          topPerformingLink: data.top_link
            ? {
                title: data.top_link.title,
                clicks: data.top_link.clicks,
                id: data.top_link.slug,
                shortUrl: `/u/${data.top_link.slug}`,
              }
            : undefined,
        };
      }),
  });
};