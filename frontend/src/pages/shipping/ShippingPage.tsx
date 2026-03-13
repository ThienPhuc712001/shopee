import { Truck, Clock, DollarSign, Package } from 'lucide-react'

function ShippingPage() {
  const shippingMethods = [
    {
      icon: <Truck className="w-8 h-8" />,
      name: 'Standard Shipping',
      time: '5-7 business days',
      cost: '$5.99 - $9.99',
      description: 'Economical shipping for non-urgent deliveries'
    },
    {
      icon: <Clock className="w-8 h-8" />,
      name: 'Express Shipping',
      time: '2-3 business days',
      cost: '$12.99 - $19.99',
      description: 'Faster delivery for time-sensitive orders'
    },
    {
      icon: <Package className="w-8 h-8" />,
      name: 'Overnight Shipping',
      time: 'Next business day',
      cost: '$24.99 - $34.99',
      description: 'Fastest shipping option available'
    }
  ]

  const infoSections = [
    {
      title: 'Shipping Destinations',
      content: 'We ship nationwide across Vietnam. International shipping is available to selected countries. Shipping costs and delivery times vary based on destination.'
    },
    {
      title: 'Order Processing',
      content: 'Orders are processed within 1-2 business days. Orders placed after 2 PM or on weekends/holidays will be processed the next business day.'
    },
    {
      title: 'Shipping Restrictions',
      content: 'Some items may have shipping restrictions due to size, weight, or hazardous materials. We will notify you if any items in your order cannot be shipped to your address.'
    },
    {
      title: 'Delivery Issues',
      content: 'If your package is damaged or lost, contact us within 48 hours of delivery. We will investigate and provide a replacement or refund as appropriate.'
    }
  ]

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">Shipping Information</h1>
          <p className="text-gray-600 text-lg">Fast, reliable delivery to your doorstep</p>
        </div>

        {/* Shipping Methods */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-16">
          {shippingMethods.map((method, index) => (
            <div key={index} className="bg-white rounded-lg shadow-lg p-8 text-center hover:shadow-xl transition">
              <div className="text-primary-500 flex justify-center mb-4">{method.icon}</div>
              <h3 className="text-xl font-bold mb-2">{method.name}</h3>
              <div className="space-y-2 mb-4">
                <div className="flex items-center justify-center gap-2 text-gray-600">
                  <Clock className="w-4 h-4" />
                  <span>{method.time}</span>
                </div>
                <div className="flex items-center justify-center gap-2 text-gray-600">
                  <DollarSign className="w-4 h-4" />
                  <span>{method.cost}</span>
                </div>
              </div>
              <p className="text-gray-600">{method.description}</p>
            </div>
          ))}
        </div>

        {/* Free Shipping Banner */}
        <div className="bg-primary-500 text-white rounded-lg shadow-lg p-8 mb-16">
          <div className="text-center">
            <h2 className="text-3xl font-bold mb-4">🎉 Free Shipping on Orders Over $50!</h2>
            <p className="text-primary-100 mb-6">
              Enjoy free standard shipping on all orders over $50. No code needed, discount applied automatically at checkout.
            </p>
            <a href="/products" className="inline-block bg-white text-primary-500 font-semibold px-8 py-4 rounded-lg hover:bg-gray-100 transition">
              Start Shopping
            </a>
          </div>
        </div>

        {/* Information Sections */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-16">
          {infoSections.map((section, index) => (
            <div key={index} className="bg-white rounded-lg shadow-lg p-8">
              <h3 className="text-xl font-bold mb-4">{section.title}</h3>
              <p className="text-gray-600">{section.content}</p>
            </div>
          ))}
        </div>

        {/* Tracking Section */}
        <div className="bg-white rounded-lg shadow-lg p-8">
          <div className="text-center">
            <h2 className="text-2xl font-bold mb-4">Track Your Order</h2>
            <p className="text-gray-600 mb-6">
              Want to know where your package is? Track your order in real-time.
            </p>
            <a href="/track-order" className="text-primary-500 hover:text-primary-600 font-semibold">
              Track Order →
            </a>
          </div>
        </div>
      </div>
    </div>
  )
}

export default ShippingPage
