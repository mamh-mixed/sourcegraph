import * as React from 'react'

import classNames from 'classnames'
import { Link, useHref } from 'react-router-dom-v5-compat'

import { useWildcardTheme } from '../../../hooks/useWildcardTheme'
import { ForwardReferenceComponent } from '../../../types'
import type { LinkProps } from '../Link'

import styles from './AnchorLink.module.scss'

type LinkType = typeof Link
export type AnchorLinkProps = LinkProps

// eslint-disable-next-line react/display-name
export const AnchorLink = React.forwardRef(({ as: Component, children, className, to, ...rest }, reference) => {
    const { isBranded } = useWildcardTheme()

    const commonProps = {
        ref: reference,
        className: classNames(isBranded && styles.anchorLink, className),
    }

    if (!Component) {
        // We may be able to get rid of this branching if we start always rendering
        // Link inside Router. Now it's not the case because for
        // tour (onboarding) components, it's rendered by renderBrandedToString
        if (typeof to === 'string') {
            return (
                <PlainLinkWithHref to={to} {...rest} {...commonProps}>
                    {children}
                </PlainLinkWithHref>
            )
        }
        return (
            <PlainLinkWithTo to={to} {...rest} {...commonProps}>
                {children}
            </PlainLinkWithTo>
        )
    }

    return (
        <Component to={to} {...rest} {...commonProps}>
            {children}
        </Component>
    )
}) as ForwardReferenceComponent<LinkType, AnchorLinkProps>
AnchorLink.displayName = 'AnchorLink'

// eslint-disable-next-line react/display-name
const PlainLinkWithTo = React.forwardRef(({ children, to, className, ...rest }, reference) => {
    const href = useHref(to)
    return (
        // eslint-disable-next-line react/forbid-elements
        <a href={href} className={className} ref={reference} {...rest}>
            {children}
        </a>
    )
}) as ForwardReferenceComponent<LinkType, AnchorLinkProps>
PlainLinkWithTo.displayName = 'PlainLinkWithTo'

// eslint-disable-next-line react/display-name
const PlainLinkWithHref = React.forwardRef(({ children, to, className, ...rest }, reference) => (
    // eslint-disable-next-line react/forbid-elements
    <a href={to} className={className} ref={reference} {...rest}>
        {children}
    </a>
)) as ForwardReferenceComponent<LinkType, AnchorLinkProps & { to: string }>
PlainLinkWithHref.displayName = 'PlainLinkWithHref'
