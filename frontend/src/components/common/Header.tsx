import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { ShoppingCart, User, Search, Menu, X, Heart, LogOut } from 'lucide-react'
import { useSelector, useDispatch } from 'react-redux'
import { RootState, AppDispatch } from '@store/store'
import { logout } from '@store/slices/authSlice'
import { motion, AnimatePresence } from 'framer-motion'

function Header() {
  const dispatch = useDispatch<AppDispatch>()
  const navigate = useNavigate()
  const { isAuthenticated, user } = useSelector((state: RootState) => state.auth)
  const { cart } = useSelector((state: RootState) => state.cart)
  const { categories } = useSelector((state: RootState) => state.products)
  
  const [searchQuery, setSearchQuery] = useState('')
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const [showUserMenu, setShowUserMenu] = useState(false)

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    if (searchQuery.trim()) {
      navigate(`/products?search=${encodeURIComponent(searchQuery)}`)
    }
  }

  const handleLogout = () => {
    dispatch(logout())
    setShowUserMenu(false)
    navigate('/')
  }

  const cartItemsCount = cart?.total_items || 0

  return (
    <header className="sticky top-0 z-50 bg-white shadow-md">
      {/* Top Bar */}
      <div className="bg-secondary-800 text-white py-2">
        <div className="container mx-auto px-4 flex justify-between items-center text-sm">
          <div className="flex gap-4">
            <span>📞 Hotline: 1900 xxxx</span>
            <span>📧 support@eshop.com</span>
          </div>
          <div className="hidden md:flex gap-4">
            <Link to="/help" className="hover:text-primary-500 transition">Help Center</Link>
            <Link to="/track-order" className="hover:text-primary-500 transition">Track Order</Link>
          </div>
        </div>
      </div>

      {/* Main Header */}
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between gap-4">
          {/* Logo */}
          <Link to="/" className="flex items-center gap-2 flex-shrink-0">
            <div className="w-10 h-10 bg-gradient-to-br from-primary-500 to-primary-700 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-xl">E</span>
            </div>
            <span className="text-2xl font-bold text-secondary-800 hidden sm:block">E-Shop</span>
          </Link>

          {/* Search Bar */}
          <form onSubmit={handleSearch} className="flex-1 max-w-2xl hidden md:block">
            <div className="relative">
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search for products..."
                className="w-full pl-12 pr-4 py-3 rounded-full border-2 border-gray-200 focus:border-primary-500 focus:ring-4 focus:ring-primary-100 outline-none transition-all"
              />
              <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
              <button
                type="submit"
                className="absolute right-2 top-1/2 -translate-y-1/2 bg-primary-500 text-white px-6 py-2 rounded-full hover:bg-primary-600 transition"
              >
                Search
              </button>
            </div>
          </form>

          {/* Actions */}
          <div className="flex items-center gap-2 sm:gap-4">
            {/* Mobile Menu Toggle */}
            <button
              onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
              className="md:hidden p-2 hover:bg-gray-100 rounded-lg transition"
            >
              {isMobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
            </button>

            {/* Wishlist */}
            <Link to="/account/wishlist" className="relative p-2 hover:bg-gray-100 rounded-lg transition hidden sm:block">
              <Heart className="w-6 h-6 text-gray-600" />
            </Link>

            {/* Cart */}
            <Link to="/cart" className="relative p-2 hover:bg-gray-100 rounded-lg transition">
              <ShoppingCart className="w-6 h-6 text-gray-600" />
              {cartItemsCount > 0 && (
                <span className="absolute -top-1 -right-1 bg-primary-500 text-white text-xs font-bold w-5 h-5 rounded-full flex items-center justify-center">
                  {cartItemsCount > 99 ? '99+' : cartItemsCount}
                </span>
              )}
            </Link>

            {/* User Menu */}
            {isAuthenticated ? (
              <div className="relative">
                <button
                  onClick={() => setShowUserMenu(!showUserMenu)}
                  className="flex items-center gap-2 p-2 hover:bg-gray-100 rounded-lg transition"
                >
                  <div className="w-8 h-8 bg-primary-100 rounded-full flex items-center justify-center">
                    <User className="w-5 h-5 text-primary-500" />
                  </div>
                  <span className="hidden lg:block text-sm font-medium">{user?.first_name}</span>
                </button>

                <AnimatePresence>
                  {showUserMenu && (
                    <motion.div
                      initial={{ opacity: 0, y: -10 }}
                      animate={{ opacity: 1, y: 0 }}
                      exit={{ opacity: 0, y: -10 }}
                      className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-xl border border-gray-100 py-2 z-50"
                    >
                      <Link
                        to="/account"
                        className="block px-4 py-2 hover:bg-gray-100 transition"
                        onClick={() => setShowUserMenu(false)}
                      >
                        👤 My Account
                      </Link>
                      <Link
                        to="/account/orders"
                        className="block px-4 py-2 hover:bg-gray-100 transition"
                        onClick={() => setShowUserMenu(false)}
                      >
                        📦 My Orders
                      </Link>
                      <button
                        onClick={handleLogout}
                        className="w-full text-left px-4 py-2 hover:bg-gray-100 transition text-red-500"
                      >
                        <LogOut className="inline w-4 h-4 mr-2" />
                        Logout
                      </button>
                    </motion.div>
                  )}
                </AnimatePresence>
              </div>
            ) : (
              <Link
                to="/login"
                className="flex items-center gap-2 bg-primary-500 text-white px-4 py-2 rounded-lg hover:bg-primary-600 transition hidden sm:flex"
              >
                <User className="w-5 h-5" />
                <span className="font-medium">Sign In</span>
              </Link>
            )}
          </div>
        </div>

        {/* Mobile Search */}
        <form onSubmit={handleSearch} className="mt-4 md:hidden">
          <div className="relative">
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search for products..."
              className="w-full pl-10 pr-4 py-3 rounded-lg border-2 border-gray-200 focus:border-primary-500 outline-none"
            />
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
          </div>
        </form>
      </div>

      {/* Categories Bar */}
      <div className="border-t hidden md:block">
        <div className="container mx-auto px-4 py-3">
          <div className="flex items-center gap-6 overflow-x-auto">
            <Link to="/products" className="flex items-center gap-2 text-primary-500 font-semibold whitespace-nowrap">
              <Menu className="w-5 h-5" />
              All Categories
            </Link>
            {categories.slice(0, 8).map((category) => (
              <Link
                key={category.id}
                to={`/category/${category.id}`}
                className="text-gray-600 hover:text-primary-500 whitespace-nowrap transition"
              >
                {category.name}
              </Link>
            ))}
          </div>
        </div>
      </div>

      {/* Mobile Menu */}
      <AnimatePresence>
        {isMobileMenuOpen && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            className="md:hidden border-t overflow-hidden"
          >
            <div className="container mx-auto px-4 py-4 space-y-4">
              <Link to="/products" className="block py-2 text-primary-500 font-semibold">
                📂 All Categories
              </Link>
              {categories.slice(0, 6).map((category) => (
                <Link
                  key={category.id}
                  to={`/category/${category.id}`}
                  className="block py-2 text-gray-600 hover:text-primary-500"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  {category.name}
                </Link>
              ))}
              <hr />
              <Link to="/account" className="block py-2">👤 My Account</Link>
              <Link to="/account/orders" className="block py-2">📦 My Orders</Link>
              <Link to="/account/wishlist" className="block py-2">❤️ Wishlist</Link>
              {!isAuthenticated && (
                <Link to="/login" className="block py-2 text-primary-500 font-semibold">
                  Sign In
                </Link>
              )}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </header>
  )
}

export default Header
