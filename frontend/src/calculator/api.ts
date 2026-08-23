export type BinaryOperation =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'power'
  | 'percentage'

export type Operation = BinaryOperation | 'square_root'

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
  const body = (await response.json()) as CalculationResponse & ErrorResponse

  if (!response.ok) {
    throw new Error(body.error?.message ?? 'The calculation failed.')
  }
  if (!Number.isFinite(body.result)) {
    throw new Error('The server returned an invalid result.')
  }

  return body.result
}
