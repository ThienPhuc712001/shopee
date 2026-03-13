function AdminDashboardPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Admin Dashboard</h1>
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        {/* Stats Cards */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-gray-600 text-sm mb-2">Total Sales</h3>
          <p className="text-3xl font-bold text-primary-500">$12,345</p>
          <p className="text-green-500 text-sm mt-2">↑ 12% from last month</p>
        </div>
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-gray-600 text-sm mb-2">Total Orders</h3>
          <p className="text-3xl font-bold text-primary-500">1,234</p>
          <p className="text-green-500 text-sm mt-2">↑ 8% from last month</p>
        </div>
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-gray-600 text-sm mb-2">Total Products</h3>
          <p className="text-3xl font-bold text-primary-500">567</p>
          <p className="text-gray-500 text-sm mt-2">Active products</p>
        </div>
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-gray-600 text-sm mb-2">Total Users</h3>
          <p className="text-3xl font-bold text-primary-500">8,901</p>
          <p className="text-green-500 text-sm mt-2">↑ 15% from last month</p>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-2xl font-bold mb-4">Admin Navigation</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <a href="/admin/products" className="p-4 bg-primary-500 text-white rounded-lg text-center hover:bg-primary-600 transition">
            📦 Products
          </a>
          <a href="/admin/orders" className="p-4 bg-secondary-800 text-white rounded-lg text-center hover:bg-secondary-700 transition">
            📋 Orders
          </a>
          <a href="/admin/users" className="p-4 bg-green-500 text-white rounded-lg text-center hover:bg-green-600 transition">
            👥 Users
          </a>
          <a href="/admin/analytics" className="p-4 bg-purple-500 text-white rounded-lg text-center hover:bg-purple-600 transition">
            📊 Analytics
          </a>
        </div>
      </div>
    </div>
  )
}

export default AdminDashboardPage
