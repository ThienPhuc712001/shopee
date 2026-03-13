import { Link } from 'react-router-dom'
import { ShoppingCart, Heart, Star, Eye } from 'lucide-react'
import { motion } from 'framer-motion'
import { Product } from '../../types'

interface ProductCardProps {
  product: Product
  onAddToCart?: (product: Product) => void
}

function ProductCard({ product, onAddToCart }: ProductCardProps) {
  const discount = product.original_price
    ? Math.round(((product.original_price - product.price) / product.original_price) * 100)
    : 0

  const handleAddToCart = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    if (onAddToCart) {
      onAddToCart(product)
    }
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="group card"
      whileHover={{ y: -5, scale: 1.02 }}
    >
      {/* Image */}
      <div className="relative overflow-hidden aspect-square bg-gradient-to-br from-gray-100 to-gray-200">
        {product.images && product.images.length > 0 && product.images[0].url ? (
          <img
            src={product.images[0].url}
            alt={product.name}
            className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
            onError={(e) => {
              (e.target as HTMLImageElement).src = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="300" height="300"%3E%3Crect fill="%23f3f4f6" width="300" height="300"/%3E%3Ctext fill="%239ca3af" font-family="sans-serif" font-size="18" x="50%25" y="50%25" dominant-baseline="middle" text-anchor="middle"%3EProduct Image%3C/text%3E%3C/svg%3E'
            }}
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-gray-400">
            <div className="text-center">
              <div className="text-6xl mb-2">📦</div>
              <p className="text-sm font-medium">No Image</p>
            </div>
          </div>
        )}
        
        {/* Discount Badge */}
        {discount > 0 && (
          <span className="absolute top-2 left-2 bg-red-500 text-white px-2 py-1 rounded-md text-sm font-bold">
            -{discount}%
          </span>
        )}

        {/* Quick Actions */}
        <div className="absolute top-2 right-2 flex flex-col gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
          <button className="w-10 h-10 bg-white rounded-full shadow-md flex items-center justify-center hover:bg-primary-500 hover:text-white transition">
            <Heart className="w-5 h-5" />
          </button>
          <button className="w-10 h-10 bg-white rounded-full shadow-md flex items-center justify-center hover:bg-primary-500 hover:text-white transition">
            <Eye className="w-5 h-5" />
          </button>
        </div>

        {/* Add to Cart Button */}
        <button 
          onClick={handleAddToCart}
          className="absolute bottom-0 left-0 right-0 bg-primary-500 text-white py-3 font-semibold opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2 hover:bg-primary-600"
        >
          <ShoppingCart className="w-5 h-5" />
          Add to Cart
        </button>
      </div>

      {/* Product Info */}
      <div className="p-4">
        <Link to={`/products/${product.id}`}>
          <h3 className="font-medium text-gray-900 mb-2 line-clamp-2 group-hover:text-primary-500 transition">
            {product.name}
          </h3>
        </Link>

        {/* Rating */}
        {product.rating && (
          <div className="flex items-center gap-1 mb-2">
            <Star className="w-4 h-4 fill-yellow-400 text-yellow-400" />
            <span className="text-sm text-gray-600">{product.rating}</span>
            <span className="text-sm text-gray-400">({product.review_count} reviews)</span>
          </div>
        )}

        {/* Price */}
        <div className="flex items-center gap-2">
          <span className="text-xl font-bold text-primary-500">
            ${product.price.toFixed(2)}
          </span>
          {product.original_price && (
            <span className="text-gray-400 line-through">
              ${product.original_price.toFixed(2)}
            </span>
          )}
        </div>

        {/* Stock Status */}
        {product.stock === 0 ? (
          <p className="text-red-500 text-sm font-medium mt-2">Out of Stock</p>
        ) : product.stock < 10 ? (
          <p className="text-orange-500 text-sm font-medium mt-2">Only {product.stock} left</p>
        ) : null}
      </div>
    </motion.div>
  )
}

export default ProductCard
