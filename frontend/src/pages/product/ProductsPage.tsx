import React, { useEffect, useState } from 'react';
import { useSearchParams, Link } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../../hooks';
import { fetchProducts, fetchCategories } from '../../store/slices/productsSlice';
import { ROUTES, SORT_OPTIONS } from '../../constants';
import ProductCard from '../../components/product/ProductCard';
import Pagination from '../../components/common/Pagination';
import LoadingSpinner from '../../components/common/LoadingSpinner';
import { FunnelIcon, ArrowsUpDownIcon } from '@heroicons/react/24/outline';
import type { Product, Category } from '../../types';

const ProductsPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const [searchParams, setSearchParams] = useSearchParams();
  const [showFilters, setShowFilters] = useState(false);
  const [apiError, setApiError] = useState<string | null>(null);
  
  const { products, categories, pagination, isLoading } = useAppSelector((state) => state.products);
  
  // Local state for demo data
  const [demoProducts, setDemoProducts] = useState<Product[]>([]);
  const [demoCategories, setDemoCategories] = useState<Category[]>([]);

  const [localFilters, setLocalFilters] = useState({
    min_price: '',
    max_price: '',
    rating: '',
  });

  useEffect(() => {
    const loadData = async () => {
      setApiError(null);
      
      try {
        const params: Record<string, string> = {};
        
        if (searchParams.get('category_id')) params.category_id = searchParams.get('category_id')!;
        if (searchParams.get('search')) params.search = searchParams.get('search')!;
        if (searchParams.get('sort')) params.sort = searchParams.get('sort')!;
        if (searchParams.get('min_price')) params.min_price = searchParams.get('min_price')!;
        if (searchParams.get('max_price')) params.max_price = searchParams.get('max_price')!;
        if (searchParams.get('page')) params.page = searchParams.get('page')!;
        
        params.limit = '20';
        
        await dispatch(fetchProducts(params as any)).unwrap();
      } catch (error: unknown) {
        const err = error as { message?: string };
        setApiError(err.message || 'Failed to load products');
        // Generate demo products
        setDemoProducts(generateDemoProducts(20));
      }
    };
    
    loadData();
  }, [dispatch, searchParams]);

  useEffect(() => {
    const loadCategories = async () => {
      try {
        await dispatch(fetchCategories()).unwrap();
      } catch {
        setDemoCategories(generateDemoCategories());
      }
    };
    loadCategories();
  }, [dispatch]);

  const handleSortChange = (value: string) => {
    const newParams = new URLSearchParams(searchParams);
    newParams.set('sort', value);
    newParams.delete('page');
    setSearchParams(newParams);
  };

  const handleCategoryFilter = (categoryId: string) => {
    const newParams = new URLSearchParams(searchParams);
    if (categoryId) {
      newParams.set('category_id', categoryId);
    } else {
      newParams.delete('category_id');
    }
    newParams.delete('page');
    setSearchParams(newParams);
  };

  const handlePriceFilter = () => {
    const newParams = new URLSearchParams(searchParams);
    if (localFilters.min_price) {
      newParams.set('min_price', localFilters.min_price);
    } else {
      newParams.delete('min_price');
    }
    if (localFilters.max_price) {
      newParams.set('max_price', localFilters.max_price);
    } else {
      newParams.delete('max_price');
    }
    newParams.delete('page');
    setSearchParams(newParams);
  };

  const handlePageChange = (page: number) => {
    const newParams = new URLSearchParams(searchParams);
    newParams.set('page', page.toString());
    setSearchParams(newParams);
  };

  const handleAddToCart = (productId: number) => {
    console.log('Add to cart:', productId);
  };

  const displayProducts = apiError ? demoProducts : products;
  const displayCategories = categories.length > 0 ? categories : demoCategories;
  const totalProducts = apiError ? demoProducts.length : pagination.total;
  const totalPages = apiError ? Math.ceil(demoProducts.length / 20) : pagination.total_pages;

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">All Products</h1>
          <p className="text-gray-600">
            {totalProducts} products found
          </p>
        </div>

        {apiError && (
          <div className="bg-yellow-50 border border-yellow-200 text-yellow-800 px-4 py-3 rounded-lg mb-6">
            <strong>Demo Mode:</strong> {apiError}. Showing demo products.
          </div>
        )}

        <div className="flex flex-col lg:flex-row gap-6">
          {/* Mobile filter toggle */}
          <button
            onClick={() => setShowFilters(!showFilters)}
            className="lg:hidden flex items-center justify-center space-x-2 px-4 py-2 bg-white border border-gray-300 rounded-lg"
          >
            <FunnelIcon className="w-5 h-5" />
            <span>Filters</span>
          </button>

          {/* Filters Sidebar */}
          <aside className={`lg:w-64 flex-shrink-0 ${showFilters ? 'block' : 'hidden lg:block'}`}>
            <div className="bg-white rounded-lg shadow-md p-4 sticky top-24">
              {/* Category Filter */}
              <div className="mb-6">
                <h3 className="font-semibold text-gray-900 mb-3">Categories</h3>
                <div className="space-y-2">
                  <label className="flex items-center cursor-pointer">
                    <input
                      type="radio"
                      name="category"
                      checked={!searchParams.get('category_id')}
                      onChange={() => handleCategoryFilter('')}
                      className="w-4 h-4 text-primary-600"
                    />
                    <span className="ml-2 text-gray-700">All Categories</span>
                  </label>
                  {displayCategories.map((category) => (
                    <label key={category.id} className="flex items-center cursor-pointer">
                      <input
                        type="radio"
                        name="category"
                        checked={searchParams.get('category_id') === category.id.toString()}
                        onChange={() => handleCategoryFilter(category.id.toString())}
                        className="w-4 h-4 text-primary-600"
                      />
                      <span className="ml-2 text-gray-700">{category.name}</span>
                    </label>
                  ))}
                </div>
              </div>

              {/* Price Filter */}
              <div className="mb-6">
                <h3 className="font-semibold text-gray-900 mb-3">Price Range</h3>
                <div className="space-y-3">
                  <input
                    type="number"
                    placeholder="Min price"
                    value={localFilters.min_price}
                    onChange={(e) => setLocalFilters({ ...localFilters, min_price: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 outline-none"
                  />
                  <input
                    type="number"
                    placeholder="Max price"
                    value={localFilters.max_price}
                    onChange={(e) => setLocalFilters({ ...localFilters, max_price: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 outline-none"
                  />
                  <button
                    onClick={handlePriceFilter}
                    className="w-full btn-primary"
                  >
                    Apply
                  </button>
                </div>
              </div>

              {/* Rating Filter */}
              <div>
                <h3 className="font-semibold text-gray-900 mb-3">Rating</h3>
                <div className="space-y-2">
                  {[4, 3, 2, 1].map((rating) => (
                    <label key={rating} className="flex items-center cursor-pointer">
                      <input
                        type="radio"
                        name="rating"
                        checked={searchParams.get('rating') === rating.toString()}
                        onChange={() => {
                          const newParams = new URLSearchParams(searchParams);
                          if (searchParams.get('rating') === rating.toString()) {
                            newParams.delete('rating');
                          } else {
                            newParams.set('rating', rating.toString());
                          }
                          setSearchParams(newParams);
                        }}
                        className="w-4 h-4 text-primary-600"
                      />
                      <span className="ml-2 text-gray-700">{rating}+ Stars</span>
                    </label>
                  ))}
                </div>
              </div>
            </div>
          </aside>

          {/* Products Grid */}
          <main className="flex-1">
            {/* Sort Bar */}
            <div className="bg-white rounded-lg shadow-md p-4 mb-4 flex justify-between items-center">
              <div className="flex items-center space-x-2 text-gray-600">
                <ArrowsUpDownIcon className="w-5 h-5" />
                <span>Sort by:</span>
              </div>
              <select
                value={searchParams.get('sort') || ''}
                onChange={(e) => handleSortChange(e.target.value)}
                className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 outline-none"
              >
                <option value="">Default</option>
                {SORT_OPTIONS.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
            </div>

            {/* Products */}
            {isLoading && !displayProducts.length ? (
              <div className="flex justify-center py-12">
                <LoadingSpinner size="lg" />
              </div>
            ) : displayProducts.length > 0 ? (
              <>
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-3 gap-4">
                  {displayProducts.map((product) => (
                    <ProductCard key={product.id} product={product} onAddToCart={handleAddToCart} />
                  ))}
                </div>

                {/* Pagination */}
                {totalPages > 1 && (
                  <div className="mt-8">
                    <Pagination
                      currentPage={pagination.page || 1}
                      totalPages={totalPages}
                      onPageChange={handlePageChange}
                      totalItems={totalProducts}
                      itemsPerPage={20}
                    />
                  </div>
                )}
              </>
            ) : (
              <div className="text-center py-12">
                <p className="text-gray-500 text-lg mb-4">No products found</p>
                <Link to={ROUTES.PRODUCTS} className="text-primary-600 hover:text-primary-700">
                  Clear filters and view all products
                </Link>
              </div>
            )}
          </main>
        </div>
      </div>
    </div>
  );
};

// Generate demo products
function generateDemoProducts(count: number): Product[] {
  return Array.from({ length: count }, (_, i) => ({
    id: i + 1,
    shop_id: 1,
    category_id: (i % 5) + 1,
    name: `Demo Product ${i + 1}`,
    slug: `demo-product-${i + 1}`,
    description: `<p>This is a demo product description for Product ${i + 1}.</p>`,
    short_description: `High quality demo product for testing.`,
    sku: `DEMO-${i + 1}`,
    brand: 'Demo Brand',
    price: 19.99 + (i * 5),
    original_price: 39.99 + (i * 5),
    discount_percent: 20 + (i % 10),
    stock: 100 - (i * 5),
    reserved_stock: 0,
    available_stock: 100 - (i * 5),
    sold_count: 50 + (i * 10),
    view_count: 500 + (i * 50),
    rating_avg: 4 + (i % 2) * 0.5,
    rating_count: 20 + (i * 5),
    review_count: 15 + (i * 3),
    status: 'active' as const,
    is_featured: false,
    is_flash_sale: i % 3 === 0,
    flash_sale_price: i % 3 === 0 ? 14.99 + (i * 3) : 0,
    flash_sale_start: null,
    flash_sale_end: null,
    images: [
      {
        id: i + 1,
        product_id: i + 1,
        url: `https://via.placeholder.com/500x500?text=Product+${i + 1}`,
        alt_text: `Product ${i + 1}`,
        is_primary: true,
        sort_order: 0,
        created_at: new Date().toISOString(),
      },
    ],
    variants: [],
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  }));
}

// Generate demo categories
function generateDemoCategories(): Category[] {
  return [
    { id: 1, parent_id: null, name: 'Electronics', slug: 'electronics', description: 'Electronic devices', icon_url: '', image_url: '', level: 0, sort_order: 1, is_active: true, created_at: '', updated_at: '' },
    { id: 2, parent_id: null, name: 'Fashion', slug: 'fashion', description: 'Clothing and accessories', icon_url: '', image_url: '', level: 0, sort_order: 2, is_active: true, created_at: '', updated_at: '' },
    { id: 3, parent_id: null, name: 'Home & Living', slug: 'home-living', description: 'Home decor', icon_url: '', image_url: '', level: 0, sort_order: 3, is_active: true, created_at: '', updated_at: '' },
    { id: 4, parent_id: null, name: 'Books', slug: 'books', description: 'Books and publications', icon_url: '', image_url: '', level: 0, sort_order: 4, is_active: true, created_at: '', updated_at: '' },
    { id: 5, parent_id: null, name: 'Sports', slug: 'sports', description: 'Sports equipment', icon_url: '', image_url: '', level: 0, sort_order: 5, is_active: true, created_at: '', updated_at: '' },
  ];
}

export default ProductsPage;
