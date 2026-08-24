import '@testing-library/jest-dom/vitest'
import { act, cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Calculator } from './components/Calculator'

afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
})

describe('Calculator', () => {
  it('calculates a binary operation', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ result: 35 }), { status: 200 }),
    )
    vi.stubGlobal('fetch', fetchMock)
    const user = userEvent.setup()
    render(<Calculator />)

    await user.click(screen.getByRole('button', { name: 'Multiply' }))
    await user.type(screen.getByLabelText('First number'), '7')
    await user.type(screen.getByLabelText('Second number'), '5')
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(await screen.findByText('35')).toBeInTheDocument()
    expect(fetchMock).toHaveBeenCalledWith('/api/calculate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ operation: 'multiply', left: 7, right: 5 }),
    })
  })

  it('does not round a large result', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ result: 1234567890123 }), { status: 200 }),
    ))
    const user = userEvent.setup()
    render(<Calculator />)

    await user.type(screen.getByLabelText('First number'), '1234567890120')
    await user.type(screen.getByLabelText('Second number'), '3')
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(await screen.findByText('1234567890123')).toBeInTheDocument()
  })

  it('locks value-changing controls while a calculation is pending', async () => {
    let resolveResponse!: (response: Response) => void
    const response = new Promise<Response>((resolve) => {
      resolveResponse = resolve
    })
    vi.stubGlobal('fetch', vi.fn().mockReturnValue(response))
    const user = userEvent.setup()
    render(<Calculator />)

    await user.type(screen.getByLabelText('First number'), '7')
    await user.type(screen.getByLabelText('Second number'), '5')
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(screen.getByLabelText('First number')).toBeDisabled()
    expect(screen.getByLabelText('Second number')).toBeDisabled()
    expect(screen.getByRole('button', { name: 'Subtract' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'Clear' })).toBeDisabled()

    await act(async () => {
      resolveResponse(new Response(JSON.stringify({ result: 12 }), { status: 200 }))
    })

    expect(await screen.findByText('12')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Clear' })).toBeEnabled()
  })

  it('sends only one operand for square root', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ result: 9 }), { status: 200 }),
    )
    vi.stubGlobal('fetch', fetchMock)
    const user = userEvent.setup()
    render(<Calculator />)

    await user.click(screen.getByRole('button', { name: 'Square root' }))
    await user.type(screen.getByLabelText('Number'), '81')
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(await screen.findByText('9')).toBeInTheDocument()
    expect(fetchMock).toHaveBeenCalledWith('/api/calculate', expect.objectContaining({
      body: JSON.stringify({ operation: 'square_root', left: 81 }),
    }))
  })

  it('validates empty fields before calling the API', async () => {
    const fetchMock = vi.fn()
    vi.stubGlobal('fetch', fetchMock)
    const user = userEvent.setup()
    render(<Calculator />)

    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(screen.getByRole('alert')).toHaveTextContent('Enter every required number.')
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('shows an API error', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ error: { message: 'Cannot divide by zero.' } }), { status: 400 }),
    ))
    const user = userEvent.setup()
    render(<Calculator />)

    await user.click(screen.getByRole('button', { name: 'Divide' }))
    await user.type(screen.getByLabelText('First number'), '7')
    await user.type(screen.getByLabelText('Second number'), '0')
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(await screen.findByRole('alert')).toHaveTextContent('Cannot divide by zero.')
  })
})
