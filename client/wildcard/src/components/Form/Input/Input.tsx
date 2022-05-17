import { useRef, forwardRef, InputHTMLAttributes, ReactNode } from 'react'

import classNames from 'classnames'
import { useMergeRefs } from 'use-callback-ref'

import { useAutoFocus } from '../../../hooks/useAutoFocus'
import { ForwardReferenceComponent } from '../../../types'
import { Label } from '../../Typography/Label'

import styles from './Input.module.scss'

export enum InputStatus {
    initial = 'initial',
    error = 'error',
    valid = 'valid',
}

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
    /** text label of input. */
    label?: ReactNode
    /** Description block shown below the input. */
    message?: ReactNode
    /** Custom class name for root label element. */
    className?: string
    /** Custom class name for input element. */
    inputClassName?: string
    /** Exclusive status */
    status?: InputStatus | `${InputStatus}`
    error?: ReactNode
    /** Disable input behavior */
    disabled?: boolean
    /** Determines the size of the input */
    variant?: 'regular' | 'small'
}

/**
 * Displays the input with description, error message, visual invalid and valid states.
 * Does not support Loader icon and status=loadind (user FormInput to get support for loading state)
 */
export const Input = forwardRef((props, reference) => {
    const {
        as: Component = 'input',
        type = 'text',
        variant = 'regular',
        label,
        message,
        className,
        inputClassName,
        disabled,
        status = InputStatus.initial,
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
            <Component
                disabled={disabled}
                type={type}
                className={classNames(styles.input, inputClassName, 'form-control', 'with-invalid-icon', {
                    'is-valid': status === InputStatus.valid,
                    'is-invalid': error || status === InputStatus.error,
                    'form-control-sm': variant === 'small',
                })}
                {...otherProps}
                ref={mergedReference}
                autoFocus={autoFocus}
            />

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
}) as ForwardReferenceComponent<'input', InputProps>

Input.displayName = 'Input'
