import { useEffect, useState } from 'react'

import { Elements } from '@stripe/react-stripe-js'
import { loadStripe, Stripe } from '@stripe/stripe-js'

import { billingPublishableKey } from '../productSubscriptions/features'

/**
 * Wraps a React tree (of elements) and injects the Stripe API.
 */
export const StripeWrapper: React.FunctionComponent<React.PropsWithChildren<{}>> = ({ children }) => {
    const [stripe, setStripe] = useState<Stripe | null>(null)

    useEffect(() => {
        if (stripe || !billingPublishableKey) {
            return
        }

        // eslint-disable-next-line @typescript-eslint/no-floating-promises
        loadStripe(billingPublishableKey).then(loadedStripe => {
            setStripe(loadedStripe)
        })
    }, [stripe])

    if (!stripe) {
        return null
    }

    return <Elements stripe={stripe}>{children}</Elements>
}
