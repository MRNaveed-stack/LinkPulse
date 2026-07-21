import { useState, Fragment } from 'react';
import { Menu as MenuIcon, Share2, Copy, ExternalLink, Check } from 'lucide-react';
import UserMenu from './UserMenu';
import { useAuthStore } from '../store/authStore';
import { Menu, Transition } from '@headlessui/react';
import { toast } from 'react-hot-toast';

interface NavbarProps {
  onMenuClick: () => void;
}

const Navbar = ({ onMenuClick }: NavbarProps) => {
  const user = useAuthStore((state) => state.user);
  const username = user?.username || '';
  const profileUrl = `${window.location.origin}/u/${username}`;
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(profileUrl);
    setCopied(true);
    toast.success('Profile URL copied!');
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="relative z-10 flex-shrink-0 flex h-16 bg-white shadow">
      <button
        className="px-4 border-r border-gray-200 text-gray-500 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500 md:hidden"
        onClick={onMenuClick}
      >
        <span className="sr-only">Open sidebar</span>
        <MenuIcon className="h-6 w-6" />
      </button>
      <div className="flex-1 px-4 flex justify-between">
        <div className="flex-1 flex items-center">
          {/* Page title or search could go here */}
        </div>
        <div className="ml-4 flex items-center md:ml-6 space-x-4">
          {username && (
            <Menu as="div" className="relative">
              <Menu.Button className="inline-flex items-center px-3.5 py-1.5 border border-gray-300 rounded-full text-xs font-semibold text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 shadow-sm cursor-pointer transition">
                <Share2 className="h-3.5 w-3.5 mr-1.5 text-indigo-600" />
                <span>Share</span>
              </Menu.Button>
              <Transition
                as={Fragment}
                enter="transition ease-out duration-100"
                enterFrom="transform opacity-0 scale-95"
                enterTo="transform opacity-100 scale-100"
                leave="transition ease-in duration-75"
                leaveFrom="transform opacity-100 scale-100"
                leaveTo="transform opacity-0 scale-95"
              >
                <Menu.Items className="origin-top-right absolute right-0 mt-2 w-56 rounded-md shadow-lg py-1 bg-white ring-1 ring-black ring-opacity-5 focus:outline-none z-50">
                  <div className="px-4 py-2 border-b">
                    <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider">My Public Page</p>
                    <p className="text-xs text-gray-500 truncate mt-0.5">{profileUrl}</p>
                  </div>
                  <Menu.Item>
                    {({ active }) => (
                      <button
                        onClick={handleCopy}
                        className={`${
                          active ? 'bg-gray-100' : ''
                        } flex w-full items-center px-4 py-2 text-sm text-gray-700`}
                      >
                        {copied ? (
                          <>
                            <Check className="mr-3 h-4 w-4 text-emerald-500" />
                            <span className="text-emerald-600 font-medium">Copied!</span>
                          </>
                        ) : (
                          <>
                            <Copy className="mr-3 h-4 w-4 text-gray-400" />
                            <span>Copy link</span>
                          </>
                        )}
                      </button>
                    )}
                  </Menu.Item>
                  <Menu.Item>
                    {({ active }) => (
                      <a
                        href={profileUrl}
                        target="_blank"
                        rel="noreferrer"
                        className={`${
                          active ? 'bg-gray-100' : ''
                        } flex w-full items-center px-4 py-2 text-sm text-gray-700`}
                      >
                        <ExternalLink className="mr-3 h-4 w-4 text-gray-400" />
                        <span>Open public page</span>
                      </a>
                    )}
                  </Menu.Item>
                </Menu.Items>
              </Transition>
            </Menu>
          )}
          <UserMenu />
        </div>
      </div>
    </div>
  );
};

export default Navbar;

