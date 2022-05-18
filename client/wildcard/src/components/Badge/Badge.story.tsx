import { Meta } from '@storybook/react'

import { BrandedStory } from '@sourcegraph/branded/src/components/BrandedStory'
import webStyles from '@sourcegraph/web/src/SourcegraphWebApp.scss'

import { Typography } from '..'

import { BADGE_VARIANTS } from './constants'

import { Badge } from '.'

const config: Meta = {
    title: 'wildcard/Badge',

    decorators: [
        story => (
            <BrandedStory styles={webStyles}>{() => <div className="container mt-3">{story()}</div>}</BrandedStory>
        ),
    ],

    parameters: {
        component: Badge,
        chromatic: {
            enableDarkMode: true,
            disableSnapshot: false,
        },
        design: [
            {
                type: 'figma',
                name: 'Figma Light',
                url: 'https://www.figma.com/file/NIsN34NH7lPu04olBzddTw/Wildcard-Design-System?node-id=908%3A6149',
            },
            {
                type: 'figma',
                name: 'Figma Dark',
                url: 'https://www.figma.com/file/NIsN34NH7lPu04olBzddTw/Wildcard-Design-System?node-id=908%3A6447',
            },
        ],
    },
}

export default config

export const Badges = () => (
    <>
        <Typography.H1>Badges</Typography.H1>
        <Typography.H2>Variants</Typography.H2>
        <Typography.Text>Our badges can be styled with different variants.</Typography.Text>
        {BADGE_VARIANTS.map(variant => (
            <Badge key={variant} variant={variant} className="mr-2">
                {variant}
            </Badge>
        ))}
        <Typography.H2 className="mt-4">Size</Typography.H2>
        <Typography.Text>We can also make our badges smaller.</Typography.Text>
        <Badge small={true}>I am a small badge</Badge>
        <Typography.H2 className="mt-4">Pills</Typography.H2>
        <Typography.Text>Commonly used to display counts, we can style badges as pills.</Typography.Text>
        <Badge pill={true} variant="secondary">
            321+
        </Badge>
        <Typography.H2 className="mt-4">Links</Typography.H2>
        <Typography.Text>For more advanced functionality, badges can also function as links.</Typography.Text>
        <Badge href="https://example.com" variant="secondary">
            I am a link
        </Badge>
    </>
)
