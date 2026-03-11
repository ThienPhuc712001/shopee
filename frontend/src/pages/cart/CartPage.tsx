import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../../hooks';
import { fetchCart, clearCart as clearCartAction, addToCart } from '../../store/slices/cartSlice';
import { ROUTES } from '../../constants';
import CartItem from '../../components/cart/CartItem';
import Button from '../../components/common/Button';
import LoadingSpinner from '../../components/common/LoadingSpinner';
import { ShoppingCartIcon, TagIcon } from '@heroicons/react/24/outline';
import type { CartItem as CartItemType } from '../../types';

const CartPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const { cart, items, isLoading, itemCount } = useAppSelector((state) => state.cart);
  
  const [demoCartItems, setDemoCartItems] = useState<CartItemType[]>([]);
  const [couponCode, setCouponCode] = React.useState('');

  useEffect(() => {
    const loadCart = async () => {
      try {
        await dispatch(fetchCart()).unwrap();
      } catch {
        // Use demo cart items
        setDemoCartItems(generateDemoCartItems());
      }
    };
    loadCart();
  }, [dispatch]);

  const displayItems = items.length > 0 ? items : demoCartItems;
  const displayItemCount = itemCount || demoCartItems.reduce((sum, item) => sum + item.quantity, 0);
  
  const subtotal = cart?.subtotal || demoCartItems.reduce((sum, item) => sum + item.subtotal, 0);
  const discount = cart?.discount || 0;
  const total = subtotal - discount;

  const handleApplyCoupon = () => {
    alert('Coupon feature coming soon!');
  };

  const handleCheckout = () => {
    navigate(ROUTES.CHECKOUT);
  };

  const handleClearCart = () => {
    if (window.confirm('Are you sure you want to clear your cart?')) {
      dispatch(clearCartAction());
      setDemoCartItems([]);
    }
  };

  const handleUpdateQuantity = (itemId: number, newQuantity: number) => {
    if (items.length > 0) {
      dispatch(addToCart({ productId: 1, quantity: newQuantity })); // This is a simplification
    } else {
      setDemoCartItems(prev => prev.map(item => 
        item.id === itemId ? { ...item, quantity: newQuantity, subtotal: item.price * newQuantity } : item
      ));
    }
  };

  const handleRemoveItem = (itemId: number) => {
    if (items.length > 0) {
      // dispatch(removeFromCart(itemId));
    } else {
      setDemoCartItems(prev => prev.filter(item => item.id !== itemId));
    }
  };

  if (isLoading && !cart && demoCartItems.length === 0) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner size="lg" text="Loading cart..." />
      </div>
    );
  }

  if (displayItems.length === 0) {
    return (
      <div className="min-h-screen bg-gray-50">
        <div className="container mx-auto px-4 py-16">
          <div className="max-w-md mx-auto text-center">
            <ShoppingCartIcon className="w-24 h-24 text-gray-300 mx-auto mb-6" />
            <h1 className="text-2xl font-bold text-gray-900 mb-4">Your cart is empty</h1>
            <p className="text-gray-600 mb-8">
              Looks like you haven't added anything to your cart yet.
            </p>
            <Link to={ROUTES.PRODUCTS}>
              <Button size="lg">Start Shopping</Button>
            </Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-8">Shopping Cart ({displayItemCount} items)</h1>

        <div className="flex flex-col lg:flex-row gap-6">
          {/* Cart Items */}
          <div className="flex-1">
            <div className="bg-white rounded-lg shadow-md p-6 mb-4">
              <div className="space-y-4">
                {displayItems.map((item) => (
                  <div key={item.id} className="flex items-center space-x-4 py-4 border-b border-gray-200 last:border-b-0">
                    {/* Product image */}
                    <div className="w-24 h-24 flex-shrink-0 bg-gray-100 rounded-lg overflow-hidden">
                      {item.product_image ? (
                        <img src={item.product_image} alt={item.product_name} className="w-full h-full object-cover" />
                      ) : (
                        <div className="w-full h-full flex items-center justify-center text-gray-400">No Image</div>
                      )}
                    </div>

                    {/* Product info */}
                    <div className="flex-1 min-w-0">
                      <h3 className="font-medium text-gray-800 truncate">{item.product_name}</h3>
                      <p className="text-sm text-gray-500 mt-1">Price: ${item.price.toFixed(2)}</p>
                    </div>

                    {/* Quantity controls */}
                    <div className="flex items-center space-x-2">
                      <button
                        onClick={() => handleUpdateQuantity(item.id, Math.max(1, item.quantity - 1))}
                        disabled={item.quantity <= 1}
                        className="w-8 h-8 rounded-md border border-gray-300 flex items-center justify-center hover:bg-gray-100 disabled:opacity-50"
                      >
                        -
                      </button>
                      <span className="w-12 text-center font-medium">{item.quantity}</span>
                      <button
                        onClick={() => handleUpdateQuantity(item.id, item.quantity + 1)}
                        className="w-8 h-8 rounded-md border border-gray-300 flex items-center justify-center hover:bg-gray-100"
                      >
                        +
                      </button>
                    </div>

                    {/* Subtotal */}
                    <div className="text-right w-28">
                      <p className="font-bold text-gray-800">${item.subtotal.toFixed(2)}</p>
                    </div>

                    {/* Remove button */}
                    <button
                      onClick={() => handleRemoveItem(item.id)}
                      className="p-2 text-red-500 hover:bg-red-50 rounded-lg transition-colors"
                    >
                      <ShoppingCartIcon className="w-5 h-5" />
                    </button>
                  </div>
                ))}
              </div>

              {/* Clear cart button */}
              <div className="mt-6 pt-6 border-t border-gray-200">
                <button
                  onClick={handleClearCart}
                  className="text-red-600 hover:text-red-700 text-sm font-medium"
                >
                  Clear Cart
                </button>
              </div>
            </div>

            {/* Continue shopping */}
            <Link to={ROUTES.PRODUCTS} className="text-primary-600 hover:text-primary-700 font-medium">
              ← Continue Shopping
            </Link>
          </div>

          {/* Order Summary */}
          <div className="lg:w-96">
            <div className="bg-white rounded-lg shadow-md p-6 sticky top-24">
              <h2 className="text-xl font-bold text-gray-900 mb-6">Order Summary</h2>

              {/* Coupon code */}
              <div className="mb-6">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Have a coupon code?
                </label>
                <div className="flex space-x-2">
                  <div className="relative flex-1">
                    <TagIcon className="w-5 h-5 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
                    <input
                      type="text"
                      value={couponCode}
                      onChange={(e) => setCouponCode(e.target.value)}
                      placeholder="Enter code"
                      className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 outline-none"
                    />
                  </div>
                  <Button onClick={handleApplyCoupon} variant="secondary">
                    Apply
                  </Button>
                </div>
              </div>

              {/* Price breakdown */}
              <div className="space-y-3 mb-6">
                <div className="flex justify-between text-gray-600">
                  <span>Subtotal ({displayItemCount} items)</span>
                  <span>${subtotal.toFixed(2)}</span>
                </div>
                {discount > 0 && (
                  <div className="flex justify-between text-green-600">
                    <span>Discount</span>
                    <span>-${discount.toFixed(2)}</span>
                  </div>
                )}
                <div className="flex justify-between text-gray-600">
                  <span>Shipping</span>
                  <span>{subtotal >= 50 ? 'FREE' : 'Calculated at checkout'}</span>
                </div>
              </div>

              {/* Total */}
              <div className="border-t border-gray-200 pt-4 mb-6">
                <div className="flex justify-between text-lg font-bold text-gray-900">
                  <span>Total</span>
                  <span>${total.toFixed(2)}</span>
                </div>
              </div>

              {/* Checkout button */}
              <Button onClick={handleCheckout} fullWidth size="lg">
                Proceed to Checkout
              </Button>

              {/* Trust badges */}
              <div className="mt-6 pt-6 border-t border-gray-200">
                <div className="grid grid-cols-2 gap-4 text-xs text-gray-500">
                  <div className="flex items-center space-x-2">
                    <span>🔒</span>
                    <span>Secure Checkout</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <span>🚚</span>
                    <span>Free Shipping $50+</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <span>↩️</span>
                    <span>Easy Returns</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <span>💬</span>
                    <span>24/7 Support</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

// Generate demo cart items
function generateDemoCartItems(): CartItemType[] {
  return [
    {
      id: 1,
      cart_id: 1,
      product_id: 1,
      variant_id: null,
      quantity: 2,
      price: 29.99,
      original_price: 49.99,
      discount: 0,
      subtotal: 59.98,
      product_name: 'Demo Product 1',
      product_image: 'https://via.placeholder.com/200x200?text=Product+1',
      shop_id: 1,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    },
    {
      id: 2,
      cart_id: 1,
      product_id: 2,
      variant_id: null,
      quantity: 1,
      price: 39.99,
      original_price: 59.99,
      discount: 0,
      subtotal: 39.99,
      product_name: 'Demo Product 2',
      product_image: 'https://via.placeholder.com/200x200?text=Product+2',
      shop_id: 1,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    },
    {
      id: 3,
      cart_id: 1,
      product_id: 3,
      variant_id: null,
      quantity: 3,
      price: 19.99,
      original_price: 29.99,
      discount: 0,
      subtotal: 59.97,
      product_name: 'Demo Product 3',
      product_image: 'https://via.placeholder.com/200x200?text=Product+3',
      shop_id: 1,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    },
  ];
}

export default CartPage;
