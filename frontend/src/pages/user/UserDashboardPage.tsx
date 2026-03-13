function UserDashboardPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">My Account</h1>
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        {/* Sidebar */}
        <aside className="md:col-span-1">
          <div className="bg-white rounded-lg shadow-md p-6">
            <div className="text-center mb-6">
              <div className="w-20 h-20 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-3">
                <span className="text-2xl font-bold text-primary-500">U</span>
              </div>
              <h3 className="font-semibold">User Name</h3>
              <p className="text-gray-600 text-sm">user@example.com</p>
            </div>
            <nav className="space-y-2">
              <a href="/account" className="block px-4 py-2 bg-primary-500 text-white rounded-lg">Dashboard</a>
              <a href="/account/orders" className="block px-4 py-2 hover:bg-gray-100 rounded-lg">Orders</a>
              <a href="/account/wishlist" className="block px-4 py-2 hover:bg-gray-100 rounded-lg">Wishlist</a>
              <a href="/account/addresses" className="block px-4 py-2 hover:bg-gray-100 rounded-lg">Addresses</a>
            </nav>
          </div>
        </aside>

        {/* Content */}
        <div className="md:col-span-3 bg-white rounded-lg shadow-md p-6">
          <h2 className="text-2xl font-bold mb-4">Dashboard</h2>
          <p className="text-gray-600">Welcome to your account dashboard. Here you can manage your orders, wishlist, and addresses.</p>
        </div>
      </div>
    </div>
  )
}

export default UserDashboardPage
