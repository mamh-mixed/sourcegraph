import { storiesOf } from '@storybook/react'

import { BrandedStory } from '@sourcegraph/branded/src/components/BrandedStory'
import { Typography } from '@sourcegraph/wildcard'

import { SyntaxHighlightedSearchQuery } from './SyntaxHighlightedSearchQuery'

const { add } = storiesOf('search-ui/SyntaxHighlightedSearchQuery', module).addParameters({
    chromatic: { viewports: [480] },
})

add('SyntaxHighlightedSearchQuery', () => (
    <BrandedStory>
        {() => (
            <Typography.Text>
                <SyntaxHighlightedSearchQuery query="test AND spec" />
                <br />
                <SyntaxHighlightedSearchQuery query="test or spec repo:sourcegraph" />
                <br />
                <SyntaxHighlightedSearchQuery query="test -lang:ts" />
            </Typography.Text>
        )}
    </BrandedStory>
))
