import React from 'react';
import { useAppSelector } from '../../hooks';
import { UserIcon } from '@heroicons/react/24/outline';

const UserProfilePage: React.FC = () => {
  const { user } = useAppSelector((state) => state.auth);

  if (!user) return null;

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">My Profile</h1>

      <div className="bg-white rounded-lg shadow-md p-6">
        {/* Profile header */}
        <div className="flex items-center space-x-4 mb-8 pb-8 border-b border-gray-200">
          <div className="w-24 h-24 bg-primary-100 rounded-full flex items-center justify-center">
            {user.avatar ? (
              <img src={user.avatar} alt={user.first_name} className="w-full h-full rounded-full object-cover" />
            ) : (
              <UserIcon className="w-12 h-12 text-primary-600" />
            )}
          </div>
          <div>
            <h2 className="text-xl font-bold text-gray-900">
              {user.first_name} {user.last_name}
            </h2>
            <p className="text-gray-600">{user.email}</p>
            <span className={`inline-block mt-2 px-3 py-1 rounded-full text-xs font-medium ${
              user.role === 'admin' ? 'bg-purple-100 text-purple-800' :
              user.role === 'seller' ? 'bg-blue-100 text-blue-800' :
              'bg-gray-100 text-gray-800'
            }`}>
              {user.role}
            </span>
          </div>
        </div>

        {/* Profile info */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label className="block text-sm font-medium text-gray-500 mb-1">First Name</label>
            <p className="text-gray-900">{user.first_name}</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-500 mb-1">Last Name</label>
            <p className="text-gray-900">{user.last_name}</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-500 mb-1">Email</label>
            <p className="text-gray-900">{user.email}</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-500 mb-1">Phone</label>
            <p className="text-gray-900">{user.phone || 'Not provided'}</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-500 mb-1">Account Status</label>
            <span className={`inline-block px-3 py-1 rounded-full text-xs font-medium ${
              user.status === 'active' ? 'bg-green-100 text-green-800' :
              user.status === 'banned' ? 'bg-red-100 text-red-800' :
              'bg-gray-100 text-gray-800'
            }`}>
              {user.status}
            </span>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-500 mb-1">Email Verified</label>
            <span className={`inline-block px-3 py-1 rounded-full text-xs font-medium ${
              user.email_verified ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
            }`}>
              {user.email_verified ? 'Verified' : 'Not Verified'}
            </span>
          </div>
        </div>

        {/* Action buttons */}
        <div className="mt-8 pt-8 border-t border-gray-200 flex flex-wrap gap-4">
          <button className="btn-primary">
            Edit Profile
          </button>
          <button className="btn-outline">
            Change Password
          </button>
          {user.email_verified === false && (
            <button className="btn-secondary">
              Verify Email
            </button>
          )}
        </div>
      </div>

      {/* Account stats */}
      <div className="mt-6 grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-sm font-medium text-gray-500 mb-2">Total Orders</h3>
          <p className="text-3xl font-bold text-gray-900">0</p>
        </div>
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-sm font-medium text-gray-500 mb-2">Total Spent</h3>
          <p className="text-3xl font-bold text-gray-900">$0.00</p>
        </div>
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-sm font-medium text-gray-500 mb-2">Reviews Written</h3>
          <p className="text-3xl font-bold text-gray-900">0</p>
        </div>
      </div>
    </div>
  );
};

export default UserProfilePage;
