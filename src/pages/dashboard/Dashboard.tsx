import { useCurrentUser } from '../../hooks/useCurrentUser';
import { useAnalyticsOverview } from '../../hooks/useAnalyticsOverview';
import Card from '../../components/common/Card';
import Skeleton from '../../components/common/Skeleton';
import Spinner from '../../components/common/Spinner';
import EmptyState from '../../components/common/EmptyState';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../store/authStore';
import { UserCircle } from 'lucide-react';
import Button from '../../components/common/Button';

import { BarChart3, Link, MousePointerClick, TrendingUp } from 'lucide-react';

const Dashboard = () => {
  const navigate = useNavigate();
  const { isLoading: userLoading } = useCurrentUser(); // If you use useQuery directly, adjust

  const userFromStore = useAuthStore((state) => state.user);

  const { data: overview, isLoading, isError, error } = useAnalyticsOverview();

  // If user data not in store yet, we might show a loading state
  if (!userFromStore && userLoading) {
    return (
      <div className="flex justify-center py-12">
        <Spinner size="lg" />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="py-12 text-center">
        <h2 className="text-lg font-medium text-red-600">Failed to load analytics</h2>
        <p className="text-sm text-gray-500 mt-2">{(error as any)?.message || 'An unexpected error occurred'}</p>
      </div>
    );
  }

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900 mb-6">Dashboard</h2>

      {/* Overview Cards */}
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
        <Card>
          <div className="px-4 py-5 sm:p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0 bg-indigo-500 rounded-md p-3">
                <MousePointerClick className="h-6 w-6 text-white" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">Total Clicks</dt>
                  <dd className="text-lg font-semibold text-gray-900">
                    {isLoading ? <Skeleton className="h-6 w-16" /> : overview?.totalClicks ?? 0}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </Card>

        <Card>
          <div className="px-4 py-5 sm:p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0 bg-green-500 rounded-md p-3">
                <Link className="h-6 w-6 text-white" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">Total Links</dt>
                  <dd className="text-lg font-semibold text-gray-900">
                    {isLoading ? <Skeleton className="h-6 w-16" /> : overview?.totalLinks ?? 0}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </Card>

        <Card>
          <div className="px-4 py-5 sm:p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0 bg-yellow-500 rounded-md p-3">
                <BarChart3 className="h-6 w-6 text-white" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">Active Links</dt>
                  <dd className="text-lg font-semibold text-gray-900">
                    {isLoading ? <Skeleton className="h-6 w-16" /> : overview?.activeLinks ?? 0}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </Card>

        <Card>
          <div className="px-4 py-5 sm:p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0 bg-purple-500 rounded-md p-3">
                <TrendingUp className="h-6 w-6 text-white" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">Top Link</dt>
                  <dd className="text-lg font-semibold text-gray-900">
                    {isLoading ? (
                      <Skeleton className="h-6 w-24" />
                    ) : overview?.topPerformingLink ? (
                      <span className="truncate block">{overview.topPerformingLink.title}</span>
                    ) : (
                      'N/A'
                    )}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </Card>
      </div>

      {/* User Profile Card */}
      <div className="mt-8">
        <Card>
          <div className="px-4 py-5 sm:p-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Profile</h3>
            {userFromStore ? (
              <div className="flex items-center space-x-4">
                <span className="inline-flex items-center justify-center h-12 w-12 rounded-full bg-indigo-500 text-white text-lg font-bold">
                  {(userFromStore.username || 'User').charAt(0).toUpperCase()}
                </span>
                <div>
                  <p className="text-sm font-medium text-gray-900">{userFromStore.username || 'User'}</p>
                  <p className="text-sm text-gray-500">{userFromStore.email}</p>
                </div>
              </div>
            ) : (
              <EmptyState title="User info not available" icon={<UserCircle className="h-8 w-8" />} />
            )}
          </div>
        </Card>
      </div>

      {/* If analytics empty (but loaded) */}
      {!isLoading && overview && overview.totalLinks === 0 && (
        <div className="mt-8">
          <EmptyState
            title="No links yet"
            description="Create your first short link to see analytics."
            action={
              <Button onClick={() => navigate('/links/new')}>Create a link</Button>
            }
          />
        </div>
      )}
    </div>
  );
};

export default Dashboard;