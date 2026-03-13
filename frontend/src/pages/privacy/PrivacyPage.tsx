function PrivacyPage() {
  const sections = [
    {
      title: '1. Information We Collect',
      content: `We collect information you provide directly to us, including your name, email address, postal address, phone number, and payment information when you make a purchase or create an account. We also automatically collect certain information about your device and browsing activity.`,
      subsections: [
        'Personal Information (name, email, address, phone)',
        'Payment Information (credit card details, billing address)',
        'Order History and Preferences',
        'Device and Browser Information',
        'Cookies and Usage Data'
      ]
    },
    {
      title: '2. How We Use Your Information',
      content: `We use the information we collect to process orders, communicate with you, improve our services, and send promotional materials (with your consent).`,
      subsections: [
        'Process and ship your orders',
        'Send order confirmations and updates',
        'Respond to your inquiries and provide support',
        'Send marketing communications (opt-in required)',
        'Improve our website and services',
        'Prevent fraud and ensure security'
      ]
    },
    {
      title: '3. Information Sharing',
      content: `We do not sell, trade, or rent your personal information to third parties. We may share your information with trusted service providers who assist us in operating our business, such as payment processors and shipping companies.`,
      subsections: [
        'Payment Processors (to process transactions)',
        'Shipping Companies (to deliver orders)',
        'Service Providers (IT, analytics, marketing)',
        'Legal Authorities (when required by law)'
      ]
    },
    {
      title: '4. Data Security',
      content: `We implement appropriate technical and organizational measures to protect your personal information against unauthorized access, alteration, disclosure, or destruction. However, no method of transmission over the Internet is 100% secure.`,
      subsections: [
        'SSL Encryption for all transactions',
        'Secure servers and databases',
        'Regular security audits',
        'Employee training on data protection'
      ]
    },
    {
      title: '5. Cookies',
      content: `We use cookies to enhance your browsing experience, analyze site traffic, and personalize content. You can control cookie settings through your browser, but disabling cookies may limit your use of certain features.`,
      subsections: [
        'Essential Cookies (required for site functionality)',
        'Analytics Cookies (to understand site usage)',
        'Marketing Cookies (for personalized ads)',
        'You can manage cookies in your browser settings'
      ]
    },
    {
      title: '6. Your Rights',
      content: `You have the right to access, correct, or delete your personal information. You can also object to or restrict certain processing of your data. Contact us to exercise these rights.`,
      subsections: [
        'Access your personal data',
        'Correct inaccurate data',
        'Request deletion of your data',
        'Opt-out of marketing communications',
        'Export your data'
      ]
    },
    {
      title: '7. Data Retention',
      content: `We retain your personal information for as long as necessary to fulfill the purposes outlined in this policy, unless a longer retention period is required by law.`,
      subsections: [
        'Account data: While your account is active',
        'Order data: As required by tax laws (7 years)',
        'Marketing data: Until you unsubscribe',
        'Analytics data: Anonymized after 26 months'
      ]
    },
    {
      title: '8. Children\'s Privacy',
      content: `Our services are not intended for children under 18 years of age. We do not knowingly collect personal information from children. If we discover that we have collected information from a child, we will delete it immediately.`,
      subsections: []
    },
    {
      title: '9. Changes to This Policy',
      content: `We may update this Privacy Policy from time to time. We will notify you of any changes by posting the new policy on this page and updating the "Last Updated" date.`,
      subsections: []
    },
    {
      title: '10. Contact Us',
      content: `If you have any questions about this Privacy Policy or our data practices, please contact us at privacy@eshop.com or visit our Contact page.`,
      subsections: []
    }
  ]

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="max-w-4xl mx-auto">
          <div className="text-center mb-12">
            <h1 className="text-4xl font-bold text-gray-900 mb-4">Privacy Policy</h1>
            <p className="text-gray-600">Last updated: March 13, 2026</p>
          </div>

          {/* Content */}
          <div className="bg-white rounded-lg shadow-lg p-8 md:p-12">
            <div className="prose prose-lg max-w-none">
              <p className="text-gray-600 mb-8">
                At E-Shop, we take your privacy seriously. This Privacy Policy explains how we collect, use, disclose, and safeguard your information when you visit our website or make purchases.
              </p>

              {sections.map((section, index) => (
                <div key={index} className="mb-8">
                  <h2 className="text-2xl font-bold text-gray-900 mb-4">{section.title}</h2>
                  <p className="text-gray-600 leading-relaxed mb-4">{section.content}</p>
                  {section.subsections && section.subsections.length > 0 && (
                    <ul className="list-disc list-inside space-y-2 text-gray-600 ml-4">
                      {section.subsections.map((item, i) => (
                        <li key={i}>{item}</li>
                      ))}
                    </ul>
                  )}
                </div>
              ))}

              {/* Acceptance */}
              <div className="mt-12 pt-8 border-t">
                <p className="text-gray-600">
                  By using E-Shop, you acknowledge that you have read and understood this Privacy Policy.
                </p>
              </div>
            </div>
          </div>

          {/* Related Links */}
          <div className="mt-8 grid grid-cols-1 md:grid-cols-2 gap-4">
            <a href="/terms" className="bg-white rounded-lg shadow-lg p-6 hover:shadow-xl transition">
              <h3 className="font-semibold text-gray-900 mb-2">Terms of Service</h3>
              <p className="text-gray-600 text-sm">Our terms and conditions</p>
            </a>
            <a href="/contact" className="bg-white rounded-lg shadow-lg p-6 hover:shadow-xl transition">
              <h3 className="font-semibold text-gray-900 mb-2">Contact Us</h3>
              <p className="text-gray-600 text-sm">Get in touch with our team</p>
            </a>
          </div>
        </div>
      </div>
    </div>
  )
}

export default PrivacyPage
