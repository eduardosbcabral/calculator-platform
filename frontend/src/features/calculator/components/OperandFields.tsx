import { isUnaryOperation, operations } from '../model'
import type { Operation } from '../model'

type Props = {
  operation: Operation
  left: string
  right: string
  disabled: boolean
  onLeftChange: (value: string) => void
  onRightChange: (value: string) => void
}

export function OperandFields({
  operation,
  left,
  right,
  disabled,
  onLeftChange,
  onRightChange,
}: Props) {
  const unary = isUnaryOperation(operation)
  const selected = operations.find((item) => item.value === operation) ?? operations[0]

  return (
    <div className={unary ? 'fields unary' : 'fields'}>
      <label>
        <span>{unary ? 'Number' : operation === 'percentage' ? 'Percentage' : 'First number'}</span>
        <input
          type="number"
          inputMode="decimal"
          step="any"
          value={left}
          disabled={disabled}
          onChange={(event) => onLeftChange(event.target.value)}
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
              disabled={disabled}
              onChange={(event) => onRightChange(event.target.value)}
            />
          </label>
        </>
      )}
    </div>
  )
}
