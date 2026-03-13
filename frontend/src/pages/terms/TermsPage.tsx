function TermsPage() {
  const sections = [
    {
      title: '1. Terms of Use',
      content: `By accessing and using E-Shop, you accept and agree to be bound by the terms and provision of this agreement. If you do not agree to abide by these terms, please do not use this service.`
    },
    {
      title: '2. User Account',
      content: `You are responsible for maintaining the confidentiality of your account and password. You agree to accept responsibility for all activities that occur under your account. You must notify us immediately of any unauthorized use of your account.`
    },
    {
      title: '3. Product Information',
      content: `We make every effort to display as accurately as possible the colors, features, specifications, and details of the products available on the Store. However, we do not guarantee that the colors, features, specifications, and details of the products will be accurate, complete, reliable, current, or free of other errors.`
    },
    {
      title: '4. Pricing',
      content: `Prices for products are subject to change without notice. We reserve the right to modify or discontinue any product or service without notice at any time. We shall not be liable to you or any third party for any modification, price change, suspension or discontinuance of the service.`
    },
    {
      title: '5. Orders and Payment',
      content: `All orders are subject to acceptance and availability. We reserve the right to refuse or cancel any order for any reason. We accept various payment methods including credit cards, debit cards, and online payment systems. You agree to provide current, complete and accurate purchase and account information.`
    },
    {
      title: '6. Shipping and Delivery',
      content: `Shipping times are estimates and not guaranteed. We are not responsible for delays caused by customs, weather, or other factors beyond our control. Risk of loss passes to you upon delivery of the items to the carrier.`
    },
    {
      title: '7. Returns and Refunds',
      content: `Items may be returned within 30 days of delivery in original condition. Refunds will be processed within 5-7 business days after receipt of returned items. Shipping costs are non-refundable unless the return is due to our error.`
    },
    {
      title: '8. Intellectual Property',
      content: `All content included on this site, such as text, graphics, logos, images, and software, is the property of E-Shop and protected by copyright laws. You may not use, reproduce, or distribute any content without our express written permission.`
    },
    {
      title: '9. Limitation of Liability',
      content: `E-Shop shall not be liable for any indirect, incidental, special, consequential or punitive damages, including without limitation, loss of profits, data, use, goodwill, or other intangible losses, resulting from your access to or use of or inability to access or use the service.`
    },
    {
      title: '10. Changes to Terms',
      content: `We reserve the right to modify these terms at any time. Continued use of the service after any such changes shall constitute your consent to such changes. It is your responsibility to check these terms periodically for updates.`
    },
    {
      title: '11. Governing Law',
      content: `These Terms shall be governed and construed in accordance with the laws of Vietnam, without regard to its conflict of law provisions. Any disputes shall be resolved in the courts of Vietnam.`
    },
    {
      title: '12. Contact Information',
      content: `For questions about these Terms of Service, please contact us at support@eshop.com or visit our Contact page.`
    }
  ]

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="max-w-4xl mx-auto">
          <div className="text-center mb-12">
            <h1 className="text-4xl font-bold text-gray-900 mb-4">Terms of Service</h1>
            <p className="text-gray-600">Last updated: March 13, 2026</p>
          </div>

          {/* Content */}
          <div className="bg-white rounded-lg shadow-lg p-8 md:p-12">
            <div className="prose prose-lg max-w-none">
              <p className="text-gray-600 mb-8">
                Welcome to E-Shop. Please read these Terms of Service carefully before using our website or services.
              </p>

              {sections.map((section, index) => (
                <div key={index} className="mb-8">
                  <h2 className="text-2xl font-bold text-gray-900 mb-4">{section.title}</h2>
                  <p className="text-gray-600 leading-relaxed">{section.content}</p>
                </div>
              ))}

              {/* Acceptance */}
              <div className="mt-12 pt-8 border-t">
                <p className="text-gray-600">
                  By using E-Shop, you acknowledge that you have read, understood, and agree to be bound by these Terms of Service.
                </p>
              </div>
            </div>
          </div>

          {/* Related Links */}
          <div className="mt-8 grid grid-cols-1 md:grid-cols-2 gap-4">
            <a href="/privacy" className="bg-white rounded-lg shadow-lg p-6 hover:shadow-xl transition">
              <h3 className="font-semibold text-gray-900 mb-2">Privacy Policy</h3>
              <p className="text-gray-600 text-sm">Learn how we protect your privacy</p>
            </a>
            <a href="/returns" className="bg-white rounded-lg shadow-lg p-6 hover:shadow-xl transition">
              <h3 className="font-semibold text-gray-900 mb-2">Returns Policy</h3>
              <p className="text-gray-600 text-sm">Information about returns and refunds</p>
            </a>
          </div>
        </div>
      </div>
    </div>
  )
}

export default TermsPage
