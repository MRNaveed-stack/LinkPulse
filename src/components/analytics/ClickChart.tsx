import {
  ResponsiveContainer,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
} from 'recharts';
import { useDailyAnalytics } from '../../hooks/useDailyAnalytics';
import Skeleton from '../common/Skeleton';

interface ClickChartProps {
  days: number;
}

const ClickChart = ({ days }: ClickChartProps) => {
  const { data, isLoading, isError } = useDailyAnalytics(days);

  if (isError) {
    return (
      <div className="p-6 text-center text-red-600 bg-red-50 rounded-lg">
        Unable to load chart data.
      </div>
    );
  }

  if (isLoading) {
    return <Skeleton className="h-64 w-full rounded-lg" />;
  }

  const chartData = data?.data || [];
  const hasData = chartData.some((d) => d.clicks > 0);

  if (!hasData) {
    return (
      <div className="p-6 text-center text-gray-500 bg-gray-50 rounded-lg h-64 flex items-center justify-center">
        No click data for this period.
      </div>
    );
  }

  return (
    <ResponsiveContainer width="100%" height={256}>
      <AreaChart data={chartData} margin={{ top: 5, right: 20, left: 0, bottom: 5 }}>
        <defs>
          <linearGradient id="colorClicks" x1="0" y1="0" x2="0" y2="1">
            <stop offset="5%" stopColor="#6366f1" stopOpacity={0.3} />
            <stop offset="95%" stopColor="#6366f1" stopOpacity={0} />
          </linearGradient>
        </defs>
        <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
        <XAxis
          dataKey="date"
          tick={{ fontSize: 12, fill: '#6b7280' }}
          tickFormatter={(date: string) => {
            const d = new Date(date);
            return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
          }}
        />
        <YAxis allowDecimals={false} tick={{ fontSize: 12, fill: '#6b7280' }} />
      <Tooltip
  contentStyle={{ borderRadius: '8px', border: '1px solid #e5e7eb' }}
  labelFormatter={(label: any) => {
    const date = new Date(label);
    if (isNaN(date.getTime())) return String(label); 
    return date.toLocaleDateString('en-US', { 
      weekday: 'short', 
      month: 'short', 
      day: 'numeric' 
    });
  }}
/>

        <Area
          type="monotone"
          dataKey="clicks"
          stroke="#6366f1"
          strokeWidth={2}
          fill="url(#colorClicks)"
        />
      </AreaChart>
    </ResponsiveContainer>
  );
};

export default ClickChart;