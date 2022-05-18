import React from 'react'

import MapSearchIcon from 'mdi-react/MapSearchIcon'

import { Typography } from '@sourcegraph/wildcard'

export const EmptyUploadRetentionMatchStatus: React.FunctionComponent<React.PropsWithChildren<unknown>> = () => (
    <Typography.Text className="text-muted text-center w-100 mb-0 mt-1">
        <MapSearchIcon className="mb-2" />
        <br />
        No retention policies matched.
    </Typography.Text>
)
