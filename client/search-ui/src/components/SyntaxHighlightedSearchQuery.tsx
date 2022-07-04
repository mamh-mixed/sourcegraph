import React, { Fragment, ReactElement, useMemo } from 'react'

import classNames from 'classnames'

import { decorate, DecoratedToken } from '@sourcegraph/shared/src/search/query/decoratedToken'
import { scanSearchQuery } from '@sourcegraph/shared/src/search/query/scanner'

interface SyntaxHighlightedSearchQueryProps extends React.HTMLAttributes<HTMLSpanElement> {
    query: string
}

function toElement(query: string, token: DecoratedToken): ReactElement<any, any> {
    switch (token.type) {
        case 'field':
            return (
                <Fragment key={token.range.start}>
                    <span className="search-filter-keyword">{token.value}</span>
                </Fragment>
            )
        case 'keyword':
            return (
                <span className="search-keyword" key={token.range.start}>
                    {token.value}
                </span>
            )
        case 'openingParen':
            return (
                <span className="search-keyword" key={token.range.start}>
                    &#40;
                </span>
            )
        case 'closingParen':
            return (
                <span className="search-keyword" key={token.range.start}>
                    &#41;
                </span>
            )
        case 'metaFilterSeparator':
            return (
                <span className="search-filter-separator" key={token.range.start}>
                    :
                </span>
            )
        case 'metaRepoRevisionSeparator':
        case 'metaContextPrefix':
            return (
                <span className="search-keyword" key={token.range.start}>
                    @
                </span>
            )
        case 'metaPath':
            return (
                <span className="search-path-separator" key={token.range.start}>
                    {token.value}
                </span>
            )
        // case metaRevision // TODO
        case 'metaRegexp': {
            let kind = ''
            switch (token.kind) {
                case 'Assertion':
                    kind = 'assertion'
                    break
                case 'Alternative':
                    kind = 'alternative'
                    break
                case 'Delimited':
                    kind = 'delimited'
                    break
                case 'EscapedCharacter':
                    kind = 'escaped-character'
                    break
                case 'CharacterSet':
                    kind = 'character-set'
                    break
                case 'CharacterClass':
                    kind = 'character-class'
                    break
                case 'CharacterClassRange':
                    kind = 'character-class-range'
                    break
                case 'CharacterClassRangeHyphen':
                    kind = 'character-class-range-hyphen'
                    break
                case 'CharacterClassMember':
                    kind = 'character-class-member'
                    break
                case 'LazyQuantifier':
                    kind = 'lazy-quantifier'
                    break
                case 'RangeQuantifier':
                    kind = 'range-quantifier'
                    break
            }
            return (
                <span className={`search-regexp-meta-${kind}`} key={token.range.start}>
                    {token.value}
                </span>
            )
        }
    }
    return <Fragment key={token.range.start}>{query.slice(token.range.start, token.range.end)}</Fragment>
}

// A read-only syntax highlighted search query
export const SyntaxHighlightedSearchQuery: React.FunctionComponent<
    React.PropsWithChildren<SyntaxHighlightedSearchQueryProps>
> = ({ query, ...otherProps }) => {
    const tokens = useMemo(() => {
        const tokens = scanSearchQuery(query)
        return tokens.type === 'success'
            ? tokens.term.flatMap(token => decorate(token).map(token => toElement(query, token)))
            : [<Fragment key="0">{query}</Fragment>]
    }, [query])

    return (
        <span {...otherProps} className={classNames('text-monospace search-query-link', otherProps.className)}>
            {tokens}
        </span>
    )
}
