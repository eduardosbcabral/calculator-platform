import { useState } from 'react'
import { calculate } from './api'
import type { CalculationRequest } from './api'
import { isUnaryOperation } from './model'
import type { Operation } from './model'

export function useCalculator() {
  const [operation, setOperation] = useState<Operation>('add')
  const [left, setLeft] = useState('')
  const [right, setRight] = useState('')
  const [result, setResult] = useState<number | null>(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  function selectOperation(next: Operation) {
    setOperation(next)
    setResult(null)
    setError('')
    if (isUnaryOperation(next)) setRight('')
  }

  async function submit() {
    setError('')
    setResult(null)

    const unary = isUnaryOperation(operation)
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

    const input: CalculationRequest = unary
      ? { operation, left: leftNumber }
      : { operation, left: leftNumber, right: rightNumber }

    setLoading(true)
    try {
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

  return {
    operation,
    left,
    right,
    result,
    error,
    loading,
    setLeft,
    setRight,
    selectOperation,
    submit,
    clear,
  }
}
