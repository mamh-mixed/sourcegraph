import React from 'react'

import { Typography } from '@sourcegraph/wildcard'

import styles from './CloseChangesetsListEmptyElement.module.scss'

export const CloseChangesetsListEmptyElement: React.FunctionComponent<React.PropsWithChildren<{}>> = () => (
    <div className={styles.closeChangesetsListEmptyElementBody}>
        <Typography.Text className="text-center text-muted font-weight-normal">
            Closing this batch change will not alter changesets and no changesets will remain open.
        </Typography.Text>
    </div>
)
