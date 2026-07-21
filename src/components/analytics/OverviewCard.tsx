import { ReactNode } from 'react';
import Card from '../common/Card';
import Skeleton from '../common/Skeleton';

interface OverviewCardProps {
  title: string;
  value: number | string;
  icon: ReactNode;
  loading?: boolean;
  className?: string;
}

const OverviewCard = ({ title, value, icon, loading, className }: OverviewCardProps) => (
  <Card className={`p-5 ${className || ''}`}>
    <div className="flex items-center">
      <div className="flex-shrink-0">{icon}</div>
      <div className="ml-5 w-0 flex-1">
        <dt className="text-sm font-medium text-gray-500 truncate">{title}</dt>
        <dd className="text-2xl font-bold text-gray-900">
          {loading ? <Skeleton className="h-8 w-20" /> : value}
        </dd>
      </div>
    </div>
  </Card>
);

export default OverviewCard;