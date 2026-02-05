'use client'

import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import {
  CircleMarker,
  MapContainer,
  Marker,
  TileLayer,
  useMap,
  useMapEvents,
} from 'react-leaflet'
import L from 'leaflet'
import Supercluster from 'supercluster'
import type { BBox, Feature, Point } from 'geojson'
import type { Store } from 'shared-configs'

import iconRetinaUrl from 'leaflet/dist/images/marker-icon-2x.png'
import iconUrl from 'leaflet/dist/images/marker-icon.png'
import shadowUrl from 'leaflet/dist/images/marker-shadow.png'

L.Icon.Default.mergeOptions({
  iconRetinaUrl,
  iconUrl,
  shadowUrl,
})

type GeoPoint = {
  lat: number
  lon: number
}

type ClusterProperties = {
  cluster: boolean
  point_count?: number
  store?: Store
}

type ClusterFeature = Feature<Point, ClusterProperties> & { id?: number }

type MapViewProps = {
  stores: Store[]
  center: [number, number]
  userLocation?: GeoPoint | null
  onBoundsChange: (bbox: BBox, zoom: number) => void
  onSelectStore: (store: Store) => void
}

function MapEvents({ onBoundsChange }: { onBoundsChange: (bbox: BBox, zoom: number) => void }) {
  const map = useMapEvents({
    moveend() {
      const bounds = map.getBounds()
      const bbox: BBox = [
        bounds.getWest(),
        bounds.getSouth(),
        bounds.getEast(),
        bounds.getNorth(),
      ]
      onBoundsChange(bbox, map.getZoom())
    },
    zoomend() {
      const bounds = map.getBounds()
      const bbox: BBox = [
        bounds.getWest(),
        bounds.getSouth(),
        bounds.getEast(),
        bounds.getNorth(),
      ]
      onBoundsChange(bbox, map.getZoom())
    },
  })

  useEffect(() => {
    const bounds = map.getBounds()
    const bbox: BBox = [
      bounds.getWest(),
      bounds.getSouth(),
      bounds.getEast(),
      bounds.getNorth(),
    ]
    onBoundsChange(bbox, map.getZoom())
  }, [map, onBoundsChange])

  return null
}

function MapCenter({ center }: { center: [number, number] }) {
  const map = useMap()
  useEffect(() => {
    const current = map.getCenter()
    const epsilon = 0.00001
    const shouldMove =
      Math.abs(current.lat - center[0]) > epsilon || Math.abs(current.lng - center[1]) > epsilon
    if (shouldMove) {
      map.setView(center, map.getZoom(), { animate: true })
    }
  }, [center, map])
  return null
}

function ClusterMarker({
  cluster,
  supercluster,
}: {
  cluster: ClusterFeature
  supercluster: Supercluster<ClusterProperties>
}) {
  const map = useMap()
  const [lon, lat] = cluster.geometry.coordinates
  const pointCount = cluster.properties.point_count ?? 0
  const size = Math.max(36, Math.min(60, 20 + pointCount))
  const icon = L.divIcon({
    html: `<div class="cluster-bubble" style="width:${size}px;height:${size}px">${pointCount}</div>`,
    className: 'cluster-marker',
    iconSize: [size, size],
  })

  return (
    <Marker
      position={[lat, lon]}
      icon={icon}
      eventHandlers={{
        click: () => {
          if (cluster.id == null) return
          const expansionZoom = supercluster.getClusterExpansionZoom(cluster.id)
          map.setView([lat, lon], expansionZoom, { animate: true })
        },
      }}
    />
  )
}

function buildPriceIcon(price: number) {
  const display = price >= 1000 ? `¥${Math.round(price)}` : `¥${price.toFixed(0)}`
  return L.divIcon({
    html: `<div class="price-marker"><span>${display}</span></div>`,
    className: 'price-marker-wrapper',
    iconSize: [52, 32],
    iconAnchor: [26, 32],
  })
}

export default function MapView({ stores, center, userLocation, onBoundsChange, onSelectStore }: MapViewProps) {
  const clusterRef = useRef<Supercluster<ClusterProperties>>(new Supercluster({ radius: 60, maxZoom: 18 }))
  const [bbox, setBbox] = useState<BBox | null>(null)
  const [zoom, setZoom] = useState(13)

  const points = useMemo<ClusterFeature[]>(() => {
    return stores.map((store) => ({
      type: 'Feature',
      properties: {
        cluster: false,
        store,
      },
      geometry: {
        type: 'Point',
        coordinates: [store.longitude, store.latitude],
      },
    }))
  }, [stores])

  useEffect(() => {
    clusterRef.current.load(points)
  }, [points])

  const clusters = useMemo(() => {
    if (!bbox) return points
    return clusterRef.current.getClusters(bbox, zoom)
  }, [bbox, zoom, points])

  const handleBoundsChange = useCallback(
    (nextBbox: BBox, nextZoom: number) => {
      setBbox(nextBbox)
      setZoom(nextZoom)
      onBoundsChange(nextBbox, nextZoom)
    },
    [onBoundsChange],
  )

  return (
    <div className="h-full w-full rounded-2xl overflow-hidden border border-slate-200">
      <MapContainer center={center} zoom={13} className="h-full w-full">
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        <MapEvents onBoundsChange={handleBoundsChange} />
        <MapCenter center={center} />

        {userLocation && (
          <CircleMarker
            center={[userLocation.lat, userLocation.lon]}
            radius={8}
            pathOptions={{ color: '#2563eb', fillColor: '#3b82f6', fillOpacity: 0.8 }}
          />
        )}

        {clusters.map((cluster) => {
          const isCluster = cluster.properties.cluster
          if (isCluster) {
            return (
              <ClusterMarker
                key={`cluster-${cluster.id}`}
                cluster={cluster as ClusterFeature}
                supercluster={clusterRef.current}
              />
            )
          }

          const store = cluster.properties.store
          if (!store) return null
          const icon =
            store.min_price !== undefined && store.min_price !== null
              ? buildPriceIcon(store.min_price)
              : undefined
          return (
            <Marker
              key={`store-${store.id}`}
              position={[store.latitude, store.longitude]}
              icon={icon}
              eventHandlers={{
                click: () => onSelectStore(store),
              }}
            />
          )
        })}
      </MapContainer>
    </div>
  )
}
