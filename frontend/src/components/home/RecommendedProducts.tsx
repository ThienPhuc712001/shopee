import ProductCard from '@components/product/ProductCard'
import { Link } from 'react-router-dom'

const recommendedProducts = [
  { id: 5, name: 'Bluetooth Speaker', price: 79.99, original_price: 129.99, rating: 4.6, review_count: 198, images: [{url: 'https://images.unsplash.com/photo-1608043152269-423dbba4e7e1?w=300&h=300&fit=crop'}], stock: 45, status: 'active' as const },
  { id: 6, name: 'Mechanical Keyboard', price: 129.99, original_price: 179.99, rating: 4.7, review_count: 267, images: [{url: 'https://images.unsplash.com/photo-1587829741301-dc798b91a603?w=300&h=300&fit=crop'}], stock: 60, status: 'active' as const },
  { id: 7, name: 'USB-C Hub', price: 49.99, original_price: 79.99, rating: 4.5, review_count: 156, images: [{url: 'https://images.unsplash.com/photo-1625841442316-65c6f7a8a0a8?w=300&h=300&fit=crop'}], stock: 100, status: 'active' as const },
  { id: 8, name: 'Wireless Charger', price: 39.99, original_price: 59.99, rating: 4.4, review_count: 234, images: [{url: 'https://images.unsplash.com/photo-1615214589305-791894887577?w=300&h=300&fit=crop'}], stock: 75, status: 'active' as const },
]

function RecommendedProducts() {
  return (
    <section className="py-12 bg-gray-50">
      <div className="container mx-auto px-4">
        <div className="text-center mb-8">
          <h2 className="text-3xl font-bold text-gray-900 mb-2">Recommended For You</h2>
          <p className="text-gray-600">Based on your browsing history and preferences</p>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {recommendedProducts.map((product) => (
            <ProductCard key={product.id} product={product as any} />
          ))}
        </div>
      </div>
    </section>
  )
}

export default RecommendedProducts
