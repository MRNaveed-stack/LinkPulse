import { useState } from 'react';
import { MousePointerClick, LinkIcon, BarChart3, Calendar, TrendingUp } from 'lucide-react';
import { useAnalyticsOverview } from '../../hooks/useAnalyticsOverview';
import OverviewCard from '../../components/analytics/OverviewCard';
import ClickChart from '../../components/analytics/ClickChart';
import ReferrerChart from '../../components/analytics/ReferrerChart';
import RecentActivityTable from '../../components/analytics/RecentActivityTable';
import AnalyticsFilter from '../../components/analytics/AnalyticsFilter';
import Card from '../../components/common/Card';
import Button from '../../components/common/Button';

export default function AnalyticsPage() {
  const [days, setDays] = useState(7);
  const { data: overview, isLoading: overviewLoading, isError: overviewError } = useAnalyticsOverview();

  return (
    <div>
      <div className="sm:flex sm:items-center sm:justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-900">Analytics</h2>
        <div className="mt-3 sm:mt-0">
          <AnalyticsFilter selectedDays={days} onDaysChange={setDays} />
        </div>
      </div>

      {/* Overview Cards */}
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 mb-8">
        <OverviewCard
          title="Total Clicks"
          value={overview?.totalClicks ?? 0}
          loading={overviewLoading}
          icon={<MousePointerClick className="h-8 w-8 text-indigo-500" />}
        />
        <OverviewCard
          title="Total Links"
          value={overview?.totalLinks ?? 0}
          loading={overviewLoading}
          icon={<LinkIcon className="h-8 w-8 text-green-500" />}
        />
        <OverviewCard
          title="Active Links"
          value={overview?.activeLinks ?? 0}
          loading={overviewLoading}
          icon={<BarChart3 className="h-8 w-8 text-yellow-500" />}
        />
        <OverviewCard
          title="Clicks Today"
          value={overview?.clicksToday ?? 0}
          loading={overviewLoading}
          icon={<Calendar className="h-8 w-8 text-purple-500" />}
        />
        <OverviewCard
          title="Top Link"
          value={overview?.topPerformingLink ? overview.topPerformingLink.title : 'N/A'}
          loading={overviewLoading}
          icon={<TrendingUp className="h-8 w-8 text-pink-500" />}
        />
      </div>

      {overviewError && (
        <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700 flex items-center justify-between">
          <span>Unable to load overview data.</span>
          <Button variant="secondary"  onClick={() => window.location.reload()}>
            Retry
          </Button>
        </div>
      )}

      {/* Charts Row */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
        <Card className="lg:col-span-2 p-5">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Click Trend</h3>
          <ClickChart days={days} />
        </Card>
        <Card className="p-5">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Referrer Sources</h3>
          <ReferrerChart />
        </Card>
      </div>

      {/* Recent Activity */}
      <Card className="p-5">
        <h3 className="text-lg font-medium text-gray-900 mb-4">Recent Activity</h3>
        <RecentActivityTable />
      </Card>
    </div>
  );
}