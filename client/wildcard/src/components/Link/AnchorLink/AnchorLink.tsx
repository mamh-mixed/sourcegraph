import * as React from 'react'

import classNames from 'classnames'
import { Link, useHref } from 'react-router-dom-v5-compat'

import { useWildcardTheme } from '../../../hooks/useWildcardTheme'
import { ForwardReferenceComponent } from '../../../types'
import type { LinkProps } from '../Link'

import styles from './AnchorLink.module.scss'

type LinkType = typeof Link
export type AnchorLinkProps = LinkProps

export const AnchorLink = React.forwardRef(({ as: Component, children, className, ...rest }, reference) => {
    const { isBranded } = useWildcardTheme()

    const commonProps = {
        ref: reference,
        className: classNames(isBranded && styles.anchorLink, className),
    }

    if (!Component) {
        return <PlainLink {...rest} {...commonProps} />
    }

    return (
        <Component {...rest} {...commonProps}>
            {children}
        </Component>
    )
}) as ForwardReferenceComponent<LinkType, AnchorLinkProps>

const PlainLink = React.forwardRef(({ children, to, className, ...rest }, reference) => {
    const href = useHref(to)
    return (
        // eslint-disable-next-line react/forbid-elements
        <a href={href} className={className} ref={reference} {...rest}>
            {children}
        </a>
    )
}) as ForwardReferenceComponent<LinkType, AnchorLinkProps>

AnchorLink.displayName = 'AnchorLink'
