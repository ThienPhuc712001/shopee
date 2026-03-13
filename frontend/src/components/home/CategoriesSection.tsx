import { useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useDispatch, useSelector } from 'react-redux'
import { AppDispatch, RootState } from '@store/store'
import { fetchCategoriesSuccess } from '@store/slices/productSlice'
import { Smartphone, Shirt, Home, Book, Trophy, Watch } from 'lucide-react'

const categoryIcons: Record<string, any> = {
  'Electronics': Smartphone,
  'Fashion': Shirt,
  'Home': Home,
  'Books': Book,
  'Sports': Trophy,
  'Watches': Watch,
}

function CategoriesSection() {
  const dispatch = useDispatch<AppDispatch>()
  const { categories } = useSelector((state: RootState) => state.products)

  useEffect(() => {
    // Mock categories data - will fetch from API
    const mockCategories = [
      { id: 1, name: 'Electronics', icon: '📱', color: 'from-blue-500 to-blue-600' },
      { id: 2, name: 'Fashion', icon: '👕', color: 'from-pink-500 to-pink-600' },
      { id: 3, name: 'Home & Living', icon: '🏠', color: 'from-green-500 to-green-600' },
      { id: 4, name: 'Books', icon: '📚', color: 'from-yellow-500 to-yellow-600' },
      { id: 5, name: 'Sports', icon: '⚽', color: 'from-red-500 to-red-600' },
      { id: 6, name: 'Watches', icon: '⌚', color: 'from-purple-500 to-purple-600' },
      { id: 7, name: 'Beauty', icon: '💄', color: 'from-rose-500 to-rose-600' },
      { id: 8, name: 'Toys', icon: '🎮', color: 'from-orange-500 to-orange-600' },
    ]
    dispatch(fetchCategoriesSuccess(mockCategories as any))
  }, [dispatch])

  return (
    <section className="py-12 bg-white">
      <div className="container mx-auto px-4">
        <div className="text-center mb-8">
          <h2 className="text-3xl font-bold text-gray-900 mb-2">Shop by Category</h2>
          <p className="text-gray-600">Browse through our wide range of categories</p>
        </div>

        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-8 gap-4">
          {categories.map((category: any) => {
            const IconComponent = categoryIcons[category.name] || null
            return (
              <Link
                key={category.id}
                to={`/category/${category.id}`}
                className="group flex flex-col items-center p-6 bg-gray-50 rounded-xl hover:bg-gradient-to-br hover:from-primary-500 hover:to-primary-600 transition-all duration-300 hover:scale-105 hover:shadow-lg"
              >
                <div className="w-16 h-16 mb-4 flex items-center justify-center text-4xl group-hover:scale-110 transition-transform">
                  {category.icon || <IconComponent className="w-10 h-10 text-white" />}
                </div>
                <span className="text-sm font-medium text-gray-700 group-hover:text-white transition-colors text-center">
                  {category.name}
                </span>
              </Link>
            )
          })}
        </div>
      </div>
    </section>
  )
}

export default CategoriesSection
