import {Fragment} from 'react';
import {Menu, Transition} from '@headlessui/react';
import {UserCircle, Settings, LogOut} from 'lucide-react';
import {useAuthStore} from '../store/authStore';
import {useNavigate} from 'react-router-dom';
import {useCurrentUser} from '../hooks/useCurrentUser';

const UserMenu = () => {
    const {user, logout} = useAuthStore();
    const navigate = useNavigate();
    useCurrentUser();

    const handleLogout = () => {
        logout();
        navigate('/login');
};

const avatarUrl = user?.avatar || null;
const initials = user?.username?.charAt(0).toUpperCase() || 'U';

return (
    <Menu as="div" className="relative">
      <Menu.Button className="max-w-xs bg-white flex items-center text-sm rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
        <span className="sr-only">Open user menu</span>
        {avatarUrl ? (
          <img className="h-8 w-8 rounded-full" src={avatarUrl} alt="" />
        ) : (
          <span className="inline-flex items-center justify-center h-8 w-8 rounded-full bg-indigo-500 text-white">
            {initials}
          </span>
        )}
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
        <Menu.Items className="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg py-1 bg-white ring-1 ring-black ring-opacity-5 focus:outline-none">
          <div className="px-4 py-2 text-sm text-gray-700 border-b">
            <p className="font-medium">{user?.username}</p>
            <p className="text-gray-500">{user?.email}</p>
          </div>

          <Menu.Item>
            {({ active }: { active: boolean }) => (
              <button
                onClick={() => navigate('/profile')}
                className={`${
                  active ? 'bg-gray-100' : ''
                } flex w-full items-center px-4 py-2 text-sm text-gray-700`}
              >
                <UserCircle className="mr-3 h-4 w-4" /> Profile
              </button>
            )}
          </Menu.Item>
          <Menu.Item>
            {({ active }: { active: boolean }) => (
              <button
                onClick={() => navigate('/settings')}
                className={`${
                  active ? 'bg-gray-100' : ''
                } flex w-full items-center px-4 py-2 text-sm text-gray-700`}
              >
                <Settings className="mr-3 h-4 w-4" /> Settings
              </button>
            )}
          </Menu.Item>
          <Menu.Item>
            {({ active }: { active: boolean }) => (
              <button
                onClick={handleLogout}
                className={`${
                  active ? 'bg-gray-100' : ''
                } flex w-full items-center px-4 py-2 text-sm text-gray-700`}
              >
                <LogOut className="mr-3 h-4 w-4" /> Logout
              </button>
            )}
          </Menu.Item>
        </Menu.Items>
      </Transition>
    </Menu>
  );
};

export default UserMenu;