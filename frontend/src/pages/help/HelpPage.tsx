import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Search, MessageCircle, Phone, Mail, ChevronDown, ChevronUp } from 'lucide-react'

function HelpPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [openFaq, setOpenFaq] = useState<number | null>(null)

  const faqs = [
    {
      id: 1,
      question: 'How do I track my order?',
      answer: 'You can track your order by visiting the "Track Order" page and entering your order number and email address. You will receive real-time updates on your order status.'
    },
    {
      id: 2,
      question: 'What are the shipping options?',
      answer: 'We offer Standard Shipping (5-7 business days), Express Shipping (2-3 business days), and Overnight Shipping (next business day). Shipping costs vary based on location and weight.'
    },
    {
      id: 3,
      question: 'How do I return an item?',
      answer: 'Items can be returned within 30 days of delivery. Visit the "Returns" page to initiate a return request. Items must be in original condition with tags attached.'
    },
    {
      id: 4,
      question: 'What payment methods do you accept?',
      answer: 'We accept Credit Cards (Visa, MasterCard, American Express), Debit Cards, PayPal, Bank Transfer, and Cash on Delivery (COD) for selected areas.'
    },
    {
      id: 5,
      question: 'How do I contact customer support?',
      answer: 'You can reach us via email at support@eshop.com, call our hotline at 1900 xxxx, or use the live chat feature available on our website during business hours (9 AM - 6 PM, Mon-Fri).'
    },
    {
      id: 6,
      question: 'Can I cancel or modify my order?',
      answer: 'Orders can be cancelled or modified within 2 hours of placement if they haven\'t been shipped yet. Contact our support team immediately for assistance.'
    }
  ]

  const contactMethods = [
    {
      icon: <MessageCircle className="w-6 h-6" />,
      title: 'Live Chat',
      description: 'Chat with our support team',
      action: 'Start Chat',
      available: 'Available 24/7'
    },
    {
      icon: <Phone className="w-6 h-6" />,
      title: 'Phone Support',
      description: 'Call us at 1900 xxxx',
      action: 'Call Now',
      available: 'Mon-Fri, 9AM-6PM'
    },
    {
      icon: <Mail className="w-6 h-6" />,
      title: 'Email Support',
      description: 'support@eshop.com',
      action: 'Send Email',
      available: 'Response within 24h'
    }
  ]

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Hero Section */}
      <div className="bg-gradient-to-r from-primary-500 to-primary-600 text-white py-16">
        <div className="container mx-auto px-4 text-center">
          <h1 className="text-4xl font-bold mb-4">How can we help you?</h1>
          <p className="text-lg text-primary-100 mb-8">Find answers to common questions or contact our support team</p>
          
          {/* Search Bar */}
          <div className="max-w-2xl mx-auto">
            <div className="relative">
              <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-6 h-6 text-gray-400" />
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search for help articles, FAQs..."
                className="w-full pl-12 pr-4 py-4 rounded-lg text-gray-900 outline-none focus:ring-4 focus:ring-primary-300"
              />
            </div>
          </div>
        </div>
      </div>

      {/* Contact Methods */}
      <div className="container mx-auto px-4 -mt-8">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {contactMethods.map((method, index) => (
            <div key={index} className="bg-white rounded-lg shadow-lg p-6 hover:shadow-xl transition">
              <div className="text-primary-500 mb-4">{method.icon}</div>
              <h3 className="text-xl font-semibold mb-2">{method.title}</h3>
              <p className="text-gray-600 mb-2">{method.description}</p>
              <p className="text-sm text-green-600 mb-4">{method.available}</p>
              <button className="w-full bg-primary-500 hover:bg-primary-600 text-white font-semibold py-2 rounded-lg transition">
                {method.action}
              </button>
            </div>
          ))}
        </div>
      </div>

      {/* FAQ Section */}
      <div className="container mx-auto px-4 py-16">
        <h2 className="text-3xl font-bold text-center mb-12">Frequently Asked Questions</h2>
        <div className="max-w-3xl mx-auto space-y-4">
          {faqs.map((faq) => (
            <div key={faq.id} className="bg-white rounded-lg shadow-md overflow-hidden">
              <button
                onClick={() => setOpenFaq(openFaq === faq.id ? null : faq.id)}
                className="w-full px-6 py-4 text-left flex items-center justify-between hover:bg-gray-50 transition"
              >
                <span className="font-semibold text-gray-900">{faq.question}</span>
                {openFaq === faq.id ? (
                  <ChevronUp className="w-5 h-5 text-gray-500 flex-shrink-0" />
                ) : (
                  <ChevronDown className="w-5 h-5 text-gray-500 flex-shrink-0" />
                )}
              </button>
              {openFaq === faq.id && (
                <div className="px-6 pb-4 text-gray-600 border-t">
                  <p className="pt-4">{faq.answer}</p>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>

      {/* Still Need Help Section */}
      <div className="bg-secondary-800 text-white py-16">
        <div className="container mx-auto px-4 text-center">
          <h2 className="text-3xl font-bold mb-4">Still need help?</h2>
          <p className="text-gray-300 mb-8 max-w-2xl mx-auto">
            Our support team is here to assist you with any questions or concerns you may have.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link to="/contact" className="bg-primary-500 hover:bg-primary-600 text-white font-semibold px-8 py-4 rounded-lg transition">
              Contact Us
            </Link>
            <Link to="/faq" className="bg-white/10 hover:bg-white/20 text-white font-semibold px-8 py-4 rounded-lg transition">
              View All FAQs
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}

export default HelpPage
