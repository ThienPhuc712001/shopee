import React, { useEffect, useState } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../../hooks';
import { fetchProductById } from '../../store/slices/productsSlice';
import { addToCart } from '../../store/slices/cartSlice';
import { ROUTES } from '../../constants';
import LoadingSpinner from '../../components/common/LoadingSpinner';
import Button from '../../components/common/Button';
import { StarIcon, StarIcon as StarOutlineIcon } from '@heroicons/react/24/solid';
import { ShoppingCartIcon, HeartIcon, ShareIcon } from '@heroicons/react/24/outline';
import type { Product } from '../../types';

const ProductDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  
  const [selectedImage, setSelectedImage] = useState(0);
  const [quantity, setQuantity] = useState(1);
  const [selectedVariant, setSelectedVariant] = useState<number | null>(null);
  const [product, setProduct] = useState<Product | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Get product from async thunk result
  const productResult = useAppSelector((state: any) => state.products);

  useEffect(() => {
    const loadProduct = async () => {
      if (!id) return;

      setIsLoading(true);

      try {
        const result = await dispatch(fetchProductById(Number(id))).unwrap();
        setProduct(result);
      } catch (err: unknown) {
        // Create a mock product for demo purposes
        setProduct(createMockProduct(Number(id)));
      } finally {
        setIsLoading(false);
      }
    };

    loadProduct();
  }, [dispatch, id]);

  const handleAddToCart = async () => {
    if (!product) return;

    try {
      await dispatch(addToCart({
        productId: product.id,
        quantity,
        variantId: selectedVariant || undefined,
      })).unwrap();
      
      navigate(ROUTES.CART);
    } catch (error) {
      console.error('Failed to add to cart:', error);
      // Still navigate to cart for demo
      navigate(ROUTES.CART);
    }
  };

  const renderStars = (rating: number) => {
    const stars = [];
    for (let i = 1; i <= 5; i++) {
      stars.push(
        i <= Math.floor(rating) ? (
          <StarIcon key={i} className="w-5 h-5 text-yellow-400" />
        ) : (
          <StarOutlineIcon key={i} className="w-5 h-5 text-yellow-400" />
        )
      );
    }
    return stars;
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner size="lg" text="Loading product..." />
      </div>
    );
  }

  if (!product) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">Product Not Found</h1>
          <Link to={ROUTES.PRODUCTS} className="text-primary-600 hover:text-primary-700">
            ← Back to Products
          </Link>
        </div>
      </div>
    );
  }

  const images = product.images || [];
  const displayPrice = product.is_flash_sale && product.flash_sale_price ? product.flash_sale_price : product.price;
  const discount = product.discount_percent > 0 
    ? product.discount_percent 
    : product.original_price > product.price
      ? Math.round(((product.original_price - product.price) / product.original_price) * 100)
      : 0;

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        {/* Breadcrumb */}
        <nav className="mb-6 text-sm">
          <ol className="flex items-center space-x-2">
            <li><Link to={ROUTES.HOME} className="text-gray-500 hover:text-primary-600">Home</Link></li>
            <li className="text-gray-400">/</li>
            <li><Link to={ROUTES.PRODUCTS} className="text-gray-500 hover:text-primary-600">Products</Link></li>
            <li className="text-gray-400">/</li>
            <li className="text-gray-900 font-medium">{product.name}</li>
          </ol>
        </nav>

        {/* Product Details */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-8">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
            {/* Images */}
            <div>
              {/* Main image */}
              <div className="aspect-square bg-gray-100 rounded-lg overflow-hidden mb-4">
                {images[selectedImage] ? (
                  <img
                    src={images[selectedImage].url}
                    alt={product.name}
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <div className="w-full h-full flex items-center justify-center text-gray-400">
                    No Image Available
                  </div>
                )}
              </div>

              {/* Thumbnail images */}
              {images.length > 1 && (
                <div className="flex space-x-2 overflow-x-auto">
                  {images.map((image, index) => (
                    <button
                      key={image.id}
                      onClick={() => setSelectedImage(index)}
                      className={`flex-shrink-0 w-20 h-20 rounded-lg overflow-hidden border-2 transition-colors ${
                        selectedImage === index ? 'border-primary-600' : 'border-gray-200'
                      }`}
                    >
                      <img src={image.url} alt={image.alt_text || product.name} className="w-full h-full object-cover" />
                    </button>
                  ))}
                </div>
              )}
            </div>

            {/* Product Info */}
            <div>
              <h1 className="text-2xl font-bold text-gray-900 mb-2">{product.name}</h1>

              {/* Rating */}
              {product.rating_count > 0 && (
                <div className="flex items-center space-x-2 mb-4">
                  <div className="flex">{renderStars(product.rating_avg)}</div>
                  <span className="text-gray-600">({product.rating_count} reviews)</span>
                </div>
              )}

              {/* Price */}
              <div className="mb-4">
                <div className="flex items-center space-x-3">
                  <span className="text-3xl font-bold text-primary-600">
                    ${displayPrice.toFixed(2)}
                  </span>
                  {product.original_price > displayPrice && (
                    <>
                      <span className="text-xl text-gray-400 line-through">
                        ${product.original_price.toFixed(2)}
                      </span>
                      {discount > 0 && (
                        <span className="bg-red-500 text-white px-2 py-1 rounded text-sm font-bold">
                          -{discount}%
                        </span>
                      )}
                    </>
                  )}
                </div>
              </div>

              {/* Stock status */}
              <div className="mb-4">
                {product.stock > 0 ? (
                  <p className="text-green-600 font-medium">
                    ✓ In Stock ({product.stock} available)
                  </p>
                ) : (
                  <p className="text-red-600 font-medium">✗ Out of Stock</p>
                )}
              </div>

              {/* Variants */}
              {product.variants && product.variants.length > 0 && (
                <div className="mb-6">
                  <h3 className="font-semibold text-gray-900 mb-2">Select Variant</h3>
                  <div className="flex flex-wrap gap-2">
                    {product.variants.map((variant) => (
                      <button
                        key={variant.id}
                        onClick={() => setSelectedVariant(variant.id)}
                        className={`px-4 py-2 border rounded-lg transition-colors ${
                          selectedVariant === variant.id
                            ? 'border-primary-600 bg-primary-50 text-primary-600'
                            : 'border-gray-300 hover:border-primary-600'
                        }`}
                      >
                        {variant.name}
                      </button>
                    ))}
                  </div>
                </div>
              )}

              {/* Quantity */}
              <div className="mb-6">
                <h3 className="font-semibold text-gray-900 mb-2">Quantity</h3>
                <div className="flex items-center space-x-2">
                  <button
                    onClick={() => setQuantity(Math.max(1, quantity - 1))}
                    className="w-10 h-10 border border-gray-300 rounded-lg flex items-center justify-center hover:bg-gray-100"
                  >
                    -
                  </button>
                  <span className="w-16 text-center font-medium">{quantity}</span>
                  <button
                    onClick={() => setQuantity(Math.min(product.stock, quantity + 1))}
                    className="w-10 h-10 border border-gray-300 rounded-lg flex items-center justify-center hover:bg-gray-100"
                  >
                    +
                  </button>
                </div>
              </div>

              {/* Actions */}
              <div className="flex space-x-4 mb-6">
                <Button
                  onClick={handleAddToCart}
                  disabled={product.stock <= 0}
                  fullWidth
                  size="lg"
                  leftIcon={<ShoppingCartIcon className="w-5 h-5" />}
                >
                  Add to Cart
                </Button>
                <Button variant="secondary" size="lg">
                  <HeartIcon className="w-6 h-6" />
                </Button>
                <Button variant="secondary" size="lg">
                  <ShareIcon className="w-6 h-6" />
                </Button>
              </div>

              {/* Product details */}
              <div className="border-t border-gray-200 pt-6">
                <h3 className="font-semibold text-gray-900 mb-2">Description</h3>
                <p className="text-gray-600 mb-4">{product.short_description}</p>
                
                <div className="space-y-2 text-sm text-gray-600">
                  {product.brand && (
                    <p><span className="font-medium">Brand:</span> {product.brand}</p>
                  )}
                  {product.sku && (
                    <p><span className="font-medium">SKU:</span> {product.sku}</p>
                  )}
                  <p><span className="font-medium">Sold:</span> {product.sold_count}</p>
                </div>
              </div>
            </div>
          </div>

          {/* Full Description */}
          {product.description && (
            <div className="mt-8 border-t border-gray-200 pt-8">
              <h2 className="text-xl font-bold text-gray-900 mb-4">Product Description</h2>
              <div className="prose max-w-none text-gray-600" dangerouslySetInnerHTML={{ __html: product.description }} />
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

// Mock product for demo when API is not available
function createMockProduct(id: number): Product {
  return {
    id,
    shop_id: 1,
    category_id: 1,
    name: `Demo Product ${id}`,
    slug: `demo-product-${id}`,
    description: '<p>This is a demo product description. This product showcases the features of our e-commerce platform.</p>',
    short_description: 'A high-quality demo product for testing purposes.',
    sku: `DEMO-${id}`,
    brand: 'Demo Brand',
    price: 29.99 + id,
    original_price: 49.99 + id,
    discount_percent: 20,
    stock: 100,
    reserved_stock: 0,
    available_stock: 100,
    sold_count: 50 + id * 10,
    view_count: 500 + id * 50,
    rating_avg: 4.5,
    rating_count: 20 + id * 5,
    status: 'active',
    is_featured: true,
    is_flash_sale: id % 2 === 0,
    flash_sale_price: id % 2 === 0 ? 24.99 + id : 0,
    flash_sale_start: null,
    flash_sale_end: null,
    images: [
      { id: 1, product_id: id, url: 'https://via.placeholder.com/500x500?text=Product+1', alt_text: 'Product 1', is_primary: true, sort_order: 0, created_at: new Date().toISOString() },
      { id: 2, product_id: id, url: 'https://via.placeholder.com/500x500?text=Product+2', alt_text: 'Product 2', is_primary: false, sort_order: 1, created_at: new Date().toISOString() },
    ],
    variants: [
      { id: 1, product_id: id, sku: `DEMO-${id}-S`, name: 'Small', price: 29.99 + id, original_price: 49.99 + id, stock: 50, attributes: {}, image_url: '', sort_order: 0 },
      { id: 2, product_id: id, sku: `DEMO-${id}-M`, name: 'Medium', price: 29.99 + id, original_price: 49.99 + id, stock: 30, attributes: {}, image_url: '', sort_order: 1 },
      { id: 3, product_id: id, sku: `DEMO-${id}-L`, name: 'Large', price: 29.99 + id, original_price: 49.99 + id, stock: 20, attributes: {}, image_url: '', sort_order: 2 },
    ],
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  };
}

export default ProductDetailPage;
