import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../hooks';
import { fetchFeaturedProducts, fetchFlashSaleProducts, fetchCategories } from '../store/slices/productsSlice';
import { ROUTES } from '../constants';
import ProductCard from '../components/product/ProductCard';
import LoadingSpinner from '../components/common/LoadingSpinner';
import { ChevronRightIcon } from '@heroicons/react/24/outline';
import type { Product, Category } from '../types';

const HomePage: React.FC = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const { featuredProducts, flashSaleProducts, categories, isLoading } = useAppSelector((state: any) => state.products);

  const [demoProducts, setDemoProducts] = useState<Product[]>([]);
  const [demoCategories, setDemoCategories] = useState<Category[]>([]);

  useEffect(() => {
    dispatch(fetchFeaturedProducts(8));
    dispatch(fetchFlashSaleProducts());
    dispatch(fetchCategories());
  }, [dispatch]);

  useEffect(() => {
    if (!featuredProducts.length) {
      setDemoProducts(generateDemoProducts(8));
    }
    if (!categories.length) {
      setDemoCategories(generateDemoCategories());
    }
  }, [featuredProducts, categories]);

  const handleAddToCart = (product: Product) => {
    console.log('Added to cart:', product);
  };

  const displayCategories = demoCategories.length > 0 ? demoCategories : categories;
  const displayFlashProducts = flashSaleProducts.length > 0 ? flashSaleProducts : demoProducts.slice(0, 4);
  const displayProducts = featuredProducts.length > 0 ? featuredProducts : demoProducts;

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Hero Section */}
      <section className="bg-gradient-to-r from-primary-600 to-primary-700 text-white">
        <div className="container mx-auto px-4 py-20">
          <div className="text-center">
            <h1 className="text-5xl font-bold mb-4">Welcome to Our Store</h1>
            <p className="text-xl mb-8">Discover amazing products at great prices</p>
            <Link
              to={ROUTES.PRODUCTS}
              className="inline-block bg-white text-primary-600 px-8 py-3 rounded-lg font-semibold hover:bg-gray-100 transition-colors"
            >
              Shop Now
            </Link>
          </div>
        </div>
      </section>

      {/* Categories Section */}
      <section className="container mx-auto px-4 py-12">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-gray-900">Shop by Category</h2>
          <Link to={ROUTES.PRODUCTS} className="text-primary-600 hover:text-primary-700 flex items-center">
            View all <ChevronRightIcon className="w-4 h-4 ml-1" />
          </Link>
        </div>

        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
          {displayCategories.map((category: any, index: number) => (
            <Link
              key={category.id || index}
              to={`${ROUTES.PRODUCTS}?category_id=${category.id}`}
              className="group"
            >
              <div className="aspect-square bg-gradient-to-br from-primary-100 to-primary-200 rounded-lg overflow-hidden mb-2 flex items-center justify-center">
                {category.image_url ? (
                  <img
                    src={category.image_url}
                    alt={category.name}
                    className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                  />
                ) : (
                  <span className="text-4xl">📦</span>
                )}
              </div>
              <p className="text-sm text-gray-700 text-center group-hover:text-primary-600">{category.name}</p>
            </Link>
          ))}
        </div>
      </section>

      {/* Flash Sale Section */}
      {displayFlashProducts.length > 0 && (
        <section className="container mx-auto px-4 py-12">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold text-gray-900">Flash Sale</h2>
            <Link to={ROUTES.PRODUCTS} className="text-primary-600 hover:text-primary-700 flex items-center">
              View all <ChevronRightIcon className="w-4 h-4 ml-1" />
            </Link>
          </div>

          <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
            {displayFlashProducts.map((product: any) => (
              <ProductCard key={product.id} product={product} onAddToCart={handleAddToCart} />
            ))}
          </div>
        </section>
      )}

      {/* Featured Products Section */}
      <section className="container mx-auto px-4 py-12">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-gray-900">Featured Products</h2>
          <Link to={ROUTES.PRODUCTS} className="text-primary-600 hover:text-primary-700 flex items-center">
            View all <ChevronRightIcon className="w-4 h-4 ml-1" />
          </Link>
        </div>

        {isLoading && !displayProducts.length ? (
          <div className="flex justify-center py-12">
            <LoadingSpinner size="lg" />
          </div>
        ) : (
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
            {displayProducts.map((product: any) => (
              <ProductCard key={product.id} product={product} onAddToCart={handleAddToCart} />
            ))}
          </div>
        )}
      </section>
    </div>
  );
};

function generateDemoProducts(count: number): Product[] {
  return Array.from({ length: count }, (_, i) => ({
    id: i + 1,
    name: `Product ${i + 1}`,
    slug: `product-${i + 1}`,
    price: 9.99 + i * 5,
    original_price: 14.99 + i * 5,
    discount_percent: 20,
    image_url: `https://via.placeholder.com/300x300?text=Product+${i + 1}`,
    rating_avg: 4.5,
    rating_count: 10 + i * 2,
    sold_count: 50 + i * 10,
    status: 'active' as const,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  }));
}

function generateDemoCategories(): Category[] {
  return [
    { id: 1, parent_id: null, name: 'Electronics', slug: 'electronics', description: 'Electronic devices', icon_url: '', image_url: '', level: 0, sort_order: 1, is_active: true, created_at: '', updated_at: '' },
    { id: 2, parent_id: null, name: 'Fashion', slug: 'fashion', description: 'Clothing and accessories', icon_url: '', image_url: '', level: 0, sort_order: 2, is_active: true, created_at: '', updated_at: '' },
    { id: 3, parent_id: null, name: 'Home & Living', slug: 'home-living', description: 'Home decor', icon_url: '', image_url: '', level: 0, sort_order: 3, is_active: true, created_at: '', updated_at: '' },
    { id: 4, parent_id: null, name: 'Books', slug: 'books', description: 'Books and publications', icon_url: '', image_url: '', level: 0, sort_order: 4, is_active: true, created_at: '', updated_at: '' },
    { id: 5, parent_id: null, name: 'Sports', slug: 'sports', description: 'Sports equipment', icon_url: '', image_url: '', level: 0, sort_order: 5, is_active: true, created_at: '', updated_at: '' },
  ];
}

export default HomePage;
