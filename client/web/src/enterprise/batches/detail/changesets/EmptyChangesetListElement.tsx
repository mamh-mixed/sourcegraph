import React from 'react'

import { Typography } from '@sourcegraph/wildcard'

import styles from './EmptyChangesetListElement.module.scss'

export const EmptyChangesetListElement: React.FunctionComponent<React.PropsWithChildren<{}>> = () => (
    <div className={styles.emptyChangesetListElementBody}>
        <Typography.H2 className="text-center mb-4">This batch change does not contain changesets</Typography.H2>
        <Typography.Text>This can occur for several reasons:</Typography.Text>
        <Typography.Text>
            <strong>
                The query specified in <span className="text-monospace">repositoriesMatchingQuery:</span> may not have
                matched any repositories.
            </strong>
        </Typography.Text>
        <Typography.Text>Test your query in the search bar and ensure it returns results.</Typography.Text>
        <Typography.Text>
            <strong>
                The code specified in <span className="text-monospace">steps:</span> may not have resulted in changes
                being made.
            </strong>
        </Typography.Text>
        <Typography.Text>
            Try the command on a local instance of one of the repositories returned in your search results. Run{' '}
            <span className="text-monospace">git status</span> and ensure it produced changed files.
        </Typography.Text>
    </div>
)
