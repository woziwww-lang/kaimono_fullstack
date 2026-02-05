import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import React from 'react'
import Home from '../app/page'

vi.mock('next/navigation', () => ({
  useRouter: () => ({
    replace: vi.fn(),
  }),
  useSearchParams: () => new URLSearchParams(),
}))

vi.mock('../components/MapView', () => ({
  default: (props: { onBoundsChange: (bbox: number[], zoom: number) => void }) => {
    React.useEffect(() => {
      props.onBoundsChange([139.0, 35.0, 139.1, 35.1], 13)
    }, [props])
    return <div data-testid="map" />
  },
}))

describe('Home page', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', (input: RequestInfo) => {
      const url = input.toString()

      if (url.includes('/api/stores')) {
        return Promise.resolve({
          ok: true,
          json: async () => ({
            data: [
              {
                id: 1,
                name: 'テスト店舗',
                address: '東京都',
                phone: '03-0000-0000',
                latitude: 35.0,
                longitude: 139.0,
                created_at: '2024-01-01T00:00:00Z',
                updated_at: '2024-01-01T00:00:00Z',
              },
            ],
          }),
        })
      }

      return Promise.resolve({ ok: false, json: async () => ({}) })
    })
  })

  it('renders stores from API', async () => {
    render(<Home />)

    expect(await screen.findByText('テスト店舗')).toBeInTheDocument()
  })
})
