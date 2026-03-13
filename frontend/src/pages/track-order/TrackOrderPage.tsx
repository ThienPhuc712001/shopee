import { useState } from 'react'
import { Search, Package, Truck, CheckCircle, Clock, AlertCircle } from 'lucide-react'

function TrackOrderPage() {
  const [orderNumber, setOrderNumber] = useState('')
  const [email, setEmail] = useState('')
  const [tracking, setTracking] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleTrack = (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    
    // Mock tracking data for demo
    setTimeout(() => {
      if (orderNumber && email) {
        setTracking({
          orderNumber: orderNumber,
          status: 'in_transit',
          estimatedDelivery: 'March 18-20, 2026',
          timeline: [
            { status: 'Order Placed', date: 'March 10, 2026', completed: true, icon: Package },
            { status: 'Processing', date: 'March 11, 2026', completed: true, icon: Clock },
            { status: 'Shipped', date: 'March 12, 2026', completed: true, icon: Truck },
            { status: 'In Transit', date: 'March 13, 2026', completed: true, icon: Truck },
            { status: 'Out for Delivery', date: 'Expected March 18', completed: false, icon: Truck },
            { status: 'Delivered', date: 'Expected March 20', completed: false, icon: CheckCircle },
          ],
          shippingAddress: {
            name: 'John Doe',
            address: '123 Main Street, Apt 4B',
            city: 'Ho Chi Minh City',
            postalCode: '700000'
          },
          items: [
            { name: 'Premium Wireless Headphones', quantity: 1, price: 199.99 },
            { name: 'Smart Watch Pro', quantity: 1, price: 349.99 }
          ]
        })
      } else {
        setError('Please enter both order number and email')
      }
      setLoading(false)
    }, 1000)
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'delivered': return 'text-green-600 bg-green-100'
      case 'in_transit': return 'text-blue-600 bg-blue-100'
      case 'processing': return 'text-yellow-600 bg-yellow-100'
      case 'cancelled': return 'text-red-600 bg-red-100'
      default: return 'text-gray-600 bg-gray-100'
    }
  }

  const getStatusText = (status: string) => {
    const statusMap: Record<string, string> = {
      'delivered': 'Delivered',
      'in_transit': 'In Transit',
      'processing': 'Processing',
      'shipped': 'Shipped',
      'cancelled': 'Cancelled'
    }
    return statusMap[status] || 'Unknown'
  }

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">Track Your Order</h1>
          <p className="text-gray-600 text-lg">Enter your order number and email to track your shipment</p>
        </div>

        {/* Tracking Form */}
        <div className="max-w-2xl mx-auto mb-12">
          <div className="bg-white rounded-lg shadow-lg p-8">
            <form onSubmit={handleTrack} className="space-y-6">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Order Number
                </label>
                <input
                  type="text"
                  value={orderNumber}
                  onChange={(e) => setOrderNumber(e.target.value)}
                  placeholder="e.g., ORD-123456"
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Email Address
                </label>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="your@email.com"
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                />
              </div>
              {error && (
                <div className="bg-red-50 text-red-700 px-4 py-3 rounded-lg flex items-center gap-2">
                  <AlertCircle className="w-5 h-5" />
                  {error}
                </div>
              )}
              <button
                type="submit"
                disabled={loading}
                className="w-full bg-primary-500 hover:bg-primary-600 disabled:bg-gray-400 text-white font-semibold py-4 rounded-lg transition flex items-center justify-center gap-2"
              >
                {loading ? (
                  <>
                    <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
                    Tracking...
                  </>
                ) : (
                  <>
                    <Search className="w-5 h-5" />
                    Track Order
                  </>
                )}
              </button>
            </form>
          </div>
        </div>

        {/* Tracking Results */}
        {tracking && (
          <div className="max-w-4xl mx-auto space-y-8">
            {/* Order Status Card */}
            <div className="bg-white rounded-lg shadow-lg p-8">
              <div className="flex items-center justify-between mb-6">
                <div>
                  <h2 className="text-2xl font-bold text-gray-900">Order #{tracking.orderNumber}</h2>
                  <p className="text-gray-600">Estimated Delivery: {tracking.estimatedDelivery}</p>
                </div>
                <span className={`px-4 py-2 rounded-full font-semibold ${getStatusColor(tracking.status)}`}>
                  {getStatusText(tracking.status)}
                </span>
              </div>

              {/* Progress Bar */}
              <div className="relative">
                <div className="absolute top-4 left-0 right-0 h-1 bg-gray-200 rounded">
                  <div 
                    className="h-full bg-primary-500 rounded transition-all"
                    style={{ width: `${(tracking.timeline.filter((t: any) => t.completed).length / tracking.timeline.length) * 100}%` }}
                  ></div>
                </div>
                <div className="relative flex justify-between">
                  {tracking.timeline.map((step: any, index: number) => {
                    const Icon = step.icon
                    return (
                      <div key={index} className="text-center">
                        <div className={`w-8 h-8 rounded-full flex items-center justify-center mx-auto mb-2 ${
                          step.completed ? 'bg-primary-500 text-white' : 'bg-gray-200 text-gray-400'
                        }`}>
                          <Icon className="w-4 h-4" />
                        </div>
                        <p className="text-xs font-medium text-gray-700">{step.status}</p>
                        <p className="text-xs text-gray-500">{step.date}</p>
                      </div>
                    )
                  })}
                </div>
              </div>
            </div>

            {/* Shipping Address */}
            <div className="bg-white rounded-lg shadow-lg p-8">
              <h3 className="text-xl font-bold text-gray-900 mb-4">Shipping Address</h3>
              <div className="text-gray-600">
                <p className="font-semibold">{tracking.shippingAddress.name}</p>
                <p>{tracking.shippingAddress.address}</p>
                <p>{tracking.shippingAddress.city} {tracking.shippingAddress.postalCode}</p>
              </div>
            </div>

            {/* Order Items */}
            <div className="bg-white rounded-lg shadow-lg p-8">
              <h3 className="text-xl font-bold text-gray-900 mb-4">Order Items</h3>
              <div className="space-y-4">
                {tracking.items.map((item: any, index: number) => (
                  <div key={index} className="flex items-center justify-between py-4 border-b last:border-0">
                    <div className="flex items-center gap-4">
                      <div className="w-16 h-16 bg-gray-100 rounded-lg"></div>
                      <div>
                        <p className="font-semibold text-gray-900">{item.name}</p>
                        <p className="text-sm text-gray-600">Qty: {item.quantity}</p>
                      </div>
                    </div>
                    <p className="font-semibold text-gray-900">${item.price}</p>
                  </div>
                ))}
              </div>
              <div className="mt-6 pt-6 border-t flex justify-between items-center">
                <span className="text-lg font-semibold text-gray-900">Total</span>
                <span className="text-2xl font-bold text-primary-500">
                  ${tracking.items.reduce((sum: number, item: any) => sum + item.price * item.quantity, 0).toFixed(2)}
                </span>
              </div>
            </div>
          </div>
        )}

        {/* Help Section */}
        <div className="max-w-2xl mx-auto mt-12 text-center">
          <div className="bg-white rounded-lg shadow-lg p-8">
            <h3 className="text-xl font-bold text-gray-900 mb-4">Need Help?</h3>
            <p className="text-gray-600 mb-6">
              Can't find your order or have questions? Our support team is here to help.
            </p>
            <a href="/help" className="text-primary-500 hover:text-primary-600 font-semibold">
              Visit Help Center →
            </a>
          </div>
        </div>
      </div>
    </div>
  )
}

export default TrackOrderPage
