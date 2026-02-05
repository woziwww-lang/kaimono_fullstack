'use client'

import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { useParams, useRouter, useSearchParams } from 'next/navigation'
import { type DailyPriceStats, type Store, type StorePriceStats } from 'shared-configs'
import { fetchApi, getApiUrl } from '@/app/lib/api'
import { formatPrice } from '@/app/lib/format'
import { PriceTrendChart } from '@/app/components/features/PriceTrendChart'

type ChartPoint = {
  date: string
  avg: number
}
export default function StoreDetailPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const params = useParams()
  const idParam = Array.isArray(params?.id) ? params?.id[0] : params?.id
  const storeId = Number(idParam)

  const queryParam = searchParams.get('q') ?? ''

  const [queryInput, setQueryInput] = useState(queryParam)

  const [store, setStore] = useState<Store | null>(null)
  const [storeError, setStoreError] = useState<string | null>(null)
  const [stats, setStats] = useState<StorePriceStats | null>(null)
  const [statsError, setStatsError] = useState<string | null>(null)
  const [statsLoading, setStatsLoading] = useState(false)

  const API_URL = getApiUrl()
  const basePath = useMemo(() => (Number.isFinite(storeId) ? `/stores/${storeId}` : '/stores'), [storeId])

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
      router.replace(query ? `${basePath}?${query}` : basePath, { scroll: false })
    },
    [router, basePath],
  )

  useEffect(() => {
    setQueryInput(queryParam)
  }, [queryParam])

  useEffect(() => {
    if (!Number.isFinite(storeId)) return
    const loadStore = async () => {
      const url = `${API_URL}/api/stores/${storeId}`
      const data = await fetchApi<Store>(url)
      if (data.error) {
        setStoreError(data.error.message)
        return
      }
      setStore(data.data ?? null)
      setStoreError(null)
    }
    void loadStore()
  }, [API_URL, storeId])

  useEffect(() => {
    if (!Number.isFinite(storeId)) return
    const loadStats = async () => {
      setStatsLoading(true)
      setStatsError(null)
      const params = new URLSearchParams()
      if (queryParam) params.set('q', queryParam)
      params.set('days', '14')
      const url = `${API_URL}/api/stores/${storeId}/price-stats?${params.toString()}`
      const data = await fetchApi<StorePriceStats>(url)
      if (data.error) {
        setStatsError(data.error.message)
        setStatsLoading(false)
        return
      }
      setStats(data.data ?? null)
      setStatsLoading(false)
    }
    void loadStats()
  }, [API_URL, storeId, queryParam])

  const daily = useMemo<DailyPriceStats[]>(() => stats?.daily ?? [], [stats])
  const chartData = useMemo<ChartPoint[]>(
    () =>
      daily.map((item) => ({
        date: item.date,
        avg: item.avg_price,
      })),
    [daily],
  )

  const summary = stats?.summary

  if (!Number.isFinite(storeId)) {
    return (
      <div className="min-h-screen bg-slate-50 p-6 text-slate-900">
        <p className="text-sm text-red-600">店舗IDが正しくありません。</p>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-slate-50 text-slate-900">
      <header className="px-6 py-5 border-b border-slate-200 bg-gradient-to-r from-white to-slate-50">
        <div className="flex items-center justify-between">
          <button
            onClick={() => router.back()}
            className="rounded-full border-2 border-slate-200 px-4 py-2 text-sm hover:border-sky-400 hover:bg-sky-50 transition-all flex items-center gap-2 group"
          >
            <span className="group-hover:-translate-x-1 transition-transform">←</span>
            <span>戻る</span>
          </button>
          <div className="flex items-center gap-2 text-sm text-slate-500">
            <span>🏪</span>
            <span className="hidden sm:inline">店舗詳細</span>
          </div>
        </div>
      </header>

      <main className="px-6 py-6 space-y-6">
        <section className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm animate-slide-up">
          {storeError && (
            <div className="flex items-center gap-2 text-sm text-red-600 bg-red-50 p-3 rounded-lg">
              <span>⚠️</span>
              <p>{storeError}</p>
            </div>
          )}
          {!store && !storeError && (
            <div className="flex items-center gap-2 text-sm text-slate-500">
              <span className="animate-pulse-slow">⏳</span>
              <p>読み込み中...</p>
            </div>
          )}
          {store && (
            <>
              <p className="text-xs uppercase tracking-[0.25em] text-slate-400 mb-2">店舗情報</p>
              <h1 className="text-2xl font-bold text-slate-900 mb-3">{store.name}</h1>
              <div className="space-y-2">
                <div className="flex items-start gap-2 text-sm hover:bg-slate-50 p-2 rounded-lg transition-colors">
                  <span className="text-base mt-0.5">📍</span>
                  <span className="text-slate-600">{store.address}</span>
                </div>
                {store.phone && (
                  <div className="flex items-center gap-2 text-sm hover:bg-sky-50 p-2 rounded-lg transition-colors">
                    <span className="text-base">📞</span>
                    <a
                      href={`tel:${store.phone}`}
                      className="text-sky-600 hover:text-sky-700 hover:underline transition-colors"
                    >
                      {store.phone}
                    </a>
                  </div>
                )}
              </div>
            </>
          )}
        </section>

        <section className="rounded-2xl border border-slate-200 bg-white p-5 space-y-4 shadow-sm animate-slide-up" style={{ animationDelay: '0.1s' }}>
          <div className="flex flex-col gap-3">
            <div className="flex flex-wrap items-center gap-3">
              <div className="relative flex-1 max-w-md">
                <span className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400">🔎</span>
                <input
                  value={queryInput}
                  onChange={(event) => setQueryInput(event.target.value)}
                  onKeyDown={(event) => {
                    if (event.key === 'Enter') {
                      updateQuery({ q: queryInput || null })
                    }
                  }}
                  placeholder="品名で絞り込み（例: りんご）"
                  className="w-full rounded-full border border-slate-300 pl-10 pr-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-sky-500 transition-all"
                />
              </div>
              <button
                onClick={() => updateQuery({ q: queryInput || null })}
                className="rounded-full bg-gradient-to-r from-sky-500 to-cyan-500 px-5 py-2.5 text-sm font-medium text-white hover:from-sky-600 hover:to-cyan-600 transition-all shadow-sm hover:shadow-md hover:scale-105 active:scale-95"
              >
                更新
              </button>
            </div>

            {queryParam && (
              <div className="flex items-center gap-2 animate-slide-up">
                <span className="text-xs text-slate-600">絞り込み中:</span>
                <button
                  onClick={() => {
                    updateQuery({ q: null })
                    setQueryInput('')
                  }}
                  className="group flex items-center gap-1.5 bg-gradient-to-r from-sky-100 to-cyan-100 hover:from-sky-200 hover:to-cyan-200 border border-sky-300 rounded-full px-3 py-1.5 text-sm font-medium text-sky-700 transition-all hover:shadow-sm"
                >
                  <span className="text-xs">🏷️</span>
                  <span>{queryParam}</span>
                  <span className="text-xs group-hover:rotate-90 transition-transform">✕</span>
                </button>
              </div>
            )}
          </div>

          {statsError && (
            <div className="flex items-center gap-2 text-sm text-red-600 bg-red-50 p-3 rounded-lg border border-red-100">
              <span>⚠️</span>
              <p>{statsError}</p>
            </div>
          )}
          {statsLoading && !statsError && (
            <div className="flex items-center gap-2 text-sm text-slate-500">
              <span className="animate-pulse-slow">⏳</span>
              <p>統計を取得中...</p>
            </div>
          )}
          {!statsLoading && !queryParam && (
            <div className="text-center py-8">
              <div className="text-4xl mb-2 animate-pulse-slow">🔍</div>
              <p className="text-sm font-medium text-slate-700 mb-1">商品名で絞り込んでください</p>
              <p className="text-xs text-slate-500">価格推移を表示します</p>
            </div>
          )}

          {queryParam && stats && (
            <>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
            <div className="rounded-xl bg-gradient-to-br from-emerald-50 to-white border border-emerald-100 p-3 transition-transform hover:scale-105">
              <div className="flex items-center gap-1.5 mb-1">
                <span>💰</span>
                <p className="text-xs text-emerald-700 font-medium">最低</p>
              </div>
              <p className="text-xl font-bold text-emerald-600">{formatPrice(summary?.min_price)}</p>
            </div>
            <div className="rounded-xl bg-gradient-to-br from-rose-50 to-white border border-rose-100 p-3 transition-transform hover:scale-105">
              <div className="flex items-center gap-1.5 mb-1">
                <span>📈</span>
                <p className="text-xs text-rose-700 font-medium">最高</p>
              </div>
              <p className="text-xl font-bold text-rose-600">{formatPrice(summary?.max_price)}</p>
            </div>
            <div className="rounded-xl bg-gradient-to-br from-sky-50 to-white border border-sky-100 p-3 transition-transform hover:scale-105">
              <div className="flex items-center gap-1.5 mb-1">
                <span>📊</span>
                <p className="text-xs text-sky-700 font-medium">平均</p>
              </div>
              <p className="text-xl font-bold text-sky-600">{formatPrice(summary?.avg_price)}</p>
            </div>
          </div>

          <div className="rounded-xl border border-slate-100 bg-gradient-to-br from-white to-slate-50 p-4">
            <div className="flex items-center justify-between mb-3">
              <div className="flex items-center gap-2">
                <span>📉</span>
                <p className="text-sm font-semibold text-slate-800">価格推移</p>
              </div>
              <span className="text-xs text-slate-500 bg-white px-2 py-1 rounded-full">{daily.length}日</span>
            </div>
            <div className="bg-white rounded-lg p-2 shadow-sm">
              <PriceTrendChart data={chartData} />
            </div>
          </div>
            </>
          )}
        </section>
      </main>
    </div>
  )
}
