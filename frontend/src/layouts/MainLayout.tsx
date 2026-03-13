import { Outlet } from 'react-router-dom'
import Header from '@components/common/Header'
import Footer from '@components/common/Footer'

function MainLayout() {
  return (
    <div className="min-h-screen flex flex-col bg-background">
      <Header />
      <main className="flex-1">
        <Outlet />
      </main>
      <Footer />
    </div>
  )
}

export default MainLayout
