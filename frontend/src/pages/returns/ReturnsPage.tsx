import { Package, RotateCcw, DollarSign, Clock } from 'lucide-react'

function ReturnsPage() {
  const steps = [
    {
      icon: <Package className="w-8 h-8" />,
      title: '1. Initiate Return',
      description: 'Visit our Returns page and enter your order number and email address.'
    },
    {
      icon: <RotateCcw className="w-8 h-8" />,
      title: '2. Select Items',
      description: 'Choose the items you want to return and select a reason for each item.'
    },
    {
      icon: <Clock className="w-8 h-8" />,
      title: '3. Print Label',
      description: 'Print the prepaid return shipping label provided via email.'
    },
    {
      icon: <DollarSign className="w-8 h-8" />,
      title: '4. Get Refund',
      description: 'Once we receive your return, refunds are processed within 5-7 business days.'
    }
  ]

  const faqs = [
    {
      question: 'What is your return policy?',
      answer: 'Items can be returned within 30 days of delivery. Products must be in original condition with tags attached and original packaging.'
    },
    {
      question: 'Which items cannot be returned?',
      answer: 'For hygiene reasons, we cannot accept returns on personal care items, underwear, swimwear (if hygiene seal is broken), and perishable goods.'
    },
    {
      question: 'Who pays for return shipping?',
      answer: 'Return shipping is free for defective or incorrect items. For other returns, a flat shipping fee of $5.99 will be deducted from your refund.'
    },
    {
      question: 'How long does it take to get a refund?',
      answer: 'Refunds are processed within 5-7 business days after we receive and inspect your return. The refund will appear in your account within 3-5 business days after processing.'
    },
    {
      question: 'Can I exchange an item?',
      answer: 'Yes! You can exchange items for the same product in a different size or color. Subject to availability. Contact our support team to arrange an exchange.'
    }
  ]

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">Returns & Refunds</h1>
          <p className="text-gray-600 text-lg">Hassle-free returns within 30 days</p>
        </div>

        {/* Policy Summary */}
        <div className="bg-primary-500 text-white rounded-lg shadow-lg p-8 mb-12">
          <div className="text-center">
            <h2 className="text-3xl font-bold mb-4">30-Day Return Policy</h2>
            <p className="text-primary-100 mb-6 max-w-2xl mx-auto">
              Not satisfied with your purchase? Return it within 30 days of delivery for a full refund. No questions asked.
            </p>
            <a href="/account/orders" className="inline-block bg-white text-primary-500 font-semibold px-8 py-4 rounded-lg hover:bg-gray-100 transition">
              Start a Return
            </a>
          </div>
        </div>

        {/* Return Steps */}
        <div className="mb-16">
          <h2 className="text-3xl font-bold text-center mb-12">How to Return an Item</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            {steps.map((step, index) => (
              <div key={index} className="bg-white rounded-lg shadow-lg p-6 text-center">
                <div className="text-primary-500 flex justify-center mb-4">{step.icon}</div>
                <h3 className="text-lg font-bold mb-3">{step.title}</h3>
                <p className="text-gray-600">{step.description}</p>
              </div>
            ))}
          </div>
        </div>

        {/* Important Information */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-16">
          <div className="bg-white rounded-lg shadow-lg p-8">
            <h3 className="text-xl font-bold mb-4">Return Conditions</h3>
            <ul className="space-y-3 text-gray-600">
              <li className="flex items-start gap-3">
                <span className="text-green-500 mt-1">✓</span>
                <span>Items must be in original condition</span>
              </li>
              <li className="flex items-start gap-3">
                <span className="text-green-500 mt-1">✓</span>
                <span>Tags and packaging must be intact</span>
              </li>
              <li className="flex items-start gap-3">
                <span className="text-green-500 mt-1">✓</span>
                <span>Return must be initiated within 30 days</span>
              </li>
              <li className="flex items-start gap-3">
                <span className="text-green-500 mt-1">✓</span>
                <span>Proof of purchase required</span>
              </li>
            </ul>
          </div>

          <div className="bg-white rounded-lg shadow-lg p-8">
            <h3 className="text-xl font-bold mb-4">Non-Returnable Items</h3>
            <ul className="space-y-3 text-gray-600">
              <li className="flex items-start gap-3">
                <span className="text-red-500 mt-1">✗</span>
                <span>Personal care products (if opened)</span>
              </li>
              <li className="flex items-start gap-3">
                <span className="text-red-500 mt-1">✗</span>
                <span>Underwear and swimwear (hygiene seal broken)</span>
              </li>
              <li className="flex items-start gap-3">
                <span className="text-red-500 mt-1">✗</span>
                <span>Perishable goods</span>
              </li>
              <li className="flex items-start gap-3">
                <span className="text-red-500 mt-1">✗</span>
                <span>Customized or personalized items</span>
              </li>
            </ul>
          </div>
        </div>

        {/* FAQs */}
        <div className="max-w-3xl mx-auto mb-16">
          <h2 className="text-3xl font-bold text-center mb-8">Return FAQs</h2>
          <div className="space-y-6">
            {faqs.map((faq, index) => (
              <div key={index} className="bg-white rounded-lg shadow-md p-6">
                <h3 className="text-lg font-semibold mb-3 text-gray-900">{faq.question}</h3>
                <p className="text-gray-600">{faq.answer}</p>
              </div>
            ))}
          </div>
        </div>

        {/* Help Section */}
        <div className="max-w-2xl mx-auto text-center">
          <div className="bg-white rounded-lg shadow-lg p-8">
            <h3 className="text-xl font-bold text-gray-900 mb-4">Need Help with Returns?</h3>
            <p className="text-gray-600 mb-6">
              Our customer support team is here to assist you with any return-related questions.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <a href="/contact" className="text-primary-500 hover:text-primary-600 font-semibold">
                Contact Support
              </a>
              <span className="text-gray-300">|</span>
              <a href="/help" className="text-primary-500 hover:text-primary-600 font-semibold">
                Visit Help Center
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default ReturnsPage
