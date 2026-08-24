export type BinaryOperation =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'power'
  | 'percentage'

export type Operation = BinaryOperation | 'square_root'

export type OperationOption = {
  value: Operation
  symbol: string
  label: string
}

export const operations: readonly OperationOption[] = [
  { value: 'add', symbol: '+', label: 'Add' },
  { value: 'subtract', symbol: '−', label: 'Subtract' },
  { value: 'multiply', symbol: '×', label: 'Multiply' },
  { value: 'divide', symbol: '÷', label: 'Divide' },
  { value: 'power', symbol: 'xʸ', label: 'Power' },
  { value: 'square_root', symbol: '√', label: 'Square root' },
  { value: 'percentage', symbol: '%', label: 'Percentage' },
]

export function isUnaryOperation(operation: Operation): operation is 'square_root' {
  return operation === 'square_root'
}
