'use client'

import {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
  type CSSProperties,
} from 'react'
import dynamic from 'next/dynamic'
import { useRouter, useSearchParams } from 'next/navigation'
import type { BBox } from 'geojson'
import { FixedSizeList } from 'react-window'
import { type Store } from 'shared-configs'
import { fetchApi, getApiUrl } from '@/app/lib/api'

const DEFAULT_CENTER: [number, number] = [35.6812, 139.7671]
const MAP_LIMIT = 200

type GeoPoint = {
  lat: number
  lon: number
}

const MapView = dynamic(() => import('./components/features/MapView'), {
  ssr: false,
  loading: () => (
    <div className="h-full w-full rounded-2xl border border-slate-200 bg-slate-100 flex items-center justify-center text-sm text-slate-500 animate-pulse">
      🗺️ 地図を読み込み中...
    </div>
  ),
})

function parseBbox(value: string | null): BBox | null {
  if (!value) return null
  const parts = value.split(',').map((part) => Number(part.trim()))
  if (parts.length !== 4 || parts.some((part) => Number.isNaN(part))) {
    return null
  }
  return [parts[0], parts[1], parts[2], parts[3]]
}

function toBboxString(bbox: BBox | null) {
  return bbox ? bbox.map((value) => value.toFixed(6)).join(',') : ''
}

function haversineKm(lat1: number, lon1: number, lat2: number, lon2: number) {
  const toRad = (value: number) => (value * Math.PI) / 180
  const dLat = toRad(lat2 - lat1)
  const dLon = toRad(lon2 - lon1)
  const a =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos(toRad(lat1)) * Math.cos(toRad(lat2)) * Math.sin(dLon / 2) * Math.sin(dLon / 2)
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a))
  return 6371 * c
}

function isSameCenter(a: [number, number], b: [number, number]) {
  const epsilon = 0.00001
  return Math.abs(a[0] - b[0]) < epsilon && Math.abs(a[1] - b[1]) < epsilon
}

function isSameBbox(a: BBox | null, b: BBox) {
  if (!a) return false
  const epsilon = 0.00001
  return (
    Math.abs(a[0] - b[0]) < epsilon &&
    Math.abs(a[1] - b[1]) < epsilon &&
    Math.abs(a[2] - b[2]) < epsilon &&
    Math.abs(a[3] - b[3]) < epsilon
  )
}

export default function Home() {
  const router = useRouter()
  const searchParams = useSearchParams()

  const queryParam = searchParams.get('q') ?? ''
  const userLatParam = searchParams.get('user_lat')
  const userLonParam = searchParams.get('user_lon')

  const bboxFromQuery = useMemo(() => parseBbox(searchParams.get('bbox')), [searchParams])

  const [center, setCenter] = useState<[number, number]>(DEFAULT_CENTER)
  const [mapBbox, setMapBbox] = useState<BBox | null>(bboxFromQuery)
  const [mapZoom, setMapZoom] = useState(13)

  const [searchInput, setSearchInput] = useState(queryParam)

  const [stores, setStores] = useState<Store[]>([])
  const [storesLoading, setStoresLoading] = useState(false)
  const [storesError, setStoresError] = useState<string | null>(null)
  const [isTransitioning, setIsTransitioning] = useState(false)

  const [sortBy, setSortBy] = useState<'price' | 'distance'>('price')

  const [locating, setLocating] = useState(false)
  const [locationNotice, setLocationNotice] = useState<string | null>(null)

  const [userLocation, setUserLocation] = useState<GeoPoint | null>(null)
  const hasInitCenter = useRef(false)

  const API_URL = getApiUrl()

  const paramsRef = useRef(searchParams.toString())
  useEffect(() => {
    paramsRef.current = searchParams.toString()
  }, [searchParams])

  const updateQuery = useCallback(
    (updates: Record<string, string | null>) => {
      const params = new URLSearchParams(paramsRef.current)
      Object.entries(updates).forEach(([key, value]) => {
        if (!value) {
          params.delete(key)
        } else {
          params.set(key, value)
        }
      })
      const query = params.toString()
      paramsRef.current = query
      router.replace(query ? `?${query}` : '/', { scroll: false })
    },
    [router],
  )

  useEffect(() => {
    if (userLatParam && userLonParam) {
      const lat = Number(userLatParam)
      const lon = Number(userLonParam)
      if (Number.isFinite(lat) && Number.isFinite(lon)) {
        setUserLocation((prev) => {
          if (!prev) return { lat, lon }
          const epsilon = 0.00001
          if (Math.abs(prev.lat - lat) < epsilon && Math.abs(prev.lon - lon) < epsilon) {
            return prev
          }
          return { lat, lon }
        })
      }
    }
  }, [userLatParam, userLonParam])

  useEffect(() => {
    setSearchInput(queryParam)
  }, [queryParam])

  useEffect(() => {
    if (userLocation) {
      const nextCenter: [number, number] = [userLocation.lat, userLocation.lon]
      setCenter((prev) => (isSameCenter(prev, nextCenter) ? prev : nextCenter))
      hasInitCenter.current = true
      return
    }
    if (!hasInitCenter.current && bboxFromQuery) {
      const nextCenter: [number, number] = [
        (bboxFromQuery[1] + bboxFromQuery[3]) / 2,
        (bboxFromQuery[0] + bboxFromQuery[2]) / 2,
      ]
      setCenter((prev) => (isSameCenter(prev, nextCenter) ? prev : nextCenter))
      hasInitCenter.current = true
    }
  }, [userLocation, bboxFromQuery])

  const debouncedBbox = useMemo(() => mapBbox, [mapBbox])
  const bboxString = useMemo(() => toBboxString(debouncedBbox), [debouncedBbox])
  const lastBboxRef = useRef<string>('')

  useEffect(() => {
    if (!bboxString) return
    if (bboxString !== lastBboxRef.current) {
      updateQuery({ bbox: bboxString })
      lastBboxRef.current = bboxString
    }
  }, [bboxString, updateQuery])

  useEffect(() => {
    const loadStores = async () => {
      // Only load stores if there's a search query
      if (!bboxString || !queryParam) {
        setStores([])
        setStoresLoading(false)
        setIsTransitioning(false)
        return
      }

      setIsTransitioning(true)
      setStoresLoading(true)
      setStoresError(null)

      const params = new URLSearchParams()
      params.set('limit', MAP_LIMIT.toString())
      params.set('offset', '0')
      params.set('bbox', bboxString)
      params.set('q', queryParam)

      if (userLocation) {
        params.set('user_lat', userLocation.lat.toString())
        params.set('user_lon', userLocation.lon.toString())
      }

      // Use user-selected sort if available, otherwise use smart defaults
      const effectiveSort = sortBy === 'distance' && userLocation ? 'distance' : 'price'
      params.set('sort', effectiveSort)
      params.set('order', 'asc')

      const url = `${API_URL}/api/stores?${params.toString()}`
      const data = await fetchApi<Store[]>(url)

      // Small delay for smooth transition
      await new Promise(resolve => setTimeout(resolve, 150))

      if (data.error) {
        setStoresError(data.error.message)
        setStoresLoading(false)
        setIsTransitioning(false)
        return
      }

      setStores(data.data ?? [])
      setStoresLoading(false)
      setTimeout(() => setIsTransitioning(false), 50)
    }

    void loadStores()
  }, [API_URL, bboxString, queryParam, userLocation, sortBy])

  const handleBoundsChange = useCallback((bbox: BBox, zoom: number) => {
    setMapBbox((prev) => (isSameBbox(prev, bbox) ? prev : bbox))
    setMapZoom((prev) => (prev === zoom ? prev : zoom))
  }, [])

  const handleUseLocation = () => {
    if (!navigator.geolocation) return
    setLocating(true)
    navigator.geolocation.getCurrentPosition(
      (position) => {
        const lat = position.coords.latitude
        const lon = position.coords.longitude
        setUserLocation({ lat, lon })
        updateQuery({ user_lat: lat.toString(), user_lon: lon.toString() })
        setCenter([lat, lon])
        setLocating(false)
        setLocationNotice('現在地を更新しました')
        window.setTimeout(() => setLocationNotice(null), 2500)
      },
      () => {
        setLocating(false)
        setLocationNotice('現在地を取得できませんでした')
        window.setTimeout(() => setLocationNotice(null), 2500)
      },
      { enableHighAccuracy: true, timeout: 8000 },
    )
  }

  const listItems = useMemo(() => {
    return stores.map((store) => {
      let distanceKm: number | null = null
      if (store.distance !== undefined && store.distance !== null) {
        distanceKm = store.distance / 1000
      } else if (userLocation) {
        distanceKm = haversineKm(userLocation.lat, userLocation.lon, store.latitude, store.longitude)
      }
      return { store, distanceKm }
    })
  }, [stores, userLocation])

  const listHeight = useMemo(() => (listItems.length > 6 ? 520 : 320), [listItems.length])

  return (
    <div className="min-h-screen bg-slate-50 text-slate-900">
      <header className="px-6 py-5 border-b border-slate-200 bg-gradient-to-r from-white to-slate-50">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <span className="text-3xl">🛒</span>
            <div>
              <h1 className="text-2xl font-bold bg-gradient-to-r from-sky-600 to-cyan-600 bg-clip-text text-transparent">
                kaimono
              </h1>
              <p className="text-xs text-slate-500">価格比較アプリ</p>
            </div>
          </div>
          <div className="hidden md:flex items-center gap-2 text-sm text-slate-500">
            <span>💡</span>
            <span>お得な店舗を見つけよう</span>
          </div>
        </div>
      </header>

      <section className="px-6 py-4 border-b border-slate-200 bg-white shadow-sm">
        <div className="flex flex-col gap-3">
          <div className="flex flex-wrap items-center gap-3">
            <div className="relative flex-1 min-w-[280px] max-w-md">
              <span className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 text-lg">🔍</span>
              <input
                value={searchInput}
                onChange={(event) => setSearchInput(event.target.value)}
                onKeyDown={(event) => {
                  if (event.key === 'Enter') {
                    updateQuery({ q: searchInput })
                  }
                }}
                placeholder="商品名を入力（例: りんご、牛乳）"
                className="w-full rounded-full border border-slate-300 pl-11 pr-4 py-2.5 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-sky-500 transition-all"
              />
            </div>
            <button
              onClick={() => updateQuery({ q: searchInput })}
              className="rounded-full bg-gradient-to-r from-sky-500 to-cyan-500 px-5 py-2.5 text-sm font-medium text-white hover:from-sky-600 hover:to-cyan-600 transition-all shadow-sm hover:shadow-md flex items-center gap-1.5"
            >
              <span>🔎</span>
              <span>検索</span>
            </button>
            <button
              onClick={handleUseLocation}
              disabled={locating}
              className="rounded-full border-2 border-slate-300 px-4 py-2 text-sm hover:border-sky-400 hover:bg-sky-50 transition-all flex items-center gap-1.5 disabled:opacity-50"
            >
              <span>{locating ? '⌛' : '📍'}</span>
              <span>{locating ? '取得中...' : '現在地'}</span>
            </button>
            {locationNotice && (
              <span className="rounded-full bg-emerald-50 px-3 py-1.5 text-xs text-emerald-700 border border-emerald-200 animate-fade-in flex items-center gap-1">
                <span>✓</span>
                {locationNotice}
              </span>
            )}
          </div>

          {/* Search tags and sorting */}
          <div className="flex flex-wrap items-center gap-3">
            {queryParam && (
              <div className="flex items-center gap-2 animate-slide-up">
                <span className="text-xs text-slate-600">検索中:</span>
                <button
                  onClick={() => {
                    updateQuery({ q: null })
                    setSearchInput('')
                  }}
                  className="group flex items-center gap-1.5 bg-gradient-to-r from-sky-100 to-cyan-100 hover:from-sky-200 hover:to-cyan-200 border border-sky-300 rounded-full px-3 py-1.5 text-sm font-medium text-sky-700 transition-all hover:shadow-sm"
                >
                  <span className="text-xs">🏷️</span>
                  <span>{queryParam}</span>
                  <span className="text-xs group-hover:rotate-90 transition-transform">✕</span>
                </button>
              </div>
            )}

            {queryParam && (
              <div className="flex items-center gap-2">
                <span className="text-xs text-slate-600 font-medium">並び替え:</span>
                <div className="flex gap-2">
                  <button
                    onClick={() => setSortBy('price')}
                    className={`px-3 py-1.5 rounded-full text-xs font-medium transition-all ${
                      sortBy === 'price'
                        ? 'bg-gradient-to-r from-emerald-500 to-teal-500 text-white shadow-sm'
                        : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                    }`}
                  >
                    💰 価格順
                  </button>
                  <button
                    onClick={() => setSortBy('distance')}
                    disabled={!userLocation}
                    className={`px-3 py-1.5 rounded-full text-xs font-medium transition-all disabled:opacity-40 disabled:cursor-not-allowed ${
                      sortBy === 'distance'
                        ? 'bg-gradient-to-r from-blue-500 to-cyan-500 text-white shadow-sm'
                        : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                    }`}
                  >
                    📏 距離順
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      </section>

      <main className="grid grid-cols-1 lg:grid-cols-[1.35fr_1fr] gap-6 px-6 py-6">
        <div className="h-[520px] lg:h-[640px]">
          <MapView
            stores={stores}
            center={center}
            userLocation={userLocation}
            onBoundsChange={handleBoundsChange}
            onSelectStore={(store) =>
              router.push(`/stores/${store.id}?q=${encodeURIComponent(queryParam)}`)
            }
          />
        </div>

        <aside className="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center gap-2">
              <span className="text-lg">🏪</span>
              <h2 className="text-lg font-semibold text-slate-800">店舗リスト</h2>
            </div>
            <span className="text-xs text-slate-500 bg-slate-100 px-2.5 py-1 rounded-full font-medium">
              {stores.length} 件
            </span>
          </div>
          {storesError && (
            <p className="mt-3 rounded-lg bg-red-50 border border-red-100 px-3 py-2 text-sm text-red-700 flex items-center gap-2">
              <span>⚠️</span>
              {storesError}
            </p>
          )}
          {storesLoading && !storesError && (
            <div className="mt-3 flex items-center gap-2 text-sm text-slate-500">
              <span className="animate-spin">⏳</span>
              <span>読み込み中...</span>
            </div>
          )}
          {!storesLoading && stores.length === 0 && !queryParam && (
            <div className="mt-8 text-center py-12">
              <div className="text-5xl mb-3 animate-pulse-slow">🔎</div>
              <p className="text-base font-medium text-slate-700 mb-2">商品名を検索してください</p>
              <p className="text-xs text-slate-500">例: りんご、牛乳、パンなど</p>
            </div>
          )}
          {!storesLoading && stores.length === 0 && queryParam && (
            <div className="mt-8 text-center py-12">
              <div className="text-5xl mb-3">😔</div>
              <p className="text-base font-medium text-slate-700 mb-2">該当する店舗がありません</p>
              <p className="text-xs text-slate-500">別のキーワードで検索してみてください</p>
            </div>
          )}
          {stores.length > 0 && (
            <div className={`mt-4 transition-opacity duration-300 ${isTransitioning ? 'opacity-0' : 'opacity-100'}`}>
              <FixedSizeList height={listHeight} itemCount={listItems.length} itemSize={90} width="100%">
                {({ index, style }: { index: number; style: CSSProperties }) => {
                  const item = listItems[index]
                  const distanceLabel =
                    item.distanceKm !== null ? `${item.distanceKm.toFixed(1)} km` : '—'
                  return (
                    <button
                      style={{
                        ...style,
                        transition: 'all 0.2s ease-out',
                      }}
                      onClick={() =>
                        router.push(`/stores/${item.store.id}?q=${encodeURIComponent(queryParam)}`)
                      }
                      className="w-full px-3 py-2 text-left hover:bg-gradient-to-r hover:from-sky-50 hover:to-cyan-50 rounded-lg transition-all group"
                    >
                      <div className="text-sm font-semibold text-slate-800 group-hover:text-sky-700 transition-colors">
                        {item.store.name}
                      </div>
                      <div className="text-xs text-slate-500 mt-0.5 flex items-start gap-1">
                        <span className="text-[10px] mt-0.5">📍</span>
                        <span className="line-clamp-1">{item.store.address}</span>
                      </div>
                      <div className="mt-1.5 flex flex-wrap items-center gap-2 text-xs">
                        <span className="text-slate-500 flex items-center gap-1">
                          <span className="text-[10px]">📏</span>
                          {distanceLabel}
                        </span>
                        {item.store.min_price !== undefined && item.store.min_price !== null && (
                          <span className="bg-gradient-to-r from-emerald-500 to-teal-500 text-white px-2 py-0.5 rounded-full font-medium">
                            ¥{item.store.min_price.toFixed(0)}
                          </span>
                        )}
                      </div>
                    </button>
                  )
                }}
              </FixedSizeList>
            </div>
          )}
        </aside>
      </main>
    </div>
  )
}
