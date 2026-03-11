import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../../hooks';
import { fetchCart } from '../../store/slices/cartSlice';
import { checkout } from '../../store/slices/ordersSlice';
import { ROUTES } from '../../constants';
import Button from '../../components/common/Button';
import Input from '../../components/common/Input';
import LoadingSpinner from '../../components/common/LoadingSpinner';
import { CreditCardIcon, TruckIcon, CheckCircleIcon } from '@heroicons/react/24/outline';

const CheckoutPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const { cart, isLoading: isCartLoading } = useAppSelector((state) => state.cart);
  const { isCheckingOut } = useAppSelector((state) => state.orders);
  
  const [step, setStep] = useState(1);
  const [shippingInfo, setShippingInfo] = useState({
    name: '',
    phone: '',
    street: '',
    ward: '',
    district: '',
    city: '',
    country: 'Vietnam',
  });
  const [paymentMethod, setPaymentMethod] = useState('credit_card');
  const [buyerNote, setBuyerNote] = useState('');

  useEffect(() => {
    dispatch(fetchCart());
  }, [dispatch]);

  const handleShippingSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setStep(2);
  };

  const handleCheckout = async () => {
    try {
      const result = await dispatch(checkout({
        shipping_address_id: 0, // Will be created or selected
        payment_method: paymentMethod,
        buyer_note: buyerNote,
      })).unwrap();
      
      // Navigate to success page
      navigate(`${ROUTES.CHECKOUT_SUCCESS}?order_id=${result.id}`);
    } catch (error) {
      console.error('Checkout failed:', error);
    }
  };

  if (isCartLoading || !cart) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner size="lg" text="Loading checkout..." />
      </div>
    );
  }

  if (!cart.items || cart.items.length === 0) {
    navigate(ROUTES.CART);
    return null;
  }

  const subtotal = cart.subtotal || 0;
  const shippingFee = subtotal >= 50 ? 0 : 5.99;
  const total = subtotal + shippingFee - (cart.discount || 0);

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-8">Checkout</h1>

        {/* Progress steps */}
        <div className="flex items-center justify-center mb-8">
          <div className="flex items-center">
            <div className={`flex items-center justify-center w-10 h-10 rounded-full ${
              step >= 1 ? 'bg-primary-600 text-white' : 'bg-gray-200 text-gray-600'
            }`}>
              1
            </div>
            <span className={`ml-2 font-medium ${step >= 1 ? 'text-gray-900' : 'text-gray-500'}`}>
              Shipping
            </span>
          </div>
          <div className={`w-16 h-1 mx-4 ${step >= 2 ? 'bg-primary-600' : 'bg-gray-200'}`} />
          <div className="flex items-center">
            <div className={`flex items-center justify-center w-10 h-10 rounded-full ${
              step >= 2 ? 'bg-primary-600 text-white' : 'bg-gray-200 text-gray-600'
            }`}>
              2
            </div>
            <span className={`ml-2 font-medium ${step >= 2 ? 'text-gray-900' : 'text-gray-500'}`}>
              Payment
            </span>
          </div>
        </div>

        <div className="flex flex-col lg:flex-row gap-6">
          {/* Main content */}
          <div className="flex-1">
            {step === 1 ? (
              /* Shipping Information */
              <div className="bg-white rounded-lg shadow-md p-6">
                <h2 className="text-xl font-bold text-gray-900 mb-6 flex items-center">
                  <TruckIcon className="w-6 h-6 mr-2" />
                  Shipping Information
                </h2>

                <form onSubmit={handleShippingSubmit} className="space-y-4">
                  <div className="grid grid-cols-2 gap-4">
                    <Input
                      label="Full Name"
                      value={shippingInfo.name}
                      onChange={(e) => setShippingInfo({ ...shippingInfo, name: e.target.value })}
                      required
                    />
                    <Input
                      label="Phone Number"
                      value={shippingInfo.phone}
                      onChange={(e) => setShippingInfo({ ...shippingInfo, phone: e.target.value })}
                      required
                    />
                  </div>

                  <Input
                    label="Street Address"
                    value={shippingInfo.street}
                    onChange={(e) => setShippingInfo({ ...shippingInfo, street: e.target.value })}
                    required
                  />

                  <div className="grid grid-cols-2 gap-4">
                    <Input
                      label="Ward"
                      value={shippingInfo.ward}
                      onChange={(e) => setShippingInfo({ ...shippingInfo, ward: e.target.value })}
                    />
                    <Input
                      label="District"
                      value={shippingInfo.district}
                      onChange={(e) => setShippingInfo({ ...shippingInfo, district: e.target.value })}
                      required
                    />
                  </div>

                  <div className="grid grid-cols-2 gap-4">
                    <Input
                      label="City"
                      value={shippingInfo.city}
                      onChange={(e) => setShippingInfo({ ...shippingInfo, city: e.target.value })}
                      required
                    />
                    <Input
                      label="Country"
                      value={shippingInfo.country}
                      onChange={(e) => setShippingInfo({ ...shippingInfo, country: e.target.value })}
                      required
                    />
                  </div>

                  <div className="pt-4">
                    <Button type="submit" fullWidth size="lg">
                      Continue to Payment
                    </Button>
                  </div>
                </form>
              </div>
            ) : (
              /* Payment Information */
              <div className="bg-white rounded-lg shadow-md p-6">
                <h2 className="text-xl font-bold text-gray-900 mb-6 flex items-center">
                  <CreditCardIcon className="w-6 h-6 mr-2" />
                  Payment Method
                </h2>

                <div className="space-y-4 mb-6">
                  {[
                    { value: 'credit_card', label: 'Credit/Debit Card', icon: '💳' },
                    { value: 'bank_transfer', label: 'Bank Transfer', icon: '🏦' },
                    { value: 'e_wallet', label: 'E-Wallet', icon: '📱' },
                    { value: 'cod', label: 'Cash on Delivery', icon: '💵' },
                  ].map((method) => (
                    <label
                      key={method.value}
                      className={`flex items-center p-4 border rounded-lg cursor-pointer transition-colors ${
                        paymentMethod === method.value
                          ? 'border-primary-600 bg-primary-50'
                          : 'border-gray-200 hover:border-primary-600'
                      }`}
                    >
                      <input
                        type="radio"
                        name="payment_method"
                        value={method.value}
                        checked={paymentMethod === method.value}
                        onChange={(e) => setPaymentMethod(e.target.value)}
                        className="w-4 h-4 text-primary-600"
                      />
                      <span className="ml-3 text-xl mr-2">{method.icon}</span>
                      <span className="font-medium">{method.label}</span>
                    </label>
                  ))}
                </div>

                <div className="mb-6">
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Order Notes (Optional)
                  </label>
                  <textarea
                    value={buyerNote}
                    onChange={(e) => setBuyerNote(e.target.value)}
                    placeholder="Any special instructions for your order?"
                    rows={3}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 outline-none resize-none"
                  />
                </div>

                <div className="flex space-x-4">
                  <Button onClick={() => setStep(1)} variant="secondary" fullWidth>
                    Back
                  </Button>
                  <Button
                    onClick={handleCheckout}
                    isLoading={isCheckingOut}
                    fullWidth
                    size="lg"
                  >
                    Place Order
                  </Button>
                </div>
              </div>
            )}
          </div>

          {/* Order Summary */}
          <div className="lg:w-96">
            <div className="bg-white rounded-lg shadow-md p-6 sticky top-24">
              <h2 className="text-xl font-bold text-gray-900 mb-6">Order Summary</h2>

              {/* Items */}
              <div className="space-y-4 mb-6 max-h-64 overflow-y-auto">
                {cart.items.map((item) => (
                  <div key={item.id} className="flex items-center space-x-3">
                    <div className="w-16 h-16 bg-gray-100 rounded-lg overflow-hidden flex-shrink-0">
                      {item.product_image ? (
                        <img src={item.product_image} alt={item.product_name} className="w-full h-full object-cover" />
                      ) : (
                        <div className="w-full h-full flex items-center justify-center text-gray-400 text-xs">
                          No Image
                        </div>
                      )}
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-gray-900 truncate">{item.product_name}</p>
                      <p className="text-sm text-gray-500">Qty: {item.quantity}</p>
                    </div>
                    <p className="text-sm font-medium text-gray-900">${item.subtotal.toFixed(2)}</p>
                  </div>
                ))}
              </div>

              {/* Price breakdown */}
              <div className="space-y-3 mb-6 pt-6 border-t border-gray-200">
                <div className="flex justify-between text-gray-600">
                  <span>Subtotal</span>
                  <span>${subtotal.toFixed(2)}</span>
                </div>
                <div className="flex justify-between text-gray-600">
                  <span>Shipping</span>
                  <span>{shippingFee === 0 ? 'FREE' : `$${shippingFee.toFixed(2)}`}</span>
                </div>
                {cart.discount && cart.discount > 0 && (
                  <div className="flex justify-between text-green-600">
                    <span>Discount</span>
                    <span>-${cart.discount.toFixed(2)}</span>
                  </div>
                )}
              </div>

              {/* Total */}
              <div className="border-t border-gray-200 pt-4">
                <div className="flex justify-between text-lg font-bold text-gray-900">
                  <span>Total</span>
                  <span>${total.toFixed(2)}</span>
                </div>
              </div>

              {/* Trust badge */}
              <div className="mt-6 flex items-center justify-center text-green-600">
                <CheckCircleIcon className="w-5 h-5 mr-2" />
                <span className="text-sm font-medium">Secure Checkout</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CheckoutPage;
