import React from 'react';
import { Link } from 'react-router-dom';
import {
  FaceSmileIcon as FacebookIcon,
  ChatBubbleLeftIcon as TwitterIcon,
  CameraIcon as InstagramIcon,
  PlayIcon as YoutubeIcon
} from '@heroicons/react/24/solid';
import { ROUTES } from '../../constants';

const Footer: React.FC = () => {
  return (
    <footer className="bg-gray-900 text-gray-300">
      {/* Main footer */}
      <div className="container mx-auto px-4 py-12">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          {/* Company info */}
          <div>
            <div className="flex items-center space-x-2 mb-4">
              <div className="w-10 h-10 bg-primary-600 rounded-lg flex items-center justify-center">
                <span className="text-white font-bold text-xl">E</span>
              </div>
              <span className="text-2xl font-bold text-white">E-Shop</span>
            </div>
            <p className="text-gray-400 mb-4">
              Your trusted online shopping destination for quality products at great prices.
            </p>
            <div className="flex space-x-4">
              <a href="#" className="hover:text-primary-500 transition-colors">
                <FacebookIcon className="w-6 h-6" />
              </a>
              <a href="#" className="hover:text-primary-500 transition-colors">
                <TwitterIcon className="w-6 h-6" />
              </a>
              <a href="#" className="hover:text-primary-500 transition-colors">
                <InstagramIcon className="w-6 h-6" />
              </a>
              <a href="#" className="hover:text-primary-500 transition-colors">
                <YoutubeIcon className="w-6 h-6" />
              </a>
            </div>
          </div>

          {/* Quick links */}
          <div>
            <h3 className="text-white font-semibold mb-4">Quick Links</h3>
            <ul className="space-y-2">
              <li>
                <Link to={ROUTES.HOME} className="hover:text-primary-500 transition-colors">
                  Home
                </Link>
              </li>
              <li>
                <Link to={ROUTES.PRODUCTS} className="hover:text-primary-500 transition-colors">
                  Products
                </Link>
              </li>
              <li>
                <Link to={ROUTES.CART} className="hover:text-primary-500 transition-colors">
                  Shopping Cart
                </Link>
              </li>
              <li>
                <Link to="/about" className="hover:text-primary-500 transition-colors">
                  About Us
                </Link>
              </li>
              <li>
                <Link to="/contact" className="hover:text-primary-500 transition-colors">
                  Contact
                </Link>
              </li>
            </ul>
          </div>

          {/* Customer service */}
          <div>
            <h3 className="text-white font-semibold mb-4">Customer Service</h3>
            <ul className="space-y-2">
              <li>
                <Link to="/help" className="hover:text-primary-500 transition-colors">
                  Help Center
                </Link>
              </li>
              <li>
                <Link to="/shipping" className="hover:text-primary-500 transition-colors">
                  Shipping Info
                </Link>
              </li>
              <li>
                <Link to="/returns" className="hover:text-primary-500 transition-colors">
                  Returns & Refunds
                </Link>
              </li>
              <li>
                <Link to="/faq" className="hover:text-primary-500 transition-colors">
                  FAQ
                </Link>
              </li>
              <li>
                <Link to="/terms" className="hover:text-primary-500 transition-colors">
                  Terms & Conditions
                </Link>
              </li>
            </ul>
          </div>

          {/* Contact info */}
          <div>
            <h3 className="text-white font-semibold mb-4">Contact Us</h3>
            <ul className="space-y-2 text-gray-400">
              <li>
                <span className="font-medium text-gray-300">Email:</span> support@eshop.com
              </li>
              <li>
                <span className="font-medium text-gray-300">Phone:</span> +1 (555) 123-4567
              </li>
              <li>
                <span className="font-medium text-gray-300">Address:</span> 123 Main St, City, Country
              </li>
              <li>
                <span className="font-medium text-gray-300">Hours:</span> Mon-Fri 9AM-6PM
              </li>
            </ul>
          </div>
        </div>
      </div>

      {/* Bottom bar */}
      <div className="border-t border-gray-800">
        <div className="container mx-auto px-4 py-4">
          <div className="flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">
            <p className="text-sm text-gray-400">
              © {new Date().getFullYear()} E-Shop. All rights reserved.
            </p>
            <div className="flex space-x-6 text-sm">
              <Link to="/privacy" className="hover:text-primary-500 transition-colors">
                Privacy Policy
              </Link>
              <Link to="/terms" className="hover:text-primary-500 transition-colors">
                Terms of Service
              </Link>
              <Link to="/cookies" className="hover:text-primary-500 transition-colors">
                Cookie Policy
              </Link>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
