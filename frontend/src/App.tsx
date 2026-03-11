import React, { Suspense, lazy } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Provider } from 'react-redux';
import { store } from './store';
import { ROUTES } from './constants';
import LoadingSpinner from './components/common/LoadingSpinner';

// Layout components
import Navbar from './components/layout/Navbar';
import Footer from './components/layout/Footer';

// Lazy loaded pages
const HomePage = lazy(() => import('./pages/HomePage'));
const LoginPage = lazy(() => import('./pages/auth/LoginPage'));
const RegisterPage = lazy(() => import('./pages/auth/RegisterPage'));
const ProductsPage = lazy(() => import('./pages/product/ProductsPage'));
const ProductDetailPage = lazy(() => import('./pages/product/ProductDetailPage'));
const CartPage = lazy(() => import('./pages/cart/CartPage'));
const CheckoutPage = lazy(() => import('./pages/checkout/CheckoutPage'));
const UserDashboard = lazy(() => import('./pages/user/UserDashboard'));
const UserOrdersPage = lazy(() => import('./pages/user/UserOrdersPage'));
const UserProfilePage = lazy(() => import('./pages/user/UserProfilePage'));

// Layout wrapper for main pages
const MainLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />
      <main className="flex-1">{children}</main>
      <Footer />
    </div>
  );
};

// Protected route component
const ProtectedRoute: React.FC<{ children: React.ReactNode; requireAuth?: boolean }> = ({ 
  children, 
  requireAuth = true 
}) => {
  const isAuthenticated = localStorage.getItem('access_token') !== null;
  
  if (requireAuth && !isAuthenticated) {
    return <Navigate to={ROUTES.LOGIN} replace />;
  }
  
  return <>{children}</>;
};

// Loading fallback
const PageLoader: React.FC = () => (
  <div className="min-h-screen flex items-center justify-center">
    <LoadingSpinner size="lg" text="Loading page..." />
  </div>
);

const App: React.FC = () => {
  return (
    <Provider store={store}>
      <BrowserRouter>
        <Suspense fallback={<PageLoader />}>
          <Routes>
            {/* Public routes with main layout */}
            <Route
              path={ROUTES.HOME}
              element={
                <MainLayout>
                  <HomePage />
                </MainLayout>
              }
            />
            <Route
              path={ROUTES.PRODUCTS}
              element={
                <MainLayout>
                  <ProductsPage />
                </MainLayout>
              }
            />
            <Route
              path={ROUTES.PRODUCT_DETAIL}
              element={
                <MainLayout>
                  <ProductDetailPage />
                </MainLayout>
              }
            />
            <Route
              path={ROUTES.CART}
              element={
                <MainLayout>
                  <CartPage />
                </MainLayout>
              }
            />
            <Route
              path={ROUTES.CHECKOUT}
              element={
                <ProtectedRoute>
                  <MainLayout>
                    <CheckoutPage />
                  </MainLayout>
                </ProtectedRoute>
              }
            />
            <Route
              path={ROUTES.CHECKOUT_SUCCESS}
              element={
                <ProtectedRoute>
                  <MainLayout>
                    <div className="container mx-auto px-4 py-16 text-center">
                      <h1 className="text-3xl font-bold text-green-600 mb-4">Order Placed Successfully!</h1>
                      <p className="text-gray-600 mb-8">Thank you for your order. We'll process it shortly.</p>
                      <a href={ROUTES.USER_ORDERS} className="btn-primary">View Orders</a>
                    </div>
                  </MainLayout>
                </ProtectedRoute>
              }
            />

            {/* Auth routes (no layout) */}
            <Route path={ROUTES.LOGIN} element={<LoginPage />} />
            <Route path={ROUTES.REGISTER} element={<RegisterPage />} />

            {/* User dashboard routes */}
            <Route
              path={ROUTES.USER_DASHBOARD}
              element={
                <ProtectedRoute>
                  <UserDashboard />
                </ProtectedRoute>
              }
            >
              <Route index element={<UserProfilePage />} />
              <Route path="profile" element={<UserProfilePage />} />
              <Route path="orders" element={<UserOrdersPage />} />
              <Route path="reviews" element={
                <div className="bg-white rounded-lg shadow-md p-6">
                  <h2 className="text-xl font-bold mb-4">My Reviews</h2>
                  <p className="text-gray-600">Review history coming soon...</p>
                </div>
              } />
              <Route path="addresses" element={
                <div className="bg-white rounded-lg shadow-md p-6">
                  <h2 className="text-xl font-bold mb-4">My Addresses</h2>
                  <p className="text-gray-600">Address management coming soon...</p>
                </div>
              } />
              <Route path="settings" element={
                <div className="bg-white rounded-lg shadow-md p-6">
                  <h2 className="text-xl font-bold mb-4">Settings</h2>
                  <p className="text-gray-600">Settings coming soon...</p>
                </div>
              } />
            </Route>

            {/* Admin routes - placeholder */}
            <Route
              path={ROUTES.ADMIN_DASHBOARD}
              element={
                <ProtectedRoute>
                  <div className="min-h-screen bg-gray-50">
                    <Navbar />
                    <div className="container mx-auto px-4 py-8">
                      <h1 className="text-3xl font-bold text-gray-900 mb-6">Admin Dashboard</h1>
                      <div className="bg-white rounded-lg shadow-md p-6">
                        <p className="text-gray-600">Admin dashboard coming soon...</p>
                      </div>
                    </div>
                    <Footer />
                  </div>
                </ProtectedRoute>
              }
            />

            {/* 404 route */}
            <Route
              path="*"
              element={
                <MainLayout>
                  <div className="container mx-auto px-4 py-16 text-center">
                    <h1 className="text-6xl font-bold text-gray-300 mb-4">404</h1>
                    <h2 className="text-2xl font-semibold text-gray-900 mb-4">Page Not Found</h2>
                    <p className="text-gray-600 mb-8">
                      The page you're looking for doesn't exist or has been moved.
                    </p>
                    <a href={ROUTES.HOME} className="btn-primary">Go Home</a>
                  </div>
                </MainLayout>
              }
            />
          </Routes>
        </Suspense>
      </BrowserRouter>
    </Provider>
  );
};

export default App;
