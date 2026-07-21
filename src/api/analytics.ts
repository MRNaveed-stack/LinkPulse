import apiClient from './client';
import type {
  AnalyticsOverview,
  DailyAnalyticsResponse,
  ReferrerAnalyticsItem,
  RecentActivityResponse,
} from '../types/analytics';

export const getOverview = () =>
  apiClient.get<AnalyticsOverview>('/analytics/overview').then((res) => res.data);

export const getDailyAnalytics = (days: number = 7) =>
  apiClient
    .get<DailyAnalyticsResponse>('/analytics/daily', { params: { days } })
    .then((res) => res.data);

export const getReferrerAnalytics = () =>
  apiClient.get<ReferrerAnalyticsItem[]>('/analytics/referrers').then((res) => res.data);

export const getRecentActivity = (limit: number = 20) =>
  apiClient
    .get<RecentActivityResponse>('/analytics/recent', { params: { limit } })
    .then((res) => res.data);