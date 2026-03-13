import { Routes, Route } from 'react-router-dom'
import MainLayout from '@layouts/MainLayout'
import HomePage from '@pages/home/HomePage'
import ProductListPage from '@pages/product/ProductListPage'
import ProductDetailPage from '@pages/product/ProductDetailPage'
import CartPage from '@pages/cart/CartPage'
import CheckoutPage from '@pages/checkout/CheckoutPage'
import UserDashboardPage from '@pages/user/UserDashboardPage'
import AdminDashboardPage from '@pages/admin/AdminDashboardPage'
import LoginPage from '@pages/auth/LoginPage'
import HelpPage from '@pages/help/HelpPage'
import TrackOrderPage from '@pages/track-order/TrackOrderPage'
import AboutPage from '@pages/about/AboutPage'
import ContactPage from '@pages/contact/ContactPage'
import FAQPage from '@pages/faq/FAQPage'
import ShippingPage from '@pages/shipping/ShippingPage'
import TermsPage from '@pages/terms/TermsPage'
import PrivacyPage from '@pages/privacy/PrivacyPage'
import ReturnsPage from '@pages/returns/ReturnsPage'

function App() {
  return (
    <Routes>
      {/* Auth Pages (no header/footer) */}
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<LoginPage />} />
      
      {/* Main Layout with Header & Footer */}
      <Route path="/" element={<MainLayout />}>
        {/* Home Page */}
        <Route index element={<HomePage />} />
        
        {/* Product Pages */}
        <Route path="products" element={<ProductListPage />} />
        <Route path="products/:id" element={<ProductDetailPage />} />
        <Route path="category/:categoryId" element={<ProductListPage />} />
        
        {/* Cart & Checkout */}
        <Route path="cart" element={<CartPage />} />
        <Route path="checkout" element={<CheckoutPage />} />
        
        {/* User Dashboard */}
        <Route path="account" element={<UserDashboardPage />} />
        <Route path="account/orders" element={<UserDashboardPage />} />
        <Route path="account/wishlist" element={<UserDashboardPage />} />
        <Route path="account/addresses" element={<UserDashboardPage />} />
        
        {/* Admin Panel */}
        <Route path="admin" element={<AdminDashboardPage />} />
        <Route path="admin/products" element={<AdminDashboardPage />} />
        <Route path="admin/orders" element={<AdminDashboardPage />} />
        <Route path="admin/users" element={<AdminDashboardPage />} />
        <Route path="admin/analytics" element={<AdminDashboardPage />} />
        
        {/* Help & Info Pages */}
        <Route path="help" element={<HelpPage />} />
        <Route path="track-order" element={<TrackOrderPage />} />
        <Route path="about" element={<AboutPage />} />
        <Route path="contact" element={<ContactPage />} />
        <Route path="faq" element={<FAQPage />} />
        <Route path="shipping" element={<ShippingPage />} />
        <Route path="terms" element={<TermsPage />} />
        <Route path="privacy" element={<PrivacyPage />} />
        <Route path="returns" element={<ReturnsPage />} />
      </Route>
    </Routes>
  )
}

export default App
