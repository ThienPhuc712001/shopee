import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../../hooks';
import { fetchUserOrders } from '../../store/slices/ordersSlice';
import { ROUTES } from '../../constants';
import OrderCard from '../../components/order/OrderCard';
import LoadingSpinner from '../../components/common/LoadingSpinner';
import Pagination from '../../components/common/Pagination';

const UserOrdersPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const [currentPage, setCurrentPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState('');
  const { orders, isLoading, pagination } = useAppSelector((state) => state.orders);

  useEffect(() => {
    dispatch(fetchUserOrders({ page: currentPage, limit: 10 }));
  }, [dispatch, currentPage]);

  const filteredOrders = statusFilter
    ? orders.filter((order) => order.status === statusFilter)
    : orders;

  const statusTabs = [
    { value: '', label: 'All Orders' },
    { value: 'pending', label: 'Pending' },
    { value: 'processing', label: 'Processing' },
    { value: 'shipped', label: 'Shipped' },
    { value: 'delivered', label: 'Delivered' },
    { value: 'cancelled', label: 'Cancelled' },
  ];

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">My Orders</h1>

      {/* Status tabs */}
      <div className="bg-white rounded-lg shadow-md mb-6 overflow-hidden">
        <div className="flex overflow-x-auto">
          {statusTabs.map((tab) => (
            <button
              key={tab.value}
              onClick={() => {
                setStatusFilter(tab.value);
                setCurrentPage(1);
              }}
              className={`px-6 py-4 font-medium whitespace-nowrap transition-colors border-b-2 ${
                statusFilter === tab.value
                  ? 'border-primary-600 text-primary-600'
                  : 'border-transparent text-gray-600 hover:text-gray-900'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      {/* Orders list */}
      {isLoading ? (
        <div className="flex justify-center py-12">
          <LoadingSpinner size="lg" />
        </div>
      ) : filteredOrders.length > 0 ? (
        <>
          <div className="space-y-4">
            {filteredOrders.map((order) => (
              <OrderCard key={order.id} order={order} />
            ))}
          </div>

          {/* Pagination */}
          {pagination.total > 10 && (
            <div className="mt-8">
              <Pagination
                currentPage={currentPage}
                totalPages={Math.ceil(pagination.total / 10)}
                onPageChange={setCurrentPage}
                totalItems={pagination.total}
                itemsPerPage={10}
              />
            </div>
          )}
        </>
      ) : (
        <div className="bg-white rounded-lg shadow-md p-12 text-center">
          <ShoppingBagIcon className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-900 mb-2">No orders found</h3>
          <p className="text-gray-600 mb-6">
            {statusFilter ? `You don't have any ${statusFilter} orders.` : "You haven't placed any orders yet."}
          </p>
          <Link to={ROUTES.PRODUCTS}>
            <Button>Start Shopping</Button>
          </Link>
        </div>
      )}
    </div>
  );
};

// Import missing components
import { ShoppingBagIcon } from '@heroicons/react/24/outline';
import Button from '../../components/common/Button';

export default UserOrdersPage;
