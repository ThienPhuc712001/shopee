import React, { useEffect } from 'react';
import { Link, Outlet, useLocation, useNavigate } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../../hooks';
import { logout } from '../../store/slices/authSlice';
import { ROUTES } from '../../constants';
import {
  UserIcon,
  ShoppingBagIcon,
  ChatBubbleLeftIcon,
  MapPinIcon,
  Cog6ToothIcon,
  ArrowLeftOnRectangleIcon,
  HomeIcon,
} from '@heroicons/react/24/outline';

const UserDashboard: React.FC = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const { user, isAuthenticated } = useAppSelector((state) => state.auth);

  useEffect(() => {
    if (!isAuthenticated) {
      navigate(ROUTES.LOGIN, { state: { from: location } });
    }
  }, [isAuthenticated, navigate, location]);

  const handleLogout = async () => {
    await dispatch(logout());
    navigate(ROUTES.HOME);
  };

  const navigation = [
    { name: 'Dashboard', href: ROUTES.USER_DASHBOARD, icon: HomeIcon },
    { name: 'My Orders', href: ROUTES.USER_ORDERS, icon: ShoppingBagIcon },
    { name: 'My Reviews', href: ROUTES.USER_REVIEWS, icon: ChatBubbleLeftIcon },
    { name: 'Addresses', href: ROUTES.USER_ADDRESSES, icon: MapPinIcon },
    { name: 'Settings', href: ROUTES.USER_SETTINGS, icon: Cog6ToothIcon },
  ];

  if (!isAuthenticated || !user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <div className="flex flex-col lg:flex-row gap-6">
          {/* Sidebar */}
          <aside className="lg:w-64 flex-shrink-0">
            <div className="bg-white rounded-lg shadow-md p-6">
              {/* User info */}
              <div className="text-center mb-6 pb-6 border-b border-gray-200">
                <div className="w-20 h-20 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-3">
                  {user.avatar ? (
                    <img src={user.avatar} alt={user.first_name} className="w-full h-full rounded-full object-cover" />
                  ) : (
                    <UserIcon className="w-10 h-10 text-primary-600" />
                  )}
                </div>
                <h2 className="font-semibold text-gray-900">
                  {user.first_name} {user.last_name}
                </h2>
                <p className="text-sm text-gray-500">{user.email}</p>
              </div>

              {/* Navigation */}
              <nav className="space-y-1">
                {navigation.map((item) => {
                  const isActive = location.pathname === item.href;
                  return (
                    <Link
                      key={item.name}
                      to={item.href}
                      className={`flex items-center px-4 py-2 rounded-lg transition-colors ${
                        isActive
                          ? 'bg-primary-50 text-primary-600'
                          : 'text-gray-700 hover:bg-gray-100'
                      }`}
                    >
                      <item.icon className="w-5 h-5 mr-3" />
                      {item.name}
                    </Link>
                  );
                })}
              </nav>

              {/* Logout */}
              <div className="mt-6 pt-6 border-t border-gray-200">
                <Link
                  to={ROUTES.HOME}
                  className="flex items-center px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
                >
                  <HomeIcon className="w-5 h-5 mr-3" />
                  Back to Home
                </Link>
                <button
                  onClick={handleLogout}
                  className="w-full flex items-center px-4 py-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                >
                  <ArrowLeftOnRectangleIcon className="w-5 h-5 mr-3" />
                  Logout
                </button>
              </div>
            </div>
          </aside>

          {/* Main content */}
          <main className="flex-1">
            <Outlet />
          </main>
        </div>
      </div>
    </div>
  );
};

export default UserDashboard;
