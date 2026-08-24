import { operations } from '../model'
import type { Operation } from '../model'

type Props = {
  value: Operation
  disabled: boolean
  onChange: (operation: Operation) => void
}

export function OperationSelector({ value, disabled, onChange }: Props) {
  return (
    <div className="operation-grid" role="group" aria-label="Operation">
      {operations.map((operation) => (
        <button
          className="operation"
          type="button"
          key={operation.value}
          aria-label={operation.label}
          aria-pressed={value === operation.value}
          disabled={disabled}
          onClick={() => onChange(operation.value)}
        >
          {operation.symbol}
        </button>
      ))}
    </div>
  )
}
