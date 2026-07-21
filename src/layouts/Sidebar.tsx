import { NavLink, useNavigate } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import {
  LayoutDashboard,
  Link,
  BarChart3,
  User,
  Settings,
  LogOut
} from 'lucide-react';

interface SidebarProps {
  mobile?: boolean;
  onClose?: () => void;
}

const navigation = [
  { name: 'Dashboard', path: '/dashboard', icon: LayoutDashboard },
  { name: 'Links', path: '/links', icon: Link },
  { name: 'Analytics', path: '/analytics', icon: BarChart3 },
  { name: 'Profile', path: '/profile', icon: User },
  { name: 'Settings', path: '/settings', icon: Settings },
];

const Sidebar = ({ mobile, onClose }: SidebarProps) => {
  const logout = useAuthStore((state) => state.logout);
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="h-full flex flex-col bg-gray-900 text-white">
      <div className="flex items-center h-16 flex-shrink-0 px-4 bg-gray-800">
        <h1 className="text-xl font-bold tracking-wider">LinkPulse</h1>
      </div>
      <div className="flex-1 flex flex-col overflow-y-auto">
        <nav className="flex-1 px-2 py-4 space-y-1">
          {navigation.map((item) => (
            <NavLink
              key={item.name}
              to={item.path}
              onClick={mobile ? onClose : undefined}
              aria-label={`Go to ${item.name} page`}
              className={({ isActive }) =>
                `group flex items-center px-2 py-2 text-sm font-medium rounded-md transition duration-150 ${
                  isActive ? 'bg-gray-800 text-white font-semibold' : 'text-gray-300 hover:bg-gray-700 hover:text-white'
                }`
              }
            >
              <item.icon className="mr-3 h-5 w-5 text-gray-400 group-hover:text-white transition" />
              {item.name}
            </NavLink>
          ))}
        </nav>
        <div className="px-2 py-4 border-t border-gray-700">
          {/* Logout */}
          <button
            onClick={handleLogout}
            aria-label="Log out of LinkPulse"
            className="w-full group flex items-center px-2 py-2 text-sm font-medium rounded-md text-gray-300 hover:bg-gray-700 hover:text-white transition cursor-pointer"
          >
            <LogOut className="mr-3 h-5 w-5 text-gray-400 group-hover:text-white" />
            Logout
          </button>
        </div>
      </div>
    </div>
  );
};

export default Sidebar;
