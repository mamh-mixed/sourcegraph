import React from 'react'

import { AnchorLink } from '../AnchorLink'

// These should no longer be needed once we update to React Router v6 and can import correct types
// from the newer version of history through react-router. In the current version of history,
// type To isn't exported and this causes issues with any down the type chain.
// Copied from: https://github.com/remix-run/history/blob/dev/packages/history/index.ts
export type TEMP_To = string | Partial<TEMP_Path>
export interface TEMP_Path {
    pathname: string
    search: string
    hash: string
}

export interface LinkProps
    extends Pick<
        React.AnchorHTMLAttributes<HTMLAnchorElement>,
        Exclude<keyof React.AnchorHTMLAttributes<HTMLAnchorElement>, 'href'>
    > {
    to: TEMP_To
    ref?: React.Ref<HTMLAnchorElement>
}

/**
 * The component used to render a link. All shared code must use this component for links—not <a>, <Link>, etc.
 *
 * Different platforms (web app vs. browser extension) require the use of different link components:
 *
 * The web app uses <RouterLinkOrAnchor>, which uses react-router-dom's <Link> for relative URLs (for page
 * navigation using the HTML history API) and <a> for absolute URLs. The react-router-dom <Link> component only
 * works inside a react-router <BrowserRouter> context, so it wouldn't work in the browser extension.
 *
 * The browser extension uses <a> for everything (because code hosts don't generally use react-router). A
 * react-router-dom <Link> wouldn't work in the browser extension, because there is no <BrowserRouter>.
 *
 * This variable must be set at initialization time by calling {@link setLinkComponent}.
 *
 * The `to` property holds the destination URL (do not use `href`). If <a> is used, the `to` property value is
 * given as the `href` property value on the <a> element.
 *
 * @see setLinkComponent
 */
export let Link: typeof AnchorLink

if (process.env.NODE_ENV !== 'production') {
    // Fail with helpful message if setLinkComponent has not been called when the <Link> component is used.

    // eslint-disable-next-line react/display-name
    Link = React.forwardRef(() => {
        throw new Error('No Link component set. You must call setLinkComponent to set the Link component to use.')
    }) as typeof Link
}

/**
 * Sets (globally) the component to use for links. This must be set at initialization time.
 *
 * @see Link
 * @see AnchorLink
 */
export function setLinkComponent(component: typeof Link): void {
    Link = component
}
