import React from 'react';
import { Link } from 'react-router-dom';
import { EyeIcon } from '@heroicons/react/24/outline';
import type { Order } from '../../types';
import { ROUTES } from '../../constants';

interface OrderCardProps {
  order: Order;
}

const OrderCard: React.FC<OrderCardProps> = ({ order }) => {
  const {
    id,
    order_number,
    status,
    payment_status,
    total_amount,
    items,
    created_at,
    shipped_at,
    delivered_at,
  } = order;

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      pending: 'bg-yellow-100 text-yellow-800',
      paid: 'bg-green-100 text-green-800',
      processing: 'bg-blue-100 text-blue-800',
      shipped: 'bg-purple-100 text-purple-800',
      delivered: 'bg-green-100 text-green-800',
      cancelled: 'bg-red-100 text-red-800',
      refunded: 'bg-gray-100 text-gray-800',
    };
    return colors[status] || 'bg-gray-100 text-gray-800';
  };

  const getStatusLabel = (status: string) => {
    const labels: Record<string, string> = {
      pending: 'Pending',
      paid: 'Paid',
      processing: 'Processing',
      shipped: 'Shipped',
      delivered: 'Delivered',
      cancelled: 'Cancelled',
      refunded: 'Refunded',
    };
    return labels[status] || status;
  };

  const totalItems = items?.reduce((sum, item) => sum + item.quantity, 0) || 0;
  const primaryImage = items?.[0]?.product_image;

  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow">
      {/* Order header */}
      <div className="bg-gray-50 px-4 py-3 border-b border-gray-200">
        <div className="flex justify-between items-center">
          <div>
            <p className="text-sm text-gray-500">Order Number</p>
            <p className="font-semibold text-gray-800">{order_number}</p>
          </div>
          <div className="text-right">
            <p className="text-sm text-gray-500">Order Date</p>
            <p className="font-medium text-gray-800">
              {new Date(created_at).toLocaleDateString()}
            </p>
          </div>
        </div>
      </div>

      {/* Order body */}
      <div className="p-4">
        <div className="flex items-start space-x-4">
          {/* Product image */}
          <div className="w-20 h-20 flex-shrink-0 bg-gray-100 rounded-lg overflow-hidden">
            {primaryImage ? (
              <img
                src={primaryImage}
                alt={items?.[0]?.product_name}
                className="w-full h-full object-cover"
              />
            ) : (
              <div className="w-full h-full flex items-center justify-center text-gray-400 text-xs">
                No Image
              </div>
            )}
          </div>

          {/* Order info */}
          <div className="flex-1">
            <div className="flex justify-between items-start mb-2">
              <div>
                <p className="font-medium text-gray-800">
                  {items?.[0]?.product_name || 'Multiple Items'}
                </p>
                {items && items.length > 1 && (
                  <p className="text-sm text-gray-500">
                    +{items.length - 1} more item{items.length - 1 > 1 ? 's' : ''}
                  </p>
                )}
              </div>
              <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(status)}`}>
                {getStatusLabel(status)}
              </span>
            </div>

            <div className="flex justify-between items-center">
              <div className="text-sm text-gray-500">
                <span>{totalItems} item{totalItems > 1 ? 's' : ''}</span>
                {payment_status === 'paid' && (
                  <span className="ml-2 text-green-600">• Paid</span>
                )}
              </div>
              <p className="font-bold text-gray-800">${total_amount.toFixed(2)}</p>
            </div>
          </div>
        </div>

        {/* Order actions */}
        <div className="mt-4 pt-4 border-t border-gray-200 flex justify-end space-x-2">
          <Link
            to={ROUTES.USER_ORDER_DETAIL.replace(':id', id.toString())}
            className="inline-flex items-center px-4 py-2 text-sm font-medium text-primary-600 border border-primary-600 rounded-lg hover:bg-primary-50 transition-colors"
          >
            <EyeIcon className="w-4 h-4 mr-2" />
            View Details
          </Link>
          {status === 'pending' && (
            <button className="px-4 py-2 text-sm font-medium text-red-600 border border-red-600 rounded-lg hover:bg-red-50 transition-colors">
              Cancel Order
            </button>
          )}
          {status === 'delivered' && (
            <button className="px-4 py-2 text-sm font-medium text-white bg-primary-600 rounded-lg hover:bg-primary-700 transition-colors">
              Write Review
            </button>
          )}
        </div>
      </div>

      {/* Order footer - tracking info */}
      {(status === 'shipped' || status === 'delivered') && (
        <div className="bg-gray-50 px-4 py-3 border-t border-gray-200">
          <div className="flex justify-between items-center text-sm">
            <span className="text-gray-500">
              {status === 'shipped' ? 'Shipped on:' : 'Delivered on:'}
            </span>
            <span className="font-medium text-gray-800">
              {status === 'shipped' 
                ? shipped_at ? new Date(shipped_at).toLocaleDateString() : 'N/A'
                : delivered_at ? new Date(delivered_at).toLocaleDateString() : 'N/A'
              }
            </span>
          </div>
        </div>
      )}
    </div>
  );
};

export default OrderCard;
