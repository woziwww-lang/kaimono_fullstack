'use client'

import { useState, useEffect } from 'react'

interface Store {
  id: number
  name: string
  address: string
  latitude: number
  longitude: number
  distance?: number
}

interface Product {
  id: number
  name: string
  category: string
  barcode: string
}

export default function Home() {
  const [stores, setStores] = useState<Store[]>([])
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

  // Fetch all stores
  const fetchStores = async () => {
    try {
      const response = await fetch(`${API_URL}/api/stores`)
      if (!response.ok) throw new Error('Failed to fetch stores')
      const data = await response.json()
      setStores(data.stores || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load stores')
    }
  }

  // Fetch all products
  const fetchProducts = async () => {
    try {
      const response = await fetch(`${API_URL}/api/products`)
      if (!response.ok) throw new Error('Failed to fetch products')
      const data = await response.json()
      setProducts(data.products || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load products')
    }
  }

  // Fetch nearby stores (using Tokyo Station as default location)
  const fetchNearbyStores = async () => {
    const lat = 35.6812 // Tokyo Station
    const lon = 139.7671
    const radius = 5000 // 5km

    try {
      const response = await fetch(
        `${API_URL}/api/stores/nearby?lat=${lat}&lon=${lon}&radius=${radius}`
      )
      if (!response.ok) throw new Error('Failed to fetch nearby stores')
      const data = await response.json()
      setStores(data.stores || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load nearby stores')
    }
  }

  useEffect(() => {
    const loadData = async () => {
      setLoading(true)
      await Promise.all([fetchStores(), fetchProducts()])
      setLoading(false)
    }
    loadData()
  }, [])

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-xl">読み込み中...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-xl text-red-600">
          エラー: {error}
          <br />
          <button
            onClick={() => window.location.reload()}
            className="mt-4 px-4 py-2 bg-primary-600 text-white rounded hover:bg-primary-700"
          >
            再読み込み
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-8">
      <section>
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-2xl font-bold">近くの店舗</h2>
          <button
            onClick={fetchNearbyStores}
            className="px-4 py-2 bg-primary-600 text-white rounded hover:bg-primary-700 transition"
          >
            東京駅周辺で検索
          </button>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {stores.map((store) => (
            <div key={store.id} className="bg-white rounded-lg shadow-md p-4 border border-gray-200">
              <h3 className="text-lg font-semibold mb-2">{store.name}</h3>
              <p className="text-gray-600 text-sm mb-2">{store.address}</p>
              {store.distance && (
                <p className="text-primary-600 text-sm font-medium">
                  距離: {(store.distance / 1000).toFixed(2)} km
                </p>
              )}
              <div className="mt-2 text-xs text-gray-500">
                <span>緯度: {store.latitude.toFixed(4)}</span>
                <br />
                <span>経度: {store.longitude.toFixed(4)}</span>
              </div>
            </div>
          ))}
        </div>
        {stores.length === 0 && (
          <p className="text-gray-500 text-center py-8">店舗が見つかりません</p>
        )}
      </section>

      <section>
        <h2 className="text-2xl font-bold mb-4">商品一覧</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {products.map((product) => (
            <div key={product.id} className="bg-white rounded-lg shadow-md p-4 border border-gray-200">
              <h3 className="text-lg font-semibold mb-2">{product.name}</h3>
              <div className="text-sm text-gray-600">
                <p>カテゴリ: {product.category}</p>
                <p className="text-xs mt-2 font-mono">{product.barcode}</p>
              </div>
            </div>
          ))}
        </div>
        {products.length === 0 && (
          <p className="text-gray-500 text-center py-8">商品が見つかりません</p>
        )}
      </section>
    </div>
  )
}
