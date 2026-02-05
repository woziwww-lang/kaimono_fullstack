import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'Price Comparison App | 価格比較アプリ',
  description: 'Compare prices across nearby stores in Japan',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ja">
      <body className="font-sans antialiased">
        {children}
      </body>
    </html>
  )
}
