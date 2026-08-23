import { useState } from 'react'
import type { FormEvent } from 'react'
import { calculate } from './api'
import type { BinaryOperation, Operation } from './api'

const operations: Array<{ value: Operation; symbol: string; label: string }> = [
  { value: 'add', symbol: '+', label: 'Add' },
  { value: 'subtract', symbol: '−', label: 'Subtract' },
  { value: 'multiply', symbol: '×', label: 'Multiply' },
  { value: 'divide', symbol: '÷', label: 'Divide' },
  { value: 'power', symbol: 'xʸ', label: 'Power' },
  { value: 'square_root', symbol: '√', label: 'Square root' },
  { value: 'percentage', symbol: '%', label: 'Percentage' },
]

export function Calculator() {
  const [operation, setOperation] = useState<Operation>('add')
  const [left, setLeft] = useState('')
  const [right, setRight] = useState('')
  const [result, setResult] = useState<number | null>(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const selected = operations.find((item) => item.value === operation) ?? operations[0]
  const unary = operation === 'square_root'

  function selectOperation(next: Operation) {
    setOperation(next)
    setResult(null)
    setError('')
    if (next === 'square_root') setRight('')
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError('')
    setResult(null)

    if (left.trim() === '' || (!unary && right.trim() === '')) {
      setError('Enter every required number.')
      return
    }

    const leftNumber = Number(left)
    const rightNumber = Number(right)
    if (!Number.isFinite(leftNumber) || (!unary && !Number.isFinite(rightNumber))) {
      setError('Enter valid numbers.')
      return
    }

    setLoading(true)
    try {
      const input = unary
        ? { operation, left: leftNumber }
        : { operation: operation as BinaryOperation, left: leftNumber, right: rightNumber }
      setResult(await calculate(input))
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : 'The calculation failed.')
    } finally {
      setLoading(false)
    }
  }

  function clear() {
    setLeft('')
    setRight('')
    setResult(null)
    setError('')
  }

  return (
    <main className="page">
      <section className="calculator" aria-labelledby="calculator-title">
        <h1 id="calculator-title">Calculator</h1>

        <div className="operation-grid" role="group" aria-label="Operation">
          {operations.map((item) => (
            <button
              className="operation"
              type="button"
              key={item.value}
              aria-label={item.label}
              aria-pressed={operation === item.value}
              onClick={() => selectOperation(item.value)}
            >
              {item.symbol}
            </button>
          ))}
        </div>

        <form onSubmit={submit} noValidate>
          <div className={unary ? 'fields unary' : 'fields'}>
            <label>
              <span>{unary ? 'Number' : operation === 'percentage' ? 'Percentage' : 'First number'}</span>
              <input
                type="number"
                inputMode="decimal"
                step="any"
                value={left}
                onChange={(event) => setLeft(event.target.value)}
                autoFocus
              />
            </label>

            {!unary && (
              <>
                <output className="operator" aria-label={selected.label}>
                  {selected.symbol}
                </output>
                <label>
                  <span>{operation === 'percentage' ? 'Of number' : 'Second number'}</span>
                  <input
                    type="number"
                    inputMode="decimal"
                    step="any"
                    value={right}
                    onChange={(event) => setRight(event.target.value)}
                  />
                </label>
              </>
            )}
          </div>

          <div className="actions">
            <button className="primary" type="submit" disabled={loading}>
              {loading ? 'Calculating…' : 'Calculate'}
            </button>
            <button className="secondary" type="button" onClick={clear}>
              Clear
            </button>
          </div>
        </form>

        <div className="feedback" aria-live="polite">
          {result !== null && (
            <p>
              <span>Result</span>
              <strong>{result.toString()}</strong>
            </p>
          )}
          {error && <p role="alert">{error}</p>}
        </div>
      </section>
    </main>
  )
}
