import { FormInputStatus } from '..'

interface GetInputStatusProps {
    isValid?: boolean
    isError?: boolean
    isLoading?: boolean
}

export function getInputStatus(props: GetInputStatusProps): FormInputStatus {
    const { isLoading, isError, isValid } = props

    if (isLoading) {
        return FormInputStatus.loading
    }

    if (isError) {
        return FormInputStatus.error
    }

    if (isValid) {
        return FormInputStatus.valid
    }

    return FormInputStatus.initial
}
