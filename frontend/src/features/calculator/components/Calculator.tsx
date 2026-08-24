import type { FormEvent } from 'react'
import { useCalculator } from '../useCalculator'
import { OperandFields } from './OperandFields'
import { OperationSelector } from './OperationSelector'

export function Calculator() {
  const calculator = useCalculator()

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    void calculator.submit()
  }

  return (
    <section className="calculator" aria-labelledby="calculator-title">
      <h1 id="calculator-title">Calculator</h1>

      <OperationSelector
        value={calculator.operation}
        disabled={calculator.loading}
        onChange={calculator.selectOperation}
      />

      <form onSubmit={submit} noValidate>
        <OperandFields
          operation={calculator.operation}
          left={calculator.left}
          right={calculator.right}
          disabled={calculator.loading}
          onLeftChange={calculator.setLeft}
          onRightChange={calculator.setRight}
        />

        <div className="actions">
          <button className="primary" type="submit" disabled={calculator.loading}>
            {calculator.loading ? 'Calculating…' : 'Calculate'}
          </button>
          <button
            className="secondary"
            type="button"
            disabled={calculator.loading}
            onClick={calculator.clear}
          >
            Clear
          </button>
        </div>
      </form>

      <div className="feedback" aria-live="polite">
        {calculator.result !== null && (
          <p>
            <span>Result</span>
            <strong>{calculator.result.toString()}</strong>
          </p>
        )}
        {calculator.error && <p role="alert">{calculator.error}</p>}
      </div>
    </section>
  )
}
