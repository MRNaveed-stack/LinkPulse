import { useState } from 'react';
import { useCurrentUser } from '../../hooks/useCurrentUser';
import { useAnalyticsOverview } from '../../hooks/useAnalyticsOverview';
import { useAuthStore } from '../../store/authStore';
import Card from '../../components/common/Card';
import Skeleton from '../../components/common/Skeleton';
import ChangePasswordForm from '../../components/profile/ChangePasswordForm';
import DeleteAccountDialog from '../../components/profile/DeleteAccountDialog';
import { Shield, CreditCard, Mail, Trash2, Calendar, Link as LinkIcon, MousePointerClick } from 'lucide-react';

export default function SettingsPage() {
  const { isLoading: userLoading } = useCurrentUser();
  const user = useAuthStore((state) => state.user);
  
  const { data: analytics, isLoading: analyticsLoading } = useAnalyticsOverview();
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);

  const joinedDate = user?.created_at
    ? new Date(user.created_at).toLocaleDateString(undefined, {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      })
    : 'N/A';

  if (userLoading) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-8">
        {/* Page Title Skeleton */}
        <div className="mb-8">
          <Skeleton className="h-8 w-64 mb-2" />
          <Skeleton className="h-4 w-96" />
        </div>

        <div className="space-y-8">
          {/* Card: Account details */}
          <Card className="p-6 bg-white border border-gray-100 shadow-sm">
            <Skeleton className="h-6 w-48 mb-6" />
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
              {[...Array(4)].map((_, i) => (
                <div key={i} className="flex items-center space-x-3 bg-gray-50 p-4 rounded-lg border border-gray-200">
                  <Skeleton className="h-5 w-5 rounded-full" />
                  <div className="flex-1">
                    <Skeleton className="h-3 w-16 mb-1" />
                    <Skeleton className="h-4 w-28" />
                  </div>
                </div>
              ))}
            </div>
          </Card>

          {/* Card: Security */}
          <Card className="p-6 bg-white border border-gray-100 shadow-sm">
            <Skeleton className="h-6 w-48 mb-6" />
            <div className="space-y-4">
              {[...Array(3)].map((_, i) => (
                <div key={i}>
                  <Skeleton className="h-4 w-28 mb-2" />
                  <Skeleton className="h-10 w-full" />
                </div>
              ))}
              <div className="flex justify-end pt-4">
                <Skeleton className="h-10 w-32" />
              </div>
            </div>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      {/* Page Title */}
      <div className="mb-8">
        <h2 className="text-3xl font-extrabold text-gray-900 tracking-tight">Account Settings</h2>
        <p className="mt-2 text-sm text-gray-600">
          Manage your subscription plan, security credentials, and danger zone configurations.
        </p>
      </div>

      <div className="space-y-8">
        {/* Section: Account Information */}
        <Card className="border border-gray-100 shadow-sm bg-white">
          <div className="p-6">
            <h3 className="text-lg font-semibold text-gray-900 border-b pb-3 mb-6 flex items-center">
              <Mail className="h-5 w-5 text-indigo-500 mr-2" />
              Account details & stats
            </h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
              <div className="flex items-center space-x-3 bg-gray-50 p-4 rounded-lg border border-gray-200">
                <Calendar className="h-5 w-5 text-gray-400" />
                <div>
                  <span className="block text-xs font-medium text-gray-400 uppercase tracking-wider">Joined Date</span>
                  <span className="text-sm font-semibold text-gray-800">{joinedDate}</span>
                </div>
              </div>

              <div className="flex items-center space-x-3 bg-gray-50 p-4 rounded-lg border border-gray-200">
                <CreditCard className="h-5 w-5 text-gray-400" />
                <div>
                  <span className="block text-xs font-medium text-gray-400 uppercase tracking-wider">Active Subscription</span>
                  <span className="text-sm font-semibold text-gray-800 capitalize">{user?.plan || 'Free'} Plan</span>
                </div>
              </div>

              <div className="flex items-center space-x-3 bg-gray-50 p-4 rounded-lg border border-gray-200">
                <LinkIcon className="h-5 w-5 text-gray-400" />
                <div>
                  <span className="block text-xs font-medium text-gray-400 uppercase tracking-wider">Links Created</span>
                  <span className="text-sm font-semibold text-gray-800">
                    {analyticsLoading ? '...' : analytics?.totalLinks ?? 0}
                  </span>
                </div>
              </div>

              <div className="flex items-center space-x-3 bg-gray-50 p-4 rounded-lg border border-gray-200">
                <MousePointerClick className="h-5 w-5 text-gray-400" />
                <div>
                  <span className="block text-xs font-medium text-gray-400 uppercase tracking-wider">Total Clicks</span>
                  <span className="text-sm font-semibold text-gray-800">
                    {analyticsLoading ? '...' : analytics?.totalClicks ?? 0}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </Card>

        {/* Section: Security */}
        <Card className="border border-gray-100 shadow-sm bg-white">
          <div className="p-6">
            <h3 className="text-lg font-semibold text-gray-900 border-b pb-3 mb-6 flex items-center">
              <Shield className="h-5 w-5 text-indigo-500 mr-2" />
              Security Credentials
            </h3>
            <ChangePasswordForm />
          </div>
        </Card>

        {/* Section: Danger Zone */}
        <Card className="border border-red-200 shadow-sm bg-red-50/50">
          <div className="p-6">
            <h3 className="text-lg font-semibold text-red-900 border-b border-red-200 pb-3 mb-6 flex items-center">
              <Trash2 className="h-5 w-5 text-red-600 mr-2" />
              Danger Zone
            </h3>
            <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
              <div>
                <h4 className="text-sm font-bold text-gray-900">Delete Account Permanently</h4>
                <p className="mt-1 text-xs text-gray-500">
                  Once deleted, your profile data, landing page, links, and click analytics cannot be recovered.
                </p>
              </div>
              <button
                onClick={() => setIsDeleteOpen(true)}
                className="w-full sm:w-auto inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 font-semibold transition"
              >
                Delete Account
              </button>
            </div>
          </div>
        </Card>
      </div>

      {/* Delete Account Modal Dialog */}
      <DeleteAccountDialog open={isDeleteOpen} onClose={() => setIsDeleteOpen(false)} />
    </div>
  );
}
