import type { PublicProfile } from '../../types/public';

interface ProfileHeaderProps {
  profile: PublicProfile;
}

const ProfileHeader = ({ profile }: ProfileHeaderProps) => {
  const initials = profile.display_name
    ? profile.display_name.charAt(0).toUpperCase()
    : profile.username.charAt(0).toUpperCase();

  return (
    <div className="flex flex-col items-center mb-8">
      {profile.avatar_url ? (
        <img
          src={profile.avatar_url}
          alt={profile.display_name || profile.username}
          className="h-24 w-24 rounded-full object-cover border-2 border-gray-200 shadow-sm mb-4"
        />
      ) : (
        <div className="h-24 w-24 rounded-full bg-indigo-500 flex items-center justify-center text-white text-3xl font-bold shadow-sm mb-4">
          {initials}
        </div>
      )}

      <h1 className="text-xl font-bold text-gray-900">
        {profile.display_name || profile.username}
      </h1>

      {profile.display_name && (
        <p className="text-sm text-gray-500 mt-1">@{profile.username}</p>
      )}

      {profile.bio && (
        <p className="text-sm text-gray-600 mt-3 text-center max-w-sm">{profile.bio}</p>
      )}
    </div>
  );
};

export default ProfileHeader;