import { useRef, forwardRef, ReactNode } from 'react'

import classNames from 'classnames'
import { useMergeRefs } from 'use-callback-ref'

import { LoaderInput } from '@sourcegraph/branded/src/components/LoaderInput'

import { useAutoFocus } from '../../../hooks/useAutoFocus'
import { ForwardReferenceComponent } from '../../../types'
import { Label } from '../../Typography/Label'

import { InputProps } from './Input'

import styles from './Input.module.scss'

export enum FormInputStatus {
    initial = 'initial',
    error = 'error',
    loading = 'loading',
    valid = 'valid',
}

export interface FormInputProps extends Omit<InputProps, 'status'> {
    status?: FormInputStatus | `${FormInputStatus}`
    /** Input icon (symbol) which render right after the input element. */
    inputSymbol?: ReactNode
}

/**
 * Displays the input with description, error message, visual invalid and valid states.
 * Renders Input component within LoaderInput to display loader icon with status=loading.
 */
export const FormInput = forwardRef((props, reference) => {
    const {
        as: Component = 'input',
        type = 'text',
        variant = 'regular',
        label,
        message,
        className,
        inputClassName,
        inputSymbol,
        disabled,
        status = FormInputStatus.initial,
        error,
        autoFocus,
        ...otherProps
    } = props

    const localReference = useRef<HTMLInputElement>(null)
    const mergedReference = useMergeRefs([localReference, reference])

    useAutoFocus({ autoFocus, reference: localReference })

    const messageClassName = 'form-text font-weight-normal mt-2'
    const inputWithMessage = (
        <>
            <LoaderInput
                className={classNames('d-flex', !label && className)}
                loading={status === FormInputStatus.loading}
            >
                <Component
                    disabled={disabled}
                    type={type}
                    className={classNames(styles.input, inputClassName, 'form-control', 'with-invalid-icon', {
                        'is-valid': status === FormInputStatus.valid,
                        'is-invalid': error || status === FormInputStatus.error,
                        'form-control-sm': variant === 'small',
                    })}
                    {...otherProps}
                    ref={mergedReference}
                    autoFocus={autoFocus}
                />

                {inputSymbol}
            </LoaderInput>

            {error && (
                <small role="alert" className={classNames('text-danger', messageClassName)}>
                    {error}
                </small>
            )}
            {!error && message && <small className={classNames('text-muted', messageClassName)}>{message}</small>}
        </>
    )

    if (label) {
        return (
            <Label className={classNames('w-100', className)}>
                {label && <div className="mb-2">{variant === 'regular' ? label : <small>{label}</small>}</div>}
                {inputWithMessage}
            </Label>
        )
    }

    return inputWithMessage
}) as ForwardReferenceComponent<'input', FormInputProps>

FormInput.displayName = 'FormInput'
