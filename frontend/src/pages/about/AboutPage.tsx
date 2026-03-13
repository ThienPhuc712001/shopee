import { Link } from 'react-router-dom'
import { Users, Target, Award, Heart } from 'lucide-react'

function AboutPage() {
  const team = [
    { name: 'John Smith', role: 'CEO & Founder', image: '👨‍💼' },
    { name: 'Sarah Johnson', role: 'CTO', image: '👩‍💻' },
    { name: 'Michael Chen', role: 'Head of Operations', image: '👨‍💼' },
    { name: 'Emily Davis', role: 'Customer Success', image: '👩‍💼' },
  ]

  const values = [
    { icon: <Target className="w-8 h-8" />, title: 'Our Mission', description: 'To provide the best online shopping experience with quality products and exceptional service.' },
    { icon: <Award className="w-8 h-8" />, title: 'Quality First', description: 'We carefully curate our products to ensure only the highest quality reaches our customers.' },
    { icon: <Heart className="w-8 h-8" />, title: 'Customer Focus', description: 'Your satisfaction is our priority. We are here to support you every step of the way.' },
    { icon: <Users className="w-8 h-8" />, title: 'Community', description: 'Building a trusted community of shoppers who value quality and authenticity.' },
  ]

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Hero Section */}
      <div className="bg-gradient-to-r from-secondary-800 to-secondary-700 text-white py-20">
        <div className="container mx-auto px-4 text-center">
          <h1 className="text-5xl font-bold mb-6">About E-Shop</h1>
          <p className="text-xl text-gray-300 max-w-3xl mx-auto">
            Your trusted partner in online shopping since 2020. We bring quality products directly to your doorstep.
          </p>
        </div>
      </div>

      {/* Our Story */}
      <div className="container mx-auto px-4 py-16">
        <div className="max-w-4xl mx-auto">
          <h2 className="text-3xl font-bold text-gray-900 mb-6 text-center">Our Story</h2>
          <div className="prose prose-lg mx-auto text-gray-600">
            <p className="mb-4">
              E-Shop was founded in 2020 with a simple mission: to make quality products accessible to everyone through a seamless online shopping experience.
            </p>
            <p className="mb-4">
              What started as a small online store has grown into a thriving e-commerce platform serving thousands of customers across the country. Our success is built on trust, quality, and customer satisfaction.
            </p>
            <p>
              Today, we continue to innovate and expand our product range while maintaining our commitment to excellent service and competitive prices.
            </p>
          </div>
        </div>
      </div>

      {/* Values */}
      <div className="bg-white py-16">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold text-center mb-12">Our Values</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            {values.map((value, index) => (
              <div key={index} className="text-center p-6">
                <div className="text-primary-500 flex justify-center mb-4">{value.icon}</div>
                <h3 className="text-xl font-semibold mb-3">{value.title}</h3>
                <p className="text-gray-600">{value.description}</p>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Team */}
      <div className="container mx-auto px-4 py-16">
        <h2 className="text-3xl font-bold text-center mb-12">Meet Our Team</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
          {team.map((member, index) => (
            <div key={index} className="bg-white rounded-lg shadow-lg p-6 text-center hover:shadow-xl transition">
              <div className="text-6xl mb-4">{member.image}</div>
              <h3 className="text-xl font-semibold mb-2">{member.name}</h3>
              <p className="text-primary-500">{member.role}</p>
            </div>
          ))}
        </div>
      </div>

      {/* Stats */}
      <div className="bg-primary-500 text-white py-16">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8 text-center">
            <div>
              <div className="text-4xl font-bold mb-2">50,000+</div>
              <div className="text-primary-100">Happy Customers</div>
            </div>
            <div>
              <div className="text-4xl font-bold mb-2">10,000+</div>
              <div className="text-primary-100">Products</div>
            </div>
            <div>
              <div className="text-4xl font-bold mb-2">100+</div>
              <div className="text-primary-100">Brands</div>
            </div>
            <div>
              <div className="text-4xl font-bold mb-2">24/7</div>
              <div className="text-primary-100">Support</div>
            </div>
          </div>
        </div>
      </div>

      {/* CTA */}
      <div className="container mx-auto px-4 py-16 text-center">
        <h2 className="text-3xl font-bold mb-4">Ready to Start Shopping?</h2>
        <p className="text-gray-600 mb-8">Join thousands of satisfied customers today.</p>
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <Link to="/products" className="bg-primary-500 hover:bg-primary-600 text-white font-semibold px-8 py-4 rounded-lg transition">
            Browse Products
          </Link>
          <Link to="/contact" className="bg-secondary-800 hover:bg-secondary-700 text-white font-semibold px-8 py-4 rounded-lg transition">
            Contact Us
          </Link>
        </div>
      </div>
    </div>
  )
}

export default AboutPage
