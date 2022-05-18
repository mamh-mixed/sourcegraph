import * as React from 'react'

import { CardElement } from '@stripe/react-stripe-js'
import classNames from 'classnames'

import { ThemeProps } from '@sourcegraph/shared/src/theme'

import styles from './PaymentTokenFormControl.module.scss'

interface Props extends ThemeProps {
    disabled?: boolean
}

/**
 * Displays a payment form control for the user to enter payment information to purchase a product subscription.
 */
export const PaymentTokenFormControl: React.FunctionComponent<React.PropsWithChildren<Props>> = props => {
    const textColor = props.isLightTheme ? '#2b3750' : 'white'

    return (
        <div className="payment-token-form-control">
            <CardElement
                className={classNames('form-control', styles.card, props.disabled && styles.cardDisabled)}
                options={{
                    style: {
                        base: {
                            fontFamily:
                                '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',
                            color: textColor,
                            ':-webkit-autofill': {
                                color: textColor,
                            },
                        },
                    },
                    disabled: props.disabled,
                }}
            />
        </div>
    )
}
