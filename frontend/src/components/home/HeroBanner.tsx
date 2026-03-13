import { useState, useEffect } from 'react'
import { ChevronLeft, ChevronRight } from 'lucide-react'
import { motion, AnimatePresence } from 'framer-motion'

const banners = [
  {
    id: 1,
    title: 'Summer Sale 2026',
    subtitle: 'Up to 70% Off',
    description: 'Discover amazing deals on your favorite products',
    image: 'https://images.unsplash.com/photo-1607082348824-0a96f2a4b9da?w=1200&h=500&fit=crop',
    cta: 'Shop Now',
    link: '/products?sale=summer',
  },
  {
    id: 2,
    title: 'New Arrivals',
    subtitle: 'Fresh Styles',
    description: 'Check out the latest trends',
    image: 'https://images.unsplash.com/photo-1441986300917-64674bd600d8?w=1200&h=500&fit=crop',
    cta: 'Explore',
    link: '/products?sort=newest',
  },
  {
    id: 3,
    title: 'Electronics Week',
    subtitle: 'Tech Deals',
    description: 'Save big on gadgets and electronics',
    image: 'https://images.unsplash.com/photo-1550009158-9ebf69173e03?w=1200&h=500&fit=crop',
    cta: 'Shop Electronics',
    link: '/category/1',
  },
]

function HeroBanner() {
  const [currentIndex, setCurrentIndex] = useState(0)

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentIndex((prev) => (prev + 1) % banners.length)
    }, 5000)
    return () => clearInterval(timer)
  }, [])

  const goToPrevious = () => {
    setCurrentIndex((prev) => (prev - 1 + banners.length) % banners.length)
  }

  const goToNext = () => {
    setCurrentIndex((prev) => (prev + 1) % banners.length)
  }

  return (
    <section className="relative h-[400px] md:h-[500px] overflow-hidden bg-gray-100">
      <AnimatePresence mode="wait">
        <motion.div
          key={currentIndex}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.5 }}
          className="absolute inset-0"
        >
          <div
            className="w-full h-full bg-cover bg-center"
            style={{ backgroundImage: `url(${banners[currentIndex].image})` }}
          >
            <div className="absolute inset-0 bg-gradient-to-r from-black/70 via-black/50 to-transparent" />
          </div>

          <div className="absolute inset-0 flex items-center">
            <div className="container mx-auto px-4">
              <motion.div
                initial={{ opacity: 0, y: 30 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.3 }}
                className="max-w-xl text-white"
              >
                <span className="inline-block bg-primary-500 text-white px-4 py-1 rounded-full text-sm font-semibold mb-4">
                  {banners[currentIndex].subtitle}
                </span>
                <h1 className="text-4xl md:text-6xl font-bold mb-4">
                  {banners[currentIndex].title}
                </h1>
                <p className="text-lg md:text-xl text-gray-200 mb-8">
                  {banners[currentIndex].description}
                </p>
                <a
                  href={banners[currentIndex].link}
                  className="inline-block bg-primary-500 hover:bg-primary-600 text-white font-semibold px-8 py-4 rounded-lg transition-all hover:scale-105"
                >
                  {banners[currentIndex].cta}
                </a>
              </motion.div>
            </div>
          </div>
        </motion.div>
      </AnimatePresence>

      {/* Navigation Arrows */}
      <button
        onClick={goToPrevious}
        className="absolute left-4 top-1/2 -translate-y-1/2 bg-white/90 hover:bg-white p-3 rounded-full transition-all hover:scale-110 shadow-lg z-10"
      >
        <ChevronLeft className="w-6 h-6 text-gray-800" />
      </button>
      <button
        onClick={goToNext}
        className="absolute right-4 top-1/2 -translate-y-1/2 bg-white/90 hover:bg-white p-3 rounded-full transition-all hover:scale-110 shadow-lg z-10"
      >
        <ChevronRight className="w-6 h-6 text-gray-800" />
      </button>

      {/* Dots Indicator */}
      <div className="absolute bottom-6 left-1/2 -translate-x-1/2 flex gap-3 z-10">
        {banners.map((_, index) => (
          <button
            key={index}
            onClick={() => setCurrentIndex(index)}
            className={`w-3 h-3 rounded-full transition-all ${
              index === currentIndex
                ? 'bg-primary-500 w-8'
                : 'bg-white/50 hover:bg-white/80'
            }`}
          />
        ))}
      </div>
    </section>
  )
}

export default HeroBanner
