function CartPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Shopping Cart</h1>
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2 bg-white rounded-lg shadow-md p-6">
          <p className="text-gray-600">Your cart items will appear here</p>
        </div>
        <div className="bg-white rounded-lg shadow-md p-6 h-fit">
          <h2 className="text-xl font-bold mb-4">Order Summary</h2>
          <div className="space-y-3">
            <div className="flex justify-between">
              <span className="text-gray-600">Subtotal</span>
              <span className="font-semibold">$0.00</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-600">Shipping</span>
              <span className="font-semibold">$0.00</span>
            </div>
            <hr />
            <div className="flex justify-between text-lg font-bold">
              <span>Total</span>
              <span className="text-primary-500">$0.00</span>
            </div>
            <button className="w-full bg-primary-500 hover:bg-primary-600 text-white font-semibold py-3 rounded-lg transition mt-4">
              Proceed to Checkout
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

export default CartPage
