import React from 'react';
import { Link } from 'react-router-dom';
import { StarIcon, ShoppingCartIcon } from '@heroicons/react/24/solid';
import { StarIcon as StarOutlineIcon } from '@heroicons/react/24/outline';
import type { Product } from '../../types';
import { ROUTES } from '../../constants';

interface ProductCardProps {
  product: Product;
  onAddToCart?: (productId: number) => void;
}

const ProductCard: React.FC<ProductCardProps> = ({ product, onAddToCart }) => {
  const {
    id,
    name,
    slug,
    price,
    original_price,
    discount_percent,
    images,
    rating_avg,
    rating_count,
    sold_count,
    stock,
    is_flash_sale,
    flash_sale_price,
  } = product;

  const primaryImage = images?.find((img) => img.is_primary) || images?.[0];
  const displayPrice = is_flash_sale && flash_sale_price ? flash_sale_price : price;
  const discount = discount_percent > 0 ? discount_percent : original_price > price 
    ? Math.round(((original_price - price) / original_price) * 100) 
    : 0;
  const isOutOfStock = stock <= 0;

  const renderStars = (rating: number) => {
    const stars = [];
    for (let i = 1; i <= 5; i++) {
      stars.push(
        i <= Math.floor(rating) ? (
          <StarIcon key={i} className="w-4 h-4 text-yellow-400" />
        ) : (
          <StarOutlineIcon key={i} className="w-4 h-4 text-yellow-400" />
        )
      );
    }
    return stars;
  };

  const handleAddToCart = (e: React.MouseEvent) => {
    e.preventDefault();
    if (onAddToCart && !isOutOfStock) {
      onAddToCart(id);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-xl transition-shadow duration-300 group">
      {/* Image */}
      <Link to={ROUTES.PRODUCT_DETAIL.replace(':id', id.toString())} className="block relative">
        <div className="aspect-square overflow-hidden bg-gray-100">
          {primaryImage ? (
            <img
              src={primaryImage.url}
              alt={name}
              className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
              loading="lazy"
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center text-gray-400">
              No Image
            </div>
          )}
        </div>

        {/* Badges */}
        <div className="absolute top-2 left-2 flex flex-col space-y-1">
          {is_flash_sale && (
            <span className="bg-red-500 text-white text-xs font-bold px-2 py-1 rounded">
              FLASH SALE
            </span>
          )}
          {discount > 0 && (
            <span className="bg-orange-500 text-white text-xs font-bold px-2 py-1 rounded">
              -{discount}%
            </span>
          )}
          {isOutOfStock && (
            <span className="bg-gray-800 text-white text-xs font-bold px-2 py-1 rounded">
              OUT OF STOCK
            </span>
          )}
        </div>

        {/* Quick add to cart button */}
        {!isOutOfStock && (
          <button
            onClick={handleAddToCart}
            className="absolute bottom-2 right-2 bg-white p-2 rounded-full shadow-md opacity-0 group-hover:opacity-100 transition-opacity duration-300 hover:bg-primary-600 hover:text-white"
          >
            <ShoppingCartIcon className="w-5 h-5" />
          </button>
        )}
      </Link>

      {/* Content */}
      <div className="p-4">
        <Link to={ROUTES.PRODUCT_DETAIL.replace(':id', id.toString())}>
          <h3 className="font-medium text-gray-800 line-clamp-2 hover:text-primary-600 transition-colors mb-2">
            {name}
          </h3>
        </Link>

        {/* Rating */}
        {rating_count > 0 && (
          <div className="flex items-center space-x-1 mb-2">
            <div className="flex">{renderStars(rating_avg)}</div>
            <span className="text-sm text-gray-500">({rating_count})</span>
          </div>
        )}

        {/* Price */}
        <div className="flex items-center space-x-2 mb-2">
          <span className="text-lg font-bold text-primary-600">
            ${displayPrice.toFixed(2)}
          </span>
          {original_price > displayPrice && (
            <span className="text-sm text-gray-400 line-through">
              ${original_price.toFixed(2)}
            </span>
          )}
        </div>

        {/* Sold count */}
        {sold_count > 0 && (
          <p className="text-sm text-gray-500">{sold_count.toLocaleString()} sold</p>
        )}

        {/* Stock indicator */}
        {!isOutOfStock && stock < 10 && (
          <p className="text-sm text-orange-500 font-medium">Only {stock} left!</p>
        )}
      </div>
    </div>
  );
};

export default ProductCard;
