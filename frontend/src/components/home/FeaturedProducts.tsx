import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import ProductCard from '@components/product/ProductCard'
import productService from '@services/productService'

function FeaturedProducts() {
  const [products, setProducts] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [apiStatus, setApiStatus] = useState<'checking' | 'connected' | 'error'>('checking')

  useEffect(() => {
    fetchFeaturedProducts()
  }, [])

  const fetchFeaturedProducts = async () => {
    try {
      setLoading(true)
      setApiStatus('checking')
      
      // Try to fetch from API
      const response = await productService.getFeaturedProducts()
      // API returns {data: [...]} or directly array depending on endpoint
      const productsData = Array.isArray(response) ? response : (response.data || [])
      setProducts(productsData)
      setApiStatus('connected')
    } catch (err: any) {
      console.error('Error fetching featured products:', err)
      setApiStatus('error')
      
      // Use mock data for demo if API fails
      const mockProducts = [
        { 
          id: 1, 
          name: 'Premium Wireless Headphones', 
          price: 199.99, 
          original_price: 299.99, 
          rating: 4.8, 
          review_count: 256, 
          images: [{url: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=300&h=300&fit=crop'}], 
          stock: 50, 
          status: 'active' as const 
        },
        { 
          id: 2, 
          name: 'Smart Watch Pro', 
          price: 349.99, 
          original_price: 449.99, 
          rating: 4.9, 
          review_count: 189, 
          images: [{url: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=300&h=300&fit=crop'}], 
          stock: 30, 
          status: 'active' as const 
        },
        { 
          id: 3, 
          name: 'Professional Camera', 
          price: 899.99, 
          original_price: 1199.99, 
          rating: 4.7, 
          review_count: 142, 
          images: [{url: 'https://images.unsplash.com/photo-1526170375885-4d8ecf77b99f?w=300&h=300&fit=crop'}], 
          stock: 15, 
          status: 'active' as const 
        },
        { 
          id: 4, 
          name: 'Gaming Laptop', 
          price: 1299.99, 
          original_price: 1599.99, 
          rating: 4.8, 
          review_count: 324, 
          images: [{url: 'https://images.unsplash.com/photo-1603302576837-37561b2e2302?w=300&h=300&fit=crop'}], 
          stock: 8, 
          status: 'active' as const 
        },
      ]
      setProducts(mockProducts)
      setApiStatus('error')
    } finally {
      setLoading(false)
    }
  }

  return (
    <section className="py-12 bg-white">
      <div className="container mx-auto px-4">
        <div className="text-center mb-8">
          <h2 className="text-3xl font-bold text-gray-900 mb-2">Featured Products</h2>
          <p className="text-gray-600">Handpicked selection of our best products</p>
          
          {/* API Status Indicator */}
          <div className="mt-4 inline-flex items-center gap-2 px-4 py-2 rounded-full bg-gray-100">
            <div className={`w-2 h-2 rounded-full ${
              apiStatus === 'connected' ? 'bg-green-500' :
              apiStatus === 'error' ? 'bg-yellow-500' : 'bg-gray-400'
            }`} />
            <span className="text-sm text-gray-600">
              {apiStatus === 'connected' ? '🟢 Backend Connected' :
               apiStatus === 'error' ? '🟡 Using Mock Data (Backend has no data)' :
               '⚪ Checking API...'}
            </span>
          </div>
        </div>

        {loading ? (
          <div className="text-center py-12">
            <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-primary-500 mx-auto"></div>
            <p className="mt-4 text-gray-600">Loading products...</p>
          </div>
        ) : (
          <>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
              {products.map((product) => (
                <ProductCard key={product.id} product={product} />
              ))}
            </div>

            <div className="text-center mt-8">
              <Link
                to="/products"
                className="inline-block bg-secondary-800 hover:bg-secondary-700 text-white font-semibold px-8 py-4 rounded-lg transition-all hover:scale-105"
              >
                View All Products →
              </Link>
            </div>
          </>
        )}
      </div>
    </section>
  )
}

export default FeaturedProducts
