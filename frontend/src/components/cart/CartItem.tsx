import React from 'react';
import { MinusIcon, PlusIcon, TrashIcon } from '@heroicons/react/24/outline';
import type { CartItem as CartItemType } from '../../types';
import { useAppDispatch } from '../../hooks';
import { updateCartItem, removeFromCart } from '../../store/slices/cartSlice';

interface CartItemProps {
  item: CartItemType;
}

const CartItem: React.FC<CartItemProps> = ({ item }) => {
  const dispatch = useAppDispatch();
  const { id, product_name, product_image, price, quantity, subtotal, product } = item;

  const handleIncrease = () => {
    if (quantity < (product?.stock || 999)) {
      dispatch(updateCartItem({ itemId: id, quantity: quantity + 1 }));
    }
  };

  const handleDecrease = () => {
    if (quantity > 1) {
      dispatch(updateCartItem({ itemId: id, quantity: quantity - 1 }));
    }
  };

  const handleRemove = () => {
    if (window.confirm('Are you sure you want to remove this item?')) {
      dispatch(removeFromCart(id));
    }
  };

  return (
    <div className="flex items-center space-x-4 py-4 border-b border-gray-200 last:border-b-0">
      {/* Product image */}
      <div className="w-24 h-24 flex-shrink-0 bg-gray-100 rounded-lg overflow-hidden">
        {product_image ? (
          <img
            src={product_image}
            alt={product_name}
            className="w-full h-full object-cover"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-gray-400">
            No Image
          </div>
        )}
      </div>

      {/* Product info */}
      <div className="flex-1 min-w-0">
        <h3 className="font-medium text-gray-800 truncate">{product_name}</h3>
        <p className="text-sm text-gray-500 mt-1">
          Price: ${price.toFixed(2)}
        </p>
        {product?.stock !== undefined && product.stock < 10 && (
          <p className="text-xs text-orange-500 mt-1">
            Only {product.stock} left in stock
          </p>
        )}
      </div>

      {/* Quantity controls */}
      <div className="flex items-center space-x-2">
        <button
          onClick={handleDecrease}
          disabled={quantity <= 1}
          className="p-1 rounded-md border border-gray-300 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          <MinusIcon className="w-4 h-4" />
        </button>
        <span className="w-12 text-center font-medium">{quantity}</span>
        <button
          onClick={handleIncrease}
          disabled={quantity >= (product?.stock || 999)}
          className="p-1 rounded-md border border-gray-300 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          <PlusIcon className="w-4 h-4" />
        </button>
      </div>

      {/* Subtotal */}
      <div className="text-right w-28">
        <p className="font-bold text-gray-800">${subtotal.toFixed(2)}</p>
      </div>

      {/* Remove button */}
      <button
        onClick={handleRemove}
        className="p-2 text-red-500 hover:bg-red-50 rounded-lg transition-colors"
        title="Remove item"
      >
        <TrashIcon className="w-5 h-5" />
      </button>
    </div>
  );
};

export default CartItem;
