export interface ApiError {
  code: string
  message: string
}

export interface Meta {
  count?: number
  limit?: number
  offset?: number
}

export interface ApiResponse<T> {
  data?: T
  meta?: Meta
  error?: ApiError
}

export interface Store {
  id: number
  name: string
  address: string
  phone?: string
  latitude: number
  longitude: number
  distance?: number
  min_price?: number
  created_at: string
  updated_at: string
}

export interface Product {
  id: number
  name: string
  category: string
  barcode: string
  created_at: string
}

export interface Price {
  id: number
  store_id: number
  product_id: number
  price: number
  currency: string
  recorded_at: string
  created_at: string
  store?: Store
  product?: Product
}

export interface PriceSummary {
  min_price?: number
  max_price?: number
  avg_price?: number
  currency?: string
}

export interface DailyPriceStats {
  date: string
  avg_price: number
  min_price: number
  max_price: number
  count: number
}

export interface StorePriceStats {
  store_id: number
  category?: string
  query?: string
  days: number
  summary: PriceSummary
  daily: DailyPriceStats[]
}
