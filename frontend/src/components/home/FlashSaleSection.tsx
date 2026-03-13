import { ShoppingBag, Clock, Zap } from 'lucide-react'

const flashProducts = [
  { id: 1, name: 'Wireless Headphones', price: 79.99, originalPrice: 149.99, discount: 47, image: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=300&h=300&fit=crop', sold: 156, total: 200 },
  { id: 2, name: 'Smart Watch', price: 199.99, originalPrice: 399.99, discount: 50, image: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=300&h=300&fit=crop', sold: 89, total: 150 },
  { id: 3, name: 'Camera Lens', price: 299.99, originalPrice: 599.99, discount: 50, image: 'https://images.unsplash.com/photo-1617005082133-548c4dd27f35?w=300&h=300&fit=crop', sold: 45, total: 100 },
  { id: 4, name: 'Gaming Mouse', price: 49.99, originalPrice: 99.99, discount: 50, image: 'https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?w=300&h=300&fit=crop', sold: 234, total: 300 },
]

function FlashSaleSection() {
  return (
    <section className="py-12 bg-gradient-to-r from-red-50 to-orange-50">
      <div className="container mx-auto px-4">
        <div className="flex flex-col md:flex-row justify-between items-center mb-8 gap-4">
          <div className="flex items-center gap-4">
            <div className="w-16 h-16 bg-red-500 rounded-full flex items-center justify-center animate-pulse">
              <Zap className="w-8 h-8 text-white" />
            </div>
            <div>
              <h2 className="text-3xl font-bold text-gray-900">Flash Sale</h2>
              <p className="text-gray-600">Limited time offers - Grab them fast!</p>
            </div>
          </div>
          <div className="flex items-center gap-4 bg-white px-6 py-4 rounded-xl shadow-md">
            <Clock className="w-6 h-6 text-red-500" />
            <div className="flex gap-2">
              <div className="text-center">
                <div className="w-12 h-12 bg-red-500 text-white rounded-lg flex items-center justify-center font-bold text-xl">02</div>
                <span className="text-xs text-gray-600 mt-1">Hours</span>
              </div>
              <div className="text-2xl font-bold text-red-500">:</div>
              <div className="text-center">
                <div className="w-12 h-12 bg-red-500 text-white rounded-lg flex items-center justify-center font-bold text-xl">45</div>
                <span className="text-xs text-gray-600 mt-1">Mins</span>
              </div>
              <div className="text-2xl font-bold text-red-500">:</div>
              <div className="text-center">
                <div className="w-12 h-12 bg-red-500 text-white rounded-lg flex items-center justify-center font-bold text-xl">30</div>
                <span className="text-xs text-gray-600 mt-1">Secs</span>
              </div>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {flashProducts.map((product) => (
            <div key={product.id} className="bg-white rounded-xl overflow-hidden shadow-lg hover:shadow-xl transition-all hover:scale-105">
              <div className="relative">
                <img src={product.image} alt={product.name} className="w-full h-48 object-cover" />
                <span className="absolute top-2 right-2 bg-red-500 text-white px-3 py-1 rounded-full text-sm font-bold">
                  -{product.discount}%
                </span>
              </div>
              <div className="p-4">
                <h3 className="font-semibold text-gray-900 mb-2">{product.name}</h3>
                <div className="flex items-center gap-2 mb-3">
                  <span className="text-2xl font-bold text-red-500">${product.price}</span>
                  <span className="text-gray-400 line-through">${product.originalPrice}</span>
                </div>
                <div className="mb-3">
                  <div className="flex justify-between text-sm mb-1">
                    <span className="text-gray-600">Sold: {product.sold}</span>
                    <span className="text-gray-600">Available: {product.total - product.sold}</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div 
                      className="bg-red-500 h-2 rounded-full transition-all"
                      style={{ width: `${(product.sold / product.total) * 100}%` }}
                    />
                  </div>
                </div>
                <button className="w-full bg-red-500 hover:bg-red-600 text-white font-semibold py-3 rounded-lg transition flex items-center justify-center gap-2">
                  <ShoppingBag className="w-5 h-5" />
                  Add to Cart
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}

export default FlashSaleSection
