import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'

// Simple example component for testing
function Welcome({ name }: { name: string }) {
  return <h1>Welcome, {name}!</h1>
}

describe('Example Test', () => {
  it('should render welcome message', () => {
    render(<Welcome name="価格比較アプリ" />)
    expect(screen.getByText('Welcome, 価格比較アプリ!')).toBeInTheDocument()
  })

  it('should perform basic math', () => {
    expect(1 + 1).toBe(2)
  })
})
