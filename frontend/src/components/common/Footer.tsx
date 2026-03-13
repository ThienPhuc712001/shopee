import { Link } from 'react-router-dom'
import { Facebook, Instagram, Twitter, Mail, Phone, MapPin } from 'lucide-react'

function Footer() {
  return (
    <footer className="bg-secondary-800 text-white mt-auto">
      {/* Main Footer */}
      <div className="container mx-auto px-4 py-12">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          {/* Company Info */}
          <div>
            <div className="flex items-center gap-2 mb-4">
              <div className="w-10 h-10 bg-gradient-to-br from-primary-500 to-primary-700 rounded-lg flex items-center justify-center">
                <span className="text-white font-bold text-xl">E</span>
              </div>
              <span className="text-2xl font-bold">E-Shop</span>
            </div>
            <p className="text-gray-300 mb-4">
              Your trusted online shopping destination for quality products at great prices.
            </p>
            <div className="flex gap-4">
              <a href="#" className="w-10 h-10 bg-gray-700 hover:bg-primary-500 rounded-full flex items-center justify-center transition">
                <Facebook className="w-5 h-5" />
              </a>
              <a href="#" className="w-10 h-10 bg-gray-700 hover:bg-primary-500 rounded-full flex items-center justify-center transition">
                <Instagram className="w-5 h-5" />
              </a>
              <a href="#" className="w-10 h-10 bg-gray-700 hover:bg-primary-500 rounded-full flex items-center justify-center transition">
                <Twitter className="w-5 h-5" />
              </a>
            </div>
          </div>

          {/* Quick Links */}
          <div>
            <h3 className="text-lg font-semibold mb-4">Quick Links</h3>
            <ul className="space-y-2">
              <li><Link to="/about" className="text-gray-300 hover:text-primary-500 transition">About Us</Link></li>
              <li><Link to="/contact" className="text-gray-300 hover:text-primary-500 transition">Contact Us</Link></li>
              <li><Link to="/faq" className="text-gray-300 hover:text-primary-500 transition">FAQ</Link></li>
              <li><Link to="/shipping" className="text-gray-300 hover:text-primary-500 transition">Shipping Info</Link></li>
              <li><Link to="/returns" className="text-gray-300 hover:text-primary-500 transition">Returns & Refunds</Link></li>
            </ul>
          </div>

          {/* Customer Service */}
          <div>
            <h3 className="text-lg font-semibold mb-4">Customer Service</h3>
            <ul className="space-y-2">
              <li><Link to="/track-order" className="text-gray-300 hover:text-primary-500 transition">Track Order</Link></li>
              <li><Link to="/account/orders" className="text-gray-300 hover:text-primary-500 transition">My Orders</Link></li>
              <li><Link to="/account/wishlist" className="text-gray-300 hover:text-primary-500 transition">Wishlist</Link></li>
              <li><Link to="/help" className="text-gray-300 hover:text-primary-500 transition">Help Center</Link></li>
              <li><Link to="/terms" className="text-gray-300 hover:text-primary-500 transition">Terms of Service</Link></li>
            </ul>
          </div>

          {/* Contact Info */}
          <div>
            <h3 className="text-lg font-semibold mb-4">Contact Us</h3>
            <ul className="space-y-3">
              <li className="flex items-start gap-3">
                <MapPin className="w-5 h-5 text-primary-500 flex-shrink-0 mt-0.5" />
                <span className="text-gray-300">123 E-Commerce Street, Tech City, Vietnam</span>
              </li>
              <li className="flex items-center gap-3">
                <Phone className="w-5 h-5 text-primary-500 flex-shrink-0" />
                <span className="text-gray-300">1900 xxxx</span>
              </li>
              <li className="flex items-center gap-3">
                <Mail className="w-5 h-5 text-primary-500 flex-shrink-0" />
                <span className="text-gray-300">support@eshop.com</span>
              </li>
            </ul>
          </div>
        </div>
      </div>

      {/* Bottom Bar */}
      <div className="border-t border-gray-700">
        <div className="container mx-auto px-4 py-4">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4">
            <p className="text-gray-400 text-sm">
              © 2026 E-Shop. All rights reserved.
            </p>
            <div className="flex gap-6">
              <Link to="/privacy" className="text-gray-400 hover:text-primary-500 text-sm transition">
                Privacy Policy
              </Link>
              <Link to="/terms" className="text-gray-400 hover:text-primary-500 text-sm transition">
                Terms of Service
              </Link>
              <Link to="/cookies" className="text-gray-400 hover:text-primary-500 text-sm transition">
                Cookie Policy
              </Link>
            </div>
          </div>
        </div>
      </div>
    </footer>
  )
}

export default Footer
