import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useSelector, useDispatch } from 'react-redux'
import { AppDispatch, RootState } from '@store/store'
import { clearCart } from '@store/slices/cartSlice'
import { CreditCard, Truck, CheckCircle, ArrowLeft, Lock } from 'lucide-react'
import orderService from '@services/orderService'

function CheckoutPage() {
  const navigate = useNavigate()
  const dispatch = useDispatch<AppDispatch>()
  const { cart } = useSelector((state: RootState) => state.cart)
  const { user } = useSelector((state: RootState) => state.auth)
  
  const [step, setStep] = useState(1)
  const [loading, setLoading] = useState(false)
  const [orderId, setOrderId] = useState<number | null>(null)
  
  const [formData, setFormData] = useState({
    // Shipping Address
    fullName: user?.first_name ? `${user.first_name} ${user.last_name}` : '',
    phone: user?.phone || '',
    address: '',
    ward: '',
    district: '',
    city: 'Ho Chi Minh City',
    
    // Payment
    paymentMethod: 'cod',
    
    // Order notes
    notes: ''
  })

  const subtotal = cart?.subtotal || 0
  const shippingFee = subtotal > 50 ? 0 : 5.99
  const discount = 0 // Will calculate based on coupon
  const total = subtotal + shippingFee - discount

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value })
  }

  const handlePlaceOrder = async () => {
    if (!cart || cart.items.length === 0) {
      alert('Your cart is empty!')
      return
    }

    if (!formData.fullName || !formData.phone || !formData.address) {
      alert('Please fill in all required fields!')
      setStep(1)
      return
    }

    setLoading(true)
    try {
      const orderData = {
        items: cart.items.map(item => ({
          product_id: item.product.id,
          quantity: item.quantity
        })),
        shipping_address: {
          name: formData.fullName,
          phone: formData.phone,
          address: `${formData.address}${formData.ward ? ', ' + formData.ward : ''}${formData.district ? ', ' + formData.district : ''}`,
          ward: formData.ward,
          district: formData.district,
          city: formData.city
        },
        payment_method: formData.paymentMethod,
        coupon_code: undefined // Add coupon logic here
      }

      const order = await orderService.checkout(orderData)
      setOrderId(order.id)
      dispatch(clearCart())
      setStep(4) // Success step
    } catch (error: any) {
      console.error('Checkout error:', error)
      alert(error.response?.data?.message || 'Failed to place order. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  const renderStepIndicator = () => (
    <div className="mb-8">
      <div className="flex items-center justify-center">
        {[
          { num: 1, title: 'Shipping', icon: Truck },
          { num: 2, title: 'Payment', icon: CreditCard },
          { num: 3, title: 'Review', icon: Lock },
          { num: 4, title: 'Complete', icon: CheckCircle }
        ].map((s, index) => (
          <div key={s.num} className="flex items-center">
            <div className={`flex items-center justify-center w-12 h-12 rounded-full border-2 ${
              step >= s.num ? 'bg-primary-500 border-primary-500 text-white' : 'bg-gray-200 border-gray-300 text-gray-500'
            }`}>
              <s.icon className="w-6 h-6" />
            </div>
            <span className={`ml-2 font-medium ${step >= s.num ? 'text-primary-500' : 'text-gray-500'}`}>
              {s.title}
            </span>
            {index < 3 && (
              <div className={`w-16 h-1 mx-4 ${step > s.num ? 'bg-primary-500' : 'bg-gray-300'}`} />
            )}
          </div>
        ))}
      </div>
    </div>
  )

  const renderShippingStep = () => (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-gray-900 mb-6">Shipping Information</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="md:col-span-2">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Full Name *
            </label>
            <input
              type="text"
              name="fullName"
              value={formData.fullName}
              onChange={handleInputChange}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
              placeholder="John Doe"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Phone Number *
            </label>
            <input
              type="tel"
              name="phone"
              value={formData.phone}
              onChange={handleInputChange}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
              placeholder="0123456789"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              City *
            </label>
            <select
              name="city"
              value={formData.city}
              onChange={handleInputChange}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
            >
              <option value="Ho Chi Minh City">Ho Chi Minh City</option>
              <option value="Ha Noi">Ha Noi</option>
              <option value="Da Nang">Da Nang</option>
              <option value="Can Tho">Can Tho</option>
              <option value="Hai Phong">Hai Phong</option>
            </select>
          </div>

          <div className="md:col-span-2">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Street Address *
            </label>
            <input
              type="text"
              name="address"
              value={formData.address}
              onChange={handleInputChange}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
              placeholder="123 Main Street, Apartment, Suite, etc."
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Ward
            </label>
            <input
              type="text"
              name="ward"
              value={formData.ward}
              onChange={handleInputChange}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
              placeholder="Ward 1"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              District
            </label>
            <input
              type="text"
              name="district"
              value={formData.district}
              onChange={handleInputChange}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
              placeholder="District 1"
            />
          </div>
        </div>
      </div>

      <div className="flex justify-end">
        <button
          onClick={() => setStep(2)}
          className="bg-primary-500 hover:bg-primary-600 text-white font-semibold px-8 py-4 rounded-lg transition"
        >
          Continue to Payment
        </button>
      </div>
    </div>
  )

  const renderPaymentStep = () => (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-gray-900 mb-6">Payment Method</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {/* Cash on Delivery */}
          <label className={`relative border-2 rounded-lg p-6 cursor-pointer transition ${
            formData.paymentMethod === 'cod' 
              ? 'border-primary-500 bg-primary-50' 
              : 'border-gray-200 hover:border-gray-300'
          }`}>
            <input
              type="radio"
              name="paymentMethod"
              value="cod"
              checked={formData.paymentMethod === 'cod'}
              onChange={handleInputChange}
              className="sr-only"
            />
            <div className="flex items-start gap-4">
              <div className={`w-6 h-6 rounded-full border-2 flex items-center justify-center ${
                formData.paymentMethod === 'cod' ? 'border-primary-500 bg-primary-500' : 'border-gray-300'
              }`}>
                {formData.paymentMethod === 'cod' && <CheckCircle className="w-4 h-4 text-white" />}
              </div>
              <div>
                <h3 className="font-semibold text-gray-900">Cash on Delivery (COD)</h3>
                <p className="text-sm text-gray-600 mt-1">Pay when you receive your order</p>
              </div>
            </div>
          </label>

          {/* Bank Transfer */}
          <label className={`relative border-2 rounded-lg p-6 cursor-pointer transition ${
            formData.paymentMethod === 'bank_transfer' 
              ? 'border-primary-500 bg-primary-50' 
              : 'border-gray-200 hover:border-gray-300'
          }`}>
            <input
              type="radio"
              name="paymentMethod"
              value="bank_transfer"
              checked={formData.paymentMethod === 'bank_transfer'}
              onChange={handleInputChange}
              className="sr-only"
            />
            <div className="flex items-start gap-4">
              <div className={`w-6 h-6 rounded-full border-2 flex items-center justify-center ${
                formData.paymentMethod === 'bank_transfer' ? 'border-primary-500 bg-primary-500' : 'border-gray-300'
              }`}>
                {formData.paymentMethod === 'bank_transfer' && <CheckCircle className="w-4 h-4 text-white" />}
              </div>
              <div>
                <h3 className="font-semibold text-gray-900">Bank Transfer</h3>
                <p className="text-sm text-gray-600 mt-1">Transfer to our bank account</p>
              </div>
            </div>
          </label>

          {/* Credit Card */}
          <label className={`relative border-2 rounded-lg p-6 cursor-pointer transition ${
            formData.paymentMethod === 'credit_card' 
              ? 'border-primary-500 bg-primary-50' 
              : 'border-gray-200 hover:border-gray-300'
          }`}>
            <input
              type="radio"
              name="paymentMethod"
              value="credit_card"
              checked={formData.paymentMethod === 'credit_card'}
              onChange={handleInputChange}
              className="sr-only"
            />
            <div className="flex items-start gap-4">
              <div className={`w-6 h-6 rounded-full border-2 flex items-center justify-center ${
                formData.paymentMethod === 'credit_card' ? 'border-primary-500 bg-primary-500' : 'border-gray-300'
              }`}>
                {formData.paymentMethod === 'credit_card' && <CheckCircle className="w-4 h-4 text-white" />}
              </div>
              <div>
                <h3 className="font-semibold text-gray-900">Credit/Debit Card</h3>
                <p className="text-sm text-gray-600 mt-1">Pay securely with card</p>
              </div>
            </div>
          </label>

          {/* PayPal */}
          <label className={`relative border-2 rounded-lg p-6 cursor-pointer transition ${
            formData.paymentMethod === 'paypal' 
              ? 'border-primary-500 bg-primary-50' 
              : 'border-gray-200 hover:border-gray-300'
          }`}>
            <input
              type="radio"
              name="paymentMethod"
              value="paypal"
              checked={formData.paymentMethod === 'paypal'}
              onChange={handleInputChange}
              className="sr-only"
            />
            <div className="flex items-start gap-4">
              <div className={`w-6 h-6 rounded-full border-2 flex items-center justify-center ${
                formData.paymentMethod === 'paypal' ? 'border-primary-500 bg-primary-500' : 'border-gray-300'
              }`}>
                {formData.paymentMethod === 'paypal' && <CheckCircle className="w-4 h-4 text-white" />}
              </div>
              <div>
                <h3 className="font-semibold text-gray-900">PayPal</h3>
                <p className="text-sm text-gray-600 mt-1">Pay with your PayPal account</p>
              </div>
            </div>
          </label>
        </div>
      </div>

      <div className="flex justify-between">
        <button
          onClick={() => setStep(1)}
          className="flex items-center gap-2 text-gray-600 hover:text-gray-900 font-medium px-6 py-4"
        >
          <ArrowLeft className="w-5 h-5" />
          Back to Shipping
        </button>
        <button
          onClick={() => setStep(3)}
          className="bg-primary-500 hover:bg-primary-600 text-white font-semibold px-8 py-4 rounded-lg transition"
        >
          Review Order
        </button>
      </div>
    </div>
  )

  const renderReviewStep = () => (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-gray-900 mb-6">Review Your Order</h2>
        
        {/* Shipping Address Review */}
        <div className="bg-gray-50 rounded-lg p-6 mb-6">
          <h3 className="font-semibold text-gray-900 mb-3">Shipping Address</h3>
          <div className="text-gray-600">
            <p className="font-medium">{formData.fullName}</p>
            <p>{formData.address}</p>
            {formData.ward && <p>{formData.ward}</p>}
            {formData.district && <p>{formData.district}</p>}
            <p>{formData.city}</p>
            <p className="mt-2">Phone: {formData.phone}</p>
          </div>
        </div>

        {/* Payment Method Review */}
        <div className="bg-gray-50 rounded-lg p-6 mb-6">
          <h3 className="font-semibold text-gray-900 mb-3">Payment Method</h3>
          <p className="text-gray-600">
            {formData.paymentMethod === 'cod' && 'Cash on Delivery (COD)'}
            {formData.paymentMethod === 'bank_transfer' && 'Bank Transfer'}
            {formData.paymentMethod === 'credit_card' && 'Credit/Debit Card'}
            {formData.paymentMethod === 'paypal' && 'PayPal'}
          </p>
        </div>

        {/* Order Items */}
        <div className="bg-gray-50 rounded-lg p-6 mb-6">
          <h3 className="font-semibold text-gray-900 mb-4">Order Items</h3>
          <div className="space-y-4">
            {cart?.items.map((item) => (
              <div key={item.id} className="flex items-center gap-4 pb-4 border-b last:border-0">
                <img 
                  src={item.product.images?.[0]?.url || 'https://via.placeholder.com/80'} 
                  alt={item.product.name}
                  className="w-20 h-20 object-cover rounded-lg"
                />
                <div className="flex-1">
                  <p className="font-medium text-gray-900">{item.product.name}</p>
                  <p className="text-sm text-gray-600">Quantity: {item.quantity}</p>
                  <p className="text-sm text-gray-600">Price: ${item.product.price.toFixed(2)}</p>
                </div>
                <p className="font-semibold text-gray-900">${item.subtotal.toFixed(2)}</p>
              </div>
            ))}
          </div>
        </div>

        {/* Order Summary */}
        <div className="bg-white border rounded-lg p-6">
          <h3 className="font-semibold text-gray-900 mb-4">Order Summary</h3>
          <div className="space-y-3">
            <div className="flex justify-between text-gray-600">
              <span>Subtotal</span>
              <span>${subtotal.toFixed(2)}</span>
            </div>
            <div className="flex justify-between text-gray-600">
              <span>Shipping Fee</span>
              <span>{shippingFee === 0 ? 'FREE' : `$${shippingFee.toFixed(2)}`}</span>
            </div>
            {shippingFee === 0 && (
              <p className="text-sm text-green-600">✓ You qualified for FREE shipping!</p>
            )}
            <div className="flex justify-between text-gray-600">
              <span>Discount</span>
              <span className="text-green-600">-${discount.toFixed(2)}</span>
            </div>
            <div className="border-t pt-3 flex justify-between text-lg font-bold">
              <span>Total</span>
              <span className="text-primary-500">${total.toFixed(2)}</span>
            </div>
          </div>
        </div>
      </div>

      <div className="flex justify-between">
        <button
          onClick={() => setStep(2)}
          className="flex items-center gap-2 text-gray-600 hover:text-gray-900 font-medium px-6 py-4"
        >
          <ArrowLeft className="w-5 h-5" />
          Back to Payment
        </button>
        <button
          onClick={handlePlaceOrder}
          disabled={loading}
          className="bg-primary-500 hover:bg-primary-600 disabled:bg-gray-400 text-white font-semibold px-8 py-4 rounded-lg transition flex items-center gap-2"
        >
          {loading ? (
            <>
              <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
              Processing...
            </>
          ) : (
            <>
              <Lock className="w-5 h-5" />
              Place Order - ${total.toFixed(2)}
            </>
          )}
        </button>
      </div>
    </div>
  )

  const renderSuccessStep = () => (
    <div className="text-center py-12">
      <div className="w-24 h-24 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6">
        <CheckCircle className="w-16 h-16 text-green-500" />
      </div>
      <h2 className="text-3xl font-bold text-gray-900 mb-4">Order Placed Successfully!</h2>
      <p className="text-gray-600 mb-2">Thank you for your purchase.</p>
      <p className="text-gray-600 mb-8">
        Order ID: <span className="font-semibold">#{orderId}</span>
      </p>
      
      <div className="bg-gray-50 rounded-lg p-6 max-w-md mx-auto mb-8">
        <h3 className="font-semibold text-gray-900 mb-4">What's Next?</h3>
        <div className="space-y-3 text-left text-gray-600">
          <div className="flex items-start gap-3">
            <span className="flex-shrink-0 w-6 h-6 bg-primary-100 text-primary-500 rounded-full flex items-center justify-center text-sm font-semibold">1</span>
            <p>Order confirmation email will be sent to your email address</p>
          </div>
          <div className="flex items-start gap-3">
            <span className="flex-shrink-0 w-6 h-6 bg-primary-100 text-primary-500 rounded-full flex items-center justify-center text-sm font-semibold">2</span>
            <p>We'll process your order within 1-2 business days</p>
          </div>
          <div className="flex items-start gap-3">
            <span className="flex-shrink-0 w-6 h-6 bg-primary-100 text-primary-500 rounded-full flex items-center justify-center text-sm font-semibold">3</span>
            <p>You'll receive tracking information via email</p>
          </div>
          <div className="flex items-start gap-3">
            <span className="flex-shrink-0 w-6 h-6 bg-primary-100 text-primary-500 rounded-full flex items-center justify-center text-sm font-semibold">4</span>
            <p>Estimated delivery: 5-7 business days</p>
          </div>
        </div>
      </div>

      <div className="flex flex-col sm:flex-row gap-4 justify-center">
        <Link
          to="/account/orders"
          className="bg-secondary-800 hover:bg-secondary-700 text-white font-semibold px-8 py-4 rounded-lg transition"
        >
          View My Orders
        </Link>
        <Link
          to="/products"
          className="bg-primary-500 hover:bg-primary-600 text-white font-semibold px-8 py-4 rounded-lg transition"
        >
          Continue Shopping
        </Link>
      </div>
    </div>
  )

  if (!user) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4">
        <div className="text-center">
          <Lock className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h2 className="text-2xl font-bold text-gray-900 mb-4">Login Required</h2>
          <p className="text-gray-600 mb-6">Please login to proceed with checkout</p>
          <Link
            to="/login"
            className="inline-block bg-primary-500 hover:bg-primary-600 text-white font-semibold px-8 py-4 rounded-lg transition"
          >
            Login Now
          </Link>
        </div>
      </div>
    )
  }

  if (!cart || cart.items.length === 0) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4">
        <div className="text-center">
          <Truck className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h2 className="text-2xl font-bold text-gray-900 mb-4">Your Cart is Empty</h2>
          <p className="text-gray-600 mb-6">Add some products before checkout</p>
          <Link
            to="/products"
            className="inline-block bg-primary-500 hover:bg-primary-600 text-white font-semibold px-8 py-4 rounded-lg transition"
          >
            Browse Products
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="container mx-auto px-4 max-w-5xl">
        {step < 4 && renderStepIndicator()}
        
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Main Content */}
          <div className="lg:col-span-2">
            <div className="bg-white rounded-lg shadow-lg p-8">
              {step === 1 && renderShippingStep()}
              {step === 2 && renderPaymentStep()}
              {step === 3 && renderReviewStep()}
              {step === 4 && renderSuccessStep()}
            </div>
          </div>

          {/* Order Summary Sidebar */}
          {step < 4 && (
            <div className="lg:col-span-1">
              <div className="bg-white rounded-lg shadow-lg p-6 sticky top-24">
                <h3 className="text-lg font-bold text-gray-900 mb-4">Order Summary</h3>
                <div className="space-y-3 mb-6">
                  {cart?.items.slice(0, 3).map((item) => (
                    <div key={item.id} className="flex gap-3">
                      <div className="relative">
                        <img 
                          src={item.product.images?.[0]?.url || 'https://via.placeholder.com/60'} 
                          alt={item.product.name}
                          className="w-16 h-16 object-cover rounded-lg"
                        />
                        <span className="absolute -top-2 -right-2 bg-gray-500 text-white text-xs w-5 h-5 rounded-full flex items-center justify-center">
                          {item.quantity}
                        </span>
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium text-gray-900 truncate">{item.product.name}</p>
                        <p className="text-sm text-gray-600">${item.product.price.toFixed(2)}</p>
                      </div>
                    </div>
                  ))}
                  {cart.items.length > 3 && (
                    <p className="text-sm text-gray-600">+ {cart.items.length - 3} more items</p>
                  )}
                </div>
                <div className="border-t pt-4 space-y-2">
                  <div className="flex justify-between text-gray-600">
                    <span>Subtotal</span>
                    <span>${subtotal.toFixed(2)}</span>
                  </div>
                  <div className="flex justify-between text-gray-600">
                    <span>Shipping</span>
                    <span>{shippingFee === 0 ? 'FREE' : `$${shippingFee.toFixed(2)}`}</span>
                  </div>
                  <div className="border-t pt-2 flex justify-between text-lg font-bold">
                    <span>Total</span>
                    <span className="text-primary-500">${total.toFixed(2)}</span>
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default CheckoutPage
