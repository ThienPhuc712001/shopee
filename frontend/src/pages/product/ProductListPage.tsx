import { useEffect, useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { useNavigate } from 'react-router-dom'
import { AppDispatch, RootState } from '@store/store'
import { fetchCartStart, fetchCartSuccess, addToCartStart, addToCartSuccess, addToCartFailure } from '@store/slices/cartSlice'
import cartService from '@services/cartService'
import productService from '@services/productService'
import ProductCard from '@components/product/ProductCard'

function ProductListPage() {
  const dispatch = useDispatch<AppDispatch>()
  const navigate = useNavigate()
  const { isAuthenticated } = useSelector((state: RootState) => state.auth)
  const { cart } = useSelector((state: RootState) => state.cart)
  const [products, setProducts] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [categories, setCategories] = useState<any[]>([])

  useEffect(() => {
    fetchProducts()
    fetchCategories()
    if (isAuthenticated) {
      fetchCart()
    }
  }, [isAuthenticated])

  const fetchCart = async () => {
    try {
      dispatch(fetchCartStart())
      const data = await cartService.getCart()
      dispatch(fetchCartSuccess(data))
    } catch (err: any) {
      console.error('Error fetching cart:', err)
    }
  }

  const handleAddToCart = async (product: any) => {
    if (!isAuthenticated) {
      alert('Please login to add items to cart')
      navigate('/login')
      return
    }

    try {
      dispatch(addToCartStart())
      const data = await cartService.addToCart(product.id, 1)
      dispatch(addToCartSuccess(data))
      alert(`✅ Added ${product.name} to cart!`)
    } catch (err: any) {
      console.error('Error adding to cart:', err)
      dispatch(addToCartFailure(err.response?.data?.message || 'Failed to add to cart'))
      alert('❌ Failed to add to cart')
    }
  }

  const fetchProducts = async () => {
    try {
      setLoading(true)
      const data = await productService.getProducts({ limit: 8 })
      // Handle both array and object with data property
      const productsData = Array.isArray(data) ? data : (data.data || [])
      setProducts(productsData)
      setError(null)
    } catch (err: any) {
      console.error('Error fetching products:', err)
      setError(err.response?.data?.message || 'Failed to load products')
      // Use empty array as fallback
      setProducts([])
    } finally {
      setLoading(false)
    }
  }

  const fetchCategories = async () => {
    try {
      const data = await productService.getCategories()
      // Handle both array and object with data property
      const categoriesData = Array.isArray(data) ? data : (data.data || [])
      setCategories(categoriesData)
    } catch (err: any) {
      console.error('Error fetching categories:', err)
      // Use empty array as fallback
      setCategories([])
    }
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">
          <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-primary-500 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading products...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex flex-col md:flex-row gap-8">
        {/* Filters Sidebar */}
        <aside className="w-full md:w-64 flex-shrink-0">
          <div className="bg-white rounded-lg shadow-md p-6 sticky top-24">
            <h3 className="font-bold text-lg mb-4">Filters</h3>
            {/* Category Filter */}
            <div className="mb-6">
              <h4 className="font-semibold mb-2">Category</h4>
              <div className="space-y-2">
                {categories.map((category) => (
                  <label key={category.id} className="flex items-center gap-2">
                    <input type="checkbox" className="rounded" />
                    <span>{category.name}</span>
                  </label>
                ))}
              </div>
            </div>

            {/* Price Range */}
            <div className="mb-6">
              <h4 className="font-semibold mb-2">Price Range</h4>
              <div className="flex gap-2">
                <input
                  type="number"
                  placeholder="Min"
                  className="w-full px-3 py-2 border rounded-lg text-sm"
                />
                <input
                  type="number"
                  placeholder="Max"
                  className="w-full px-3 py-2 border rounded-lg text-sm"
                />
              </div>
            </div>

            <button 
              onClick={fetchProducts}
              className="w-full bg-primary-500 hover:bg-primary-600 text-white font-semibold py-3 rounded-lg transition"
            >
              Apply Filters
            </button>
          </div>
        </aside>

        {/* Products Grid */}
        <div className="flex-1">
          <div className="flex justify-between items-center mb-6">
            <p className="text-gray-600">
              {products.length > 0 
                ? `${products.length} products found` 
                : 'No products available (Backend may not have data)'}
            </p>
            <select className="px-4 py-2 border rounded-lg outline-none focus:ring-2 focus:ring-primary-500">
              <option>Sort by: Featured</option>
              <option>Price: Low to High</option>
              <option>Price: High to Low</option>
              <option>Newest First</option>
              <option>Best Selling</option>
            </select>
          </div>

          {error && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4">
              <p className="font-semibold">API Error:</p>
              <p>{error}</p>
              <p className="text-sm mt-2">Note: This is expected if backend has no data in database</p>
            </div>
          )}

          {products.length > 0 ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              {products.map((product) => (
                <ProductCard 
                  key={product.id} 
                  product={product}
                  onAddToCart={handleAddToCart}
                />
              ))}
            </div>
          ) : (
            <div className="bg-white rounded-lg shadow-md p-12 text-center">
              <div className="text-6xl mb-4">📦</div>
              <h3 className="text-xl font-bold text-gray-900 mb-2">No Products Yet</h3>
              <p className="text-gray-600 mb-6">
                The backend API is connected, but there are no products in the database.
              </p>
              <div className="bg-gray-50 p-4 rounded-lg text-left">
                <p className="font-semibold mb-2">API Status:</p>
                <ul className="text-sm text-gray-600 space-y-1">
                  <li>✅ Backend: http://localhost:8080</li>
                  <li>✅ Frontend: http://localhost:3000</li>
                  <li>✅ API Proxy: Configured</li>
                  <li>⚠️ Database: No products found</li>
                </ul>
              </div>
            </div>
          )}

          {/* Pagination */}
          {products.length > 0 && (
            <div className="flex justify-center mt-8 gap-2">
              <button className="px-4 py-2 border rounded-lg hover:bg-gray-50 transition">Previous</button>
              <button className="px-4 py-2 bg-primary-500 text-white rounded-lg">1</button>
              <button className="px-4 py-2 border rounded-lg hover:bg-gray-50 transition">2</button>
              <button className="px-4 py-2 border rounded-lg hover:bg-gray-50 transition">3</button>
              <button className="px-4 py-2 border rounded-lg hover:bg-gray-50 transition">Next</button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default ProductListPage
