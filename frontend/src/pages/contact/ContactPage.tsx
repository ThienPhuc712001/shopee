import { useState } from 'react'
import { Mail, Phone, MapPin, Send, Clock } from 'lucide-react'

function ContactPage() {
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    subject: '',
    message: ''
  })
  const [submitted, setSubmitted] = useState(false)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitted(true)
    // In real app, send to backend
    setTimeout(() => {
      setSubmitted(false)
      setFormData({ name: '', email: '', subject: '', message: '' })
      alert('Thank you! We will get back to you soon.')
    }, 2000)
  }

  const contactInfo = [
    {
      icon: <MapPin className="w-6 h-6" />,
      title: 'Visit Us',
      lines: ['123 E-Commerce Street', 'Tech City, Vietnam 700000']
    },
    {
      icon: <Phone className="w-6 h-6" />,
      title: 'Call Us',
      lines: ['Hotline: 1900 xxxx', 'Phone: +84 123 456 789']
    },
    {
      icon: <Mail className="w-6 h-6" />,
      title: 'Email Us',
      lines: ['support@eshop.com', 'sales@eshop.com']
    },
    {
      icon: <Clock className="w-6 h-6" />,
      title: 'Business Hours',
      lines: ['Mon-Fri: 9AM - 6PM', 'Sat: 9AM - 12PM']
    }
  ]

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">Contact Us</h1>
          <p className="text-gray-600 text-lg">Get in touch with our team</p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 max-w-6xl mx-auto">
          {/* Contact Information */}
          <div className="lg:col-span-1 space-y-6">
            {contactInfo.map((info, index) => (
              <div key={index} className="bg-white rounded-lg shadow-lg p-6">
                <div className="text-primary-500 mb-4">{info.icon}</div>
                <h3 className="text-lg font-semibold mb-3">{info.title}</h3>
                {info.lines.map((line, i) => (
                  <p key={i} className="text-gray-600">{line}</p>
                ))}
              </div>
            ))}

            {/* Map Placeholder */}
            <div className="bg-white rounded-lg shadow-lg p-6">
              <h3 className="text-lg font-semibold mb-4">Our Location</h3>
              <div className="bg-gray-200 rounded-lg h-48 flex items-center justify-center">
                <MapPin className="w-12 h-12 text-gray-400" />
              </div>
            </div>
          </div>

          {/* Contact Form */}
          <div className="lg:col-span-2">
            <div className="bg-white rounded-lg shadow-lg p-8">
              <h2 className="text-2xl font-bold mb-6">Send us a Message</h2>
              <form onSubmit={handleSubmit} className="space-y-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Your Name
                    </label>
                    <input
                      type="text"
                      value={formData.name}
                      onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                      className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                      placeholder="John Doe"
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Email Address
                    </label>
                    <input
                      type="email"
                      value={formData.email}
                      onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                      className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                      placeholder="john@example.com"
                      required
                    />
                  </div>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Subject
                  </label>
                  <input
                    type="text"
                    value={formData.subject}
                    onChange={(e) => setFormData({ ...formData, subject: e.target.value })}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                    placeholder="How can we help?"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Message
                  </label>
                  <textarea
                    value={formData.message}
                    onChange={(e) => setFormData({ ...formData, message: e.target.value })}
                    rows={6}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none resize-none"
                    placeholder="Tell us more about your inquiry..."
                    required
                  ></textarea>
                </div>
                <button
                  type="submit"
                  disabled={submitted}
                  className="w-full bg-primary-500 hover:bg-primary-600 disabled:bg-gray-400 text-white font-semibold py-4 rounded-lg transition flex items-center justify-center gap-2"
                >
                  {submitted ? (
                    <>Sending...</>
                  ) : (
                    <>
                      <Send className="w-5 h-5" />
                      Send Message
                    </>
                  )}
                </button>
              </form>
            </div>
          </div>
        </div>

        {/* FAQ Link */}
        <div className="max-w-2xl mx-auto mt-12 text-center">
          <div className="bg-white rounded-lg shadow-lg p-8">
            <h3 className="text-xl font-bold text-gray-900 mb-4">Have a Question?</h3>
            <p className="text-gray-600 mb-6">
              Check our FAQ section for quick answers to common questions.
            </p>
            <a href="/faq" className="text-primary-500 hover:text-primary-600 font-semibold">
              View FAQs →
            </a>
          </div>
        </div>
      </div>
    </div>
  )
}

export default ContactPage
