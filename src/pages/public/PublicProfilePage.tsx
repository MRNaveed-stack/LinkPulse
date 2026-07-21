import { useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { usePublicProfile } from '../../hooks/usePublicProfile';
import ProfileHeader from '../../components/public/ProfileHeader';
import LinkButton from '../../components/public/LinkButton';
import ProfileSkeleton from '../../components/public/ProfileSkeleton';
import EmptyProfile from '../../components/public/EmptyProfile';
import ErrorPage from '../../components/public/ErrorPage';

const PublicProfilePage = () => {
  const { username } = useParams<{ username: string }>();
  const { data: profile, isLoading, isError, error } = usePublicProfile(username || '');

  useEffect(() => {
    if (profile) {
      document.title = `${profile.display_name || profile.username} | LinkPulse`;
    } else if (!isLoading && !isError) {
      document.title = 'LinkPulse';
    }
    return () => {
      document.title = 'LinkPulse';
    };
  }, [profile, isLoading, isError]);

  if (isLoading) {
    return <ProfileSkeleton />;
  }

  if (isError) {
    const axiosError = error as any;
    if (axiosError?.response?.status === 404) {
      return <ErrorPage status={404} message="Profile not found" />;
    }
    return <ErrorPage status={500} message="Something went wrong. Please try again." />;
  }

  if (!profile) {
    return <ErrorPage />;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-md mx-auto px-4 py-12">
        <ProfileHeader profile={profile} />

        {/* Link Buttons */}
        <div className="space-y-3">
          {profile.links.length > 0 ? (
            profile.links.map((link) => (
              <LinkButton key={link.slug} link={link} username={profile.username} />
            ))
          ) : (
            <EmptyProfile />
          )}
        </div>

        {/* Branding Footer */}
        <div className="mt-10 text-center">
          <p className="text-xs text-gray-400">
            Powered by{' '}
            <a
              href="/"
              className="text-indigo-500 hover:text-indigo-600 font-medium"
            >
              LinkPulse
            </a>
          </p>
        </div>
      </div>
    </div>
  );
};

export default PublicProfilePage;