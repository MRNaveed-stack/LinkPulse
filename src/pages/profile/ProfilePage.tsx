import { useState, useEffect } from 'react';
import { useCurrentUser } from '../../hooks/useCurrentUser';
import { usePublicProfile } from '../../hooks/usePublicProfile';
import { useLinks } from '../../hooks/useLinks';
import { useUpdateProfile } from '../../hooks/useUpdateProfile';
import { useAuthStore } from '../../store/authStore';
import Card from '../../components/common/Card';
import Skeleton from '../../components/common/Skeleton';
import ProfileForm from '../../components/profile/ProfileForm';
import ProfilePreview from '../../components/profile/ProfilePreview';

export default function ProfilePage() {
  const { isLoading: userLoading } = useCurrentUser();
  const user = useAuthStore((state) => state.user);
  
  const username = user?.username || '';
  const { data: profile, isLoading: profileLoading, error: profileError } = usePublicProfile(username);
  const { data: links = [] } = useLinks();
  const updateProfileMutation = useUpdateProfile();

  const [previewData, setPreviewData] = useState({
    display_name: '',
    username: '',
    bio: '',
    avatar_url: '',
  });

  // Populate preview data once loaded
  useEffect(() => {
    if (profile) {
      setPreviewData({
        display_name: profile.display_name || '',
        username: username,
        bio: profile.bio || '',
        avatar_url: profile.avatar_url || '',
      });
    } else if (user) {
      setPreviewData({
        display_name: user.username || '',
        username: username,
        bio: '',
        avatar_url: user.avatar || '',
      });
    }
  }, [profile, user, username]);

  const handleProfileSubmit = (data: any) => {
    updateProfileMutation.mutate({
      display_name: data.display_name,
      username: data.username,
      email: data.email,
      bio: data.bio,
      avatar_url: data.avatar_url,
    });
  };

  const isProfile404 = (profileError as any)?.response?.status === 404;

  if (userLoading || (profileLoading && !isProfile404)) {
    return (
      <div className="max-w-6xl mx-auto px-4 py-8">
        {/* Page Header Skeleton */}
        <div className="mb-8">
          <Skeleton className="h-8 w-64 mb-2" />
          <Skeleton className="h-4 w-96" />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 items-start">
          {/* Profile Edit Form Skeleton */}
          <div className="lg:col-span-2">
            <Card className="p-6 bg-white border border-gray-100 shadow-sm">
              <Skeleton className="h-6 w-40 mb-6" />
              <div className="space-y-6">
                {[...Array(5)].map((_, i) => (
                  <div key={i}>
                    <Skeleton className="h-4 w-24 mb-2" />
                    <Skeleton className="h-10 w-full" />
                  </div>
                ))}
                <div className="flex justify-end pt-4 border-t border-gray-100">
                  <Skeleton className="h-10 w-32" />
                </div>
              </div>
            </Card>
          </div>

          {/* Live Preview Skeleton */}
          <div className="lg:col-span-1">
            <Card className="p-6 bg-white border border-gray-100 shadow-sm text-center flex flex-col items-center">
              <Skeleton className="h-24 w-24 rounded-full mb-4" />
              <Skeleton className="h-6 w-36 mb-2" />
              <Skeleton className="h-4 w-48 mb-6" />
              <div className="space-y-3 w-full">
                {[...Array(3)].map((_, i) => (
                  <Skeleton key={i} className="h-12 w-full" />
                ))}
              </div>
            </Card>
          </div>
        </div>
      </div>
    );
  }

  const initialProfileData = {
    display_name: previewData.display_name,
    username: username,
    email: user?.email || '',
    bio: previewData.bio,
    avatar_url: previewData.avatar_url,
  };

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      {/* Page Header */}
      <div className="mb-8">
        <h2 className="text-3xl font-extrabold text-gray-900 tracking-tight">Profile Customization</h2>
        <p className="mt-2 text-sm text-gray-600">
          Personalize your public LinkPulse bio landing page and update your avatar.
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 items-start">
        {/* Profile Edit Form */}
        <div className="lg:col-span-2">
          <Card className="border border-gray-100 shadow-sm bg-white overflow-visible">
            <div className="p-6">
              <h3 className="text-lg font-semibold text-gray-900 border-b pb-3 mb-6">
                Edit Profile Info
              </h3>
              {user && (
                <ProfileForm
                  initialData={initialProfileData}
                  onSubmit={handleProfileSubmit}
                  isLoading={updateProfileMutation.isPending}
                  onFormChange={setPreviewData} // Stable reference passed here to prevent loop
                />
              )}
            </div>
          </Card>
        </div>

        {/* Live Preview Column */}
        <div className="lg:col-span-1">
          <ProfilePreview data={previewData} links={links} />
        </div>
      </div>
    </div>
  );
}
