import HeroBanner from '@components/home/HeroBanner'
import CategoriesSection from '@components/home/CategoriesSection'
import FlashSaleSection from '@components/home/FlashSaleSection'
import FeaturedProducts from '@components/home/FeaturedProducts'
import RecommendedProducts from '@components/home/RecommendedProducts'

function HomePage() {
  return (
    <div className="animate-fade-in">
      <HeroBanner />
      <CategoriesSection />
      <FlashSaleSection />
      <FeaturedProducts />
      <RecommendedProducts />
    </div>
  )
}

export default HomePage
