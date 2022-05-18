import { screen } from '@testing-library/react'
import { createMemoryHistory } from 'history'

import { renderWithBrandedContext } from '@sourcegraph/shared/src/testing'
import { MockedTestProvider } from '@sourcegraph/shared/src/testing/apollo'

import { ProductSubscriptionForm } from './ProductSubscriptionForm'

jest.mock('../../dotcom/productSubscriptions/features', () => ({
    billingPublishableKey: 'public-key',
}))

jest.mock('@stripe/stripe-js', () => ({
    ...jest.requireActual('@stripe/stripe-js'),
    loadStripe: () =>
        Promise.resolve({
            elements: () => {},
            createToken: () => {},
            createPaymentMethod: () => {},
            confirmCardPayment: () => {},
        }),
}))

jest.mock('@stripe/react-stripe-js', () => ({
    ...jest.requireActual('@stripe/react-stripe-js'),
    CardElement: 'cardelement',
}))

describe('ProductSubscriptionForm', () => {
    test('new subscription for anonymous viewer (no account)', async () => {
        const history = createMemoryHistory()
        const { asFragment } = renderWithBrandedContext(
            <MockedTestProvider>
                <ProductSubscriptionForm
                    accountID={null}
                    subscriptionID={null}
                    onSubmit={() => undefined}
                    submissionState={undefined}
                    primaryButtonText="Submit"
                    isLightTheme={false}
                    history={history}
                />
            </MockedTestProvider>,
            { history }
        )

        expect(await screen.findByText(/submit/i)).toBeInTheDocument()
        expect(asFragment()).toMatchSnapshot()
    })

    test('new subscription for existing account', () => {
        const history = createMemoryHistory()
        expect(
            renderWithBrandedContext(
                <ProductSubscriptionForm
                    accountID="a"
                    subscriptionID={null}
                    onSubmit={() => undefined}
                    submissionState={undefined}
                    primaryButtonText="Submit"
                    isLightTheme={false}
                    history={history}
                />,
                { history }
            ).asFragment()
        ).toMatchSnapshot()
    })

    test('edit existing subscription', () => {
        const history = createMemoryHistory()
        expect(
            renderWithBrandedContext(
                <ProductSubscriptionForm
                    accountID="a"
                    subscriptionID="s"
                    initialValue={{ userCount: 123, billingPlanID: 'p' }}
                    onSubmit={() => undefined}
                    submissionState={undefined}
                    primaryButtonText="Submit"
                    isLightTheme={false}
                    history={history}
                />,
                { history }
            ).asFragment()
        ).toMatchSnapshot()
    })
})
