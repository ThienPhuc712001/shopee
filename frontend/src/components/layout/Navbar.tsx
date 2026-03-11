import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { ShoppingCartIcon, UserIcon, Bars3Icon, XMarkIcon, BellIcon, MagnifyingGlassIcon } from '@heroicons/react/24/outline';
import { useAppSelector } from '../../hooks';
import { ROUTES } from '../../constants';

const Navbar: React.FC = () => {
  const navigate = useNavigate();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = React.useState(false);
  const [searchQuery, setSearchQuery] = React.useState('');
  const { isAuthenticated, user } = useAppSelector((state) => state.auth);
  const { itemCount } = useAppSelector((state) => state.cart);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      navigate(`${ROUTES.PRODUCTS}?search=${encodeURIComponent(searchQuery)}`);
    }
  };

  return (
    <nav className="bg-white shadow-md sticky top-0 z-50">
      {/* Top bar */}
      <div className="bg-primary-600 text-white py-2">
        <div className="container mx-auto px-4 flex justify-between items-center text-sm">
          <div className="flex items-center space-x-4">
            <span>Welcome to E-Commerce Store!</span>
          </div>
          <div className="hidden md:flex items-center space-x-4">
            {isAuthenticated ? (
              <>
                <Link to={ROUTES.USER_ORDERS} className="hover:underline">
                  Track Order
                </Link>
                <span>|</span>
                <Link to={ROUTES.USER_DASHBOARD} className="hover:underline">
                  My Account
                </Link>
              </>
            ) : (
              <>
                <Link to={ROUTES.LOGIN} className="hover:underline">
                  Login
                </Link>
                <span>|</span>
                <Link to={ROUTES.REGISTER} className="hover:underline">
                  Sign Up
                </Link>
              </>
            )}
          </div>
        </div>
      </div>

      {/* Main navbar */}
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          {/* Logo */}
          <Link to={ROUTES.HOME} className="flex items-center space-x-2">
            <div className="w-10 h-10 bg-primary-600 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-xl">E</span>
            </div>
            <span className="text-2xl font-bold text-gray-800 hidden sm:block">E-Shop</span>
          </Link>

          {/* Search bar */}
          <form onSubmit={handleSearch} className="flex-1 max-w-2xl mx-4 hidden md:flex">
            <div className="relative w-full">
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search for products..."
                className="w-full px-4 py-2 pl-10 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
              />
              <MagnifyingGlassIcon className="w-5 h-5 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
            </div>
            <button
              type="submit"
              className="ml-2 px-6 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
            >
              Search
            </button>
          </form>

          {/* Right actions */}
          <div className="flex items-center space-x-4">
            {/* Mobile search button */}
            <button className="md:hidden p-2 hover:bg-gray-100 rounded-lg">
              <MagnifyingGlassIcon className="w-6 h-6 text-gray-600" />
            </button>

            {/* Notifications */}
            <button className="relative p-2 hover:bg-gray-100 rounded-lg hidden sm:block">
              <BellIcon className="w-6 h-6 text-gray-600" />
              <span className="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"></span>
            </button>

            {/* Cart */}
            <Link to={ROUTES.CART} className="relative p-2 hover:bg-gray-100 rounded-lg">
              <ShoppingCartIcon className="w-6 h-6 text-gray-600" />
              {itemCount > 0 && (
                <span className="absolute -top-1 -right-1 w-5 h-5 bg-primary-600 text-white text-xs rounded-full flex items-center justify-center">
                  {itemCount}
                </span>
              )}
            </Link>

            {/* User */}
            {isAuthenticated ? (
              <Link to={ROUTES.USER_DASHBOARD} className="flex items-center space-x-2 p-2 hover:bg-gray-100 rounded-lg">
                <UserIcon className="w-6 h-6 text-gray-600" />
                <span className="hidden lg:block text-sm font-medium">
                  {user?.first_name || user?.email}
                </span>
              </Link>
            ) : (
              <Link to={ROUTES.LOGIN} className="p-2 hover:bg-gray-100 rounded-lg hidden sm:block">
                <UserIcon className="w-6 h-6 text-gray-600" />
              </Link>
            )}

            {/* Mobile menu button */}
            <button
              onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
              className="md:hidden p-2 hover:bg-gray-100 rounded-lg"
            >
              {isMobileMenuOpen ? (
                <XMarkIcon className="w-6 h-6 text-gray-600" />
              ) : (
                <Bars3Icon className="w-6 h-6 text-gray-600" />
              )}
            </button>
          </div>
        </div>

        {/* Mobile search */}
        <form onSubmit={handleSearch} className="mt-4 md:hidden">
          <div className="relative">
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search for products..."
              className="w-full px-4 py-2 pl-10 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
            />
            <MagnifyingGlassIcon className="w-5 h-5 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
          </div>
        </form>
      </div>

      {/* Mobile menu */}
      {isMobileMenuOpen && (
        <div className="md:hidden border-t border-gray-200 bg-white">
          <div className="container mx-auto px-4 py-4 space-y-4">
            {isAuthenticated ? (
              <>
                <Link
                  to={ROUTES.USER_DASHBOARD}
                  className="block py-2 text-gray-700 hover:text-primary-600"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  My Account
                </Link>
                <Link
                  to={ROUTES.USER_ORDERS}
                  className="block py-2 text-gray-700 hover:text-primary-600"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  My Orders
                </Link>
              </>
            ) : (
              <>
                <Link
                  to={ROUTES.LOGIN}
                  className="block py-2 text-gray-700 hover:text-primary-600"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  Login
                </Link>
                <Link
                  to={ROUTES.REGISTER}
                  className="block py-2 text-gray-700 hover:text-primary-600"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  Sign Up
                </Link>
              </>
            )}
          </div>
        </div>
      )}
    </nav>
  );
};

export default Navbar;
