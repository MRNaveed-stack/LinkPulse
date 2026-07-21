import { ResponsiveContainer, PieChart, Pie, Cell, Tooltip, Legend } from 'recharts';
import { useReferrers } from '../../hooks/useReferrers';
import Skeleton from '../common/Skeleton';

const COLORS = ['#6366f1', '#8b5cf6', '#a855f7', '#d946ef', '#ec4899', '#f43f5e', '#f97316', '#eab308'];

const ReferrerChart = () => {
  const { data, isLoading, isError } = useReferrers();

  if (isError) {
    return (
      <div className="p-6 text-center text-red-600 bg-red-50 rounded-lg">
        Unable to load referrer data.
      </div>
    );
  }

  if (isLoading) {
    return <Skeleton className="h-64 w-full rounded-lg" />;
  }

  const referrers = data || [];
  const hasData = referrers.length > 0;

  if (!hasData) {
    return (
      <div className="p-6 text-center text-gray-500 bg-gray-50 rounded-lg h-64 flex items-center justify-center">
        No referrer data yet.
      </div>
    );
  }

  return (
    <ResponsiveContainer width="100%" height={256}>
      <PieChart>
        <Pie
          data={referrers}
          dataKey="clicks"
          nameKey="referrer"
          cx="50%"
          cy="50%"
          outerRadius={80}
          innerRadius={50}
          paddingAngle={3}
          label={({ referrer, percent }: any) => `${referrer} ${(percent * 100).toFixed(0)}%`}
        >
          {referrers.map((_, index) => (
            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
          ))}
        </Pie>
        <Tooltip />
        <Legend />
      </PieChart>
    </ResponsiveContainer>
  );
};

export default ReferrerChart;
