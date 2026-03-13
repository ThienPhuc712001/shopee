import { useState } from 'react'
import { Search, ChevronDown, ChevronUp, HelpCircle } from 'lucide-react'

function FAQPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [openFaq, setOpenFaq] = useState<number | null>(null)
  const [activeCategory, setActiveCategory] = useState('all')

  const categories = [
    { id: 'all', name: 'All Topics' },
    { id: 'orders', name: 'Orders & Shipping' },
    { id: 'payment', name: 'Payment' },
    { id: 'returns', name: 'Returns & Refunds' },
    { id: 'account', name: 'Account' },
    { id: 'products', name: 'Products' }
  ]

  const faqs = [
    {
      id: 1,
      category: 'orders',
      question: 'How do I place an order?',
      answer: 'Browse our products, add items to your cart, proceed to checkout, enter your shipping address and payment information, then confirm your order. You will receive an email confirmation shortly after.'
    },
    {
      id: 2,
      category: 'orders',
      question: 'How long does shipping take?',
      answer: 'Standard Shipping: 5-7 business days, Express Shipping: 2-3 business days, Overnight Shipping: Next business day. Delivery times may vary based on your location.'
    },
    {
      id: 3,
      category: 'orders',
      question: 'Can I track my order?',
      answer: 'Yes! Once your order ships, you will receive a tracking number via email. You can also track your order from the "Track Order" page using your order number and email.'
    },
    {
      id: 4,
      category: 'payment',
      question: 'What payment methods do you accept?',
      answer: 'We accept Visa, MasterCard, American Express, PayPal, Bank Transfer, and Cash on Delivery (COD) for selected areas.'
    },
    {
      id: 5,
      category: 'payment',
      question: 'Is my payment information secure?',
      answer: 'Absolutely! We use industry-standard SSL encryption to protect your payment information. Your credit card details are never stored on our servers.'
    },
    {
      id: 6,
      category: 'returns',
      question: 'What is your return policy?',
      answer: 'Items can be returned within 30 days of delivery. Products must be in original condition with tags attached. Some items like personal care products and underwear cannot be returned for hygiene reasons.'
    },
    {
      id: 7,
      category: 'returns',
      question: 'How do I return an item?',
      answer: 'Visit the Returns page, enter your order number, select the items you want to return, choose a reason, and print the prepaid return label. Drop off the package at any authorized location.'
    },
    {
      id: 8,
      category: 'returns',
      question: 'When will I receive my refund?',
      answer: 'Refunds are processed within 5-7 business days after we receive your return. The refund will be credited to your original payment method.'
    },
    {
      id: 9,
      category: 'account',
      question: 'How do I create an account?',
      answer: 'Click on "Sign In" in the header, then select "Create Account". Enter your email, create a password, and fill in your details. You can also sign up using your social media accounts.'
    },
    {
      id: 10,
      category: 'account',
      question: 'I forgot my password. What should I do?',
      answer: 'On the login page, click "Forgot Password". Enter your email address and we will send you a link to reset your password.'
    },
    {
      id: 11,
      category: 'products',
      question: 'Are your products authentic?',
      answer: 'Yes! All our products are 100% authentic and sourced directly from authorized brands and distributors.'
    },
    {
      id: 12,
      category: 'products',
      question: 'Do you offer gift wrapping?',
      answer: 'Yes, gift wrapping is available for select items. You can add gift wrapping during checkout for an additional fee.'
    }
  ]

  const filteredFaqs = faqs.filter(faq => {
    const matchesCategory = activeCategory === 'all' || faq.category === activeCategory
    const matchesSearch = faq.question.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         faq.answer.toLowerCase().includes(searchQuery.toLowerCase())
    return matchesCategory && matchesSearch
  })

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">Frequently Asked Questions</h1>
          <p className="text-gray-600 text-lg">Find answers to common questions</p>
        </div>

        {/* Search Bar */}
        <div className="max-w-2xl mx-auto mb-12">
          <div className="relative">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search for answers..."
              className="w-full pl-12 pr-4 py-4 rounded-lg border border-gray-300 focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
            />
          </div>
        </div>

        {/* Categories */}
        <div className="flex flex-wrap justify-center gap-3 mb-12">
          {categories.map((category) => (
            <button
              key={category.id}
              onClick={() => setActiveCategory(category.id)}
              className={`px-6 py-3 rounded-full font-medium transition ${
                activeCategory === category.id
                  ? 'bg-primary-500 text-white'
                  : 'bg-white text-gray-700 hover:bg-gray-100'
              }`}
            >
              {category.name}
            </button>
          ))}
        </div>

        {/* FAQs */}
        <div className="max-w-3xl mx-auto space-y-4">
          {filteredFaqs.length > 0 ? (
            filteredFaqs.map((faq) => (
              <div key={faq.id} className="bg-white rounded-lg shadow-md overflow-hidden">
                <button
                  onClick={() => setOpenFaq(openFaq === faq.id ? null : faq.id)}
                  className="w-full px-6 py-4 text-left flex items-center justify-between hover:bg-gray-50 transition"
                >
                  <div className="flex items-center gap-3">
                    <HelpCircle className="w-5 h-5 text-primary-500 flex-shrink-0" />
                    <span className="font-semibold text-gray-900">{faq.question}</span>
                  </div>
                  {openFaq === faq.id ? (
                    <ChevronUp className="w-5 h-5 text-gray-500 flex-shrink-0" />
                  ) : (
                    <ChevronDown className="w-5 h-5 text-gray-500 flex-shrink-0" />
                  )}
                </button>
                {openFaq === faq.id && (
                  <div className="px-6 pb-4 pl-14 text-gray-600 border-t">
                    <p className="pt-4">{faq.answer}</p>
                  </div>
                )}
              </div>
            ))
          ) : (
            <div className="text-center py-12">
              <HelpCircle className="w-16 h-16 text-gray-300 mx-auto mb-4" />
              <p className="text-gray-600">No results found for "{searchQuery}"</p>
              <button
                onClick={() => { setSearchQuery(''); setActiveCategory('all') }}
                className="mt-4 text-primary-500 hover:text-primary-600 font-medium"
              >
                Clear filters
              </button>
            </div>
          )}
        </div>

        {/* Still Need Help */}
        <div className="max-w-2xl mx-auto mt-12 text-center">
          <div className="bg-white rounded-lg shadow-lg p-8">
            <h3 className="text-xl font-bold text-gray-900 mb-4">Still have questions?</h3>
            <p className="text-gray-600 mb-6">
              Our support team is here to help you with any questions.
            </p>
            <a href="/contact" className="text-primary-500 hover:text-primary-600 font-semibold">
              Contact Support →
            </a>
          </div>
        </div>
      </div>
    </div>
  )
}

export default FAQPage
