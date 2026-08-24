import type { BinaryOperation } from './model'

export type CalculationRequest =
  | { operation: BinaryOperation; left: number; right: number }
  | { operation: 'square_root'; left: number }

type CalculationResponse = { result: number }
type ErrorResponse = { error?: { message?: string } }

export async function calculate(input: CalculationRequest): Promise<number> {
  const response = await fetch('/api/calculate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const body = (await response.json().catch(() => null)) as
    | (CalculationResponse & ErrorResponse)
    | null

  if (!response.ok) {
    throw new Error(body?.error?.message ?? 'The calculation failed.')
  }
  if (!body || !Number.isFinite(body.result)) {
    throw new Error('The server returned an invalid result.')
  }

  return body.result
}
