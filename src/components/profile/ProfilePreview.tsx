import type { Link } from '../../types/link';
import { Eye, ExternalLink } from 'lucide-react';

interface ProfilePreviewProps {
  data: {
    display_name: string;
    username: string;
    bio: string;
    avatar_url: string;
  };
  links: Link[];
}

export default function ProfilePreview({ data, links }: ProfilePreviewProps) {
  const { display_name, username, bio, avatar_url } = data;
  const initials = display_name
    ? display_name.charAt(0).toUpperCase()
    : username.charAt(0).toUpperCase() || 'U';

  const activeLinks = links ? links.filter((l) => l.is_active) : [];

  return (
    <div className="sticky top-6">
      <div className="flex items-center space-x-2 text-sm font-semibold text-gray-700 mb-3 px-1">
        <Eye className="h-4 w-4 text-indigo-500" />
        <span>Public Profile Preview</span>
      </div>

      <div className="border border-gray-200 shadow-xl bg-white overflow-hidden rounded-2xl max-w-sm mx-auto">
        <div className="bg-gradient-to-r from-indigo-500 via-purple-500 to-pink-500 h-24 w-full"></div>
        <div className="px-6 pb-8 pt-0 flex flex-col items-center -mt-12">
          {/* Avatar Rendering */}
          {avatar_url ? (
            <img
              className="h-24 w-24 rounded-full border-4 border-white shadow-md object-cover bg-gray-50"
              src={avatar_url}
              alt={display_name || username}
              onError={(e) => {
                (e.target as HTMLElement).style.display = 'none';
              }}
            />
          ) : (
            <span className="inline-flex items-center justify-center h-24 w-24 rounded-full border-4 border-white shadow-md bg-indigo-600 text-white text-3xl font-bold">
              {initials}
            </span>
          )}

          {/* User Info */}
          <h4 className="mt-4 text-xl font-bold text-gray-900 text-center truncate w-full px-2">
            {display_name || username || 'Muhammad Naveed'}
          </h4>
          
          {username && (
            <p className="text-xs text-gray-400 font-medium mt-0.5">@{username}</p>
          )}

          {/* Bio */}
          <p className="mt-3.5 text-sm text-gray-500 text-center max-w-xs leading-relaxed italic px-2">
            {bio || "This user hasn't written a bio yet."}
          </p>

          {/* Links Render */}
          <div className="w-full mt-6 space-y-3 max-h-56 overflow-y-auto px-1">
            {activeLinks.length > 0 ? (
              activeLinks.map((link) => (
                <a
                  key={link.id}
                  href={`/u/${username}/${link.slug}`}
                  target="_blank"
                  rel="noreferrer"
                  className="group w-full flex items-center justify-between bg-gray-50 border border-gray-200 rounded-lg py-2.5 px-4 text-sm font-medium text-gray-700 hover:bg-gray-100 hover:border-gray-300 transition duration-150 shadow-sm"
                >
                  <span className="truncate pr-2">{link.title}</span>
                  <ExternalLink className="h-3.5 w-3.5 text-gray-400 group-hover:text-indigo-500 flex-shrink-0" />
                </a>
              ))
            ) : (
              <div className="text-center py-6 text-xs text-gray-400 border border-dashed border-gray-200 rounded-lg">
                No active links to display.
              </div>
            )}
          </div>

          <div className="mt-6 text-center border-t border-gray-100 w-full pt-4">
            <a
              href={`/u/${username}`}
              target="_blank"
              rel="noreferrer"
              className="text-xs text-indigo-500 hover:text-indigo-600 font-semibold underline flex items-center justify-center space-x-1"
            >
              <span>View Public Page</span>
            </a>
          </div>
        </div>
      </div>
    </div>
  );
}
