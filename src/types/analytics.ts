export interface AnalyticsOverview {
    total_clicks: number;
    total_links : number;
    active_links: number;   
    clicks_today: number;
    top_link : {
        title: string;
        slug: string;
        clicks: number;
    } | null;
}

export interface DailyAnalyticsItem {
    date: string;
    clicks: number;
}

export interface DailyAnalyticsResponse {
    data: DailyAnalyticsItem[];
}

export interface ReferrerAnalyticsItem {
    referrer: string;
    clicks: number;
}

export interface RecentActivityItem {
    link_title: string;
    slug: string;
    clicked_at: string;
    ip_address:string;
    referrer: string;
}

export interface RecentActivityResponse {
    data: RecentActivityItem[];
}

