import Skeleton from '../common/Skeleton';

const ProfileSkeleton = () => (
  <div className="max-w-md mx-auto px-4 py-12 flex flex-col items-center">
    {/* Avatar */}
    <Skeleton className="h-24 w-24 rounded-full mb-4" />
    {/* Display Name */}
    <Skeleton className="h-7 w-40 mb-2" />
    {/* Username */}
    <Skeleton className="h-5 w-28 mb-3" />
    {/* Bio */}
    <Skeleton className="h-4 w-64 mb-8" />
    {/* Links */}
    {[...Array(5)].map((_, i) => (
      <Skeleton key={i} className="h-14 w-full rounded-lg mb-3" />
    ))}
    {/* Footer */}
    <Skeleton className="h-4 w-36 mt-8" />
  </div>
);

export default ProfileSkeleton;