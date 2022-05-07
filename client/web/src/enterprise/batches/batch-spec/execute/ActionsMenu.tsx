import React, { useCallback, useState } from 'react'

import ChevronDownIcon from 'mdi-react/ChevronDownIcon'
import CloseIcon from 'mdi-react/CloseIcon'
import PencilIcon from 'mdi-react/PencilIcon'
import SyncIcon from 'mdi-react/SyncIcon'
import { useHistory, useLocation } from 'react-router'

import { useMutation } from '@sourcegraph/http-client'
import { Button, Icon, Link, Menu, MenuButton, MenuItem, MenuList, Position, useMeasure } from '@sourcegraph/wildcard'

import {
    BatchSpecExecutionFields,
    BatchSpecState,
    CancelBatchSpecExecutionResult,
    CancelBatchSpecExecutionVariables,
    RetryBatchSpecExecutionResult,
    RetryBatchSpecExecutionVariables,
} from '../../../../graphql-operations'
import { useBatchSpecContext } from '../BatchSpecContext'

import { CANCEL_BATCH_SPEC_EXECUTION, RETRY_BATCH_SPEC_EXECUTION } from './backend'
import { CancelExecutionModal } from './CancelExecutionModal'

import styles from './ActionsMenu.module.scss'

export const ActionsMenu: React.FunctionComponent<React.PropsWithChildren<{}>> = () => {
    const history = useHistory()
    const location = useLocation()

    const { batchChange, batchSpec, setActionsError } = useBatchSpecContext<BatchSpecExecutionFields>()
    const { url } = batchChange
    const { isExecuting } = batchSpec

    const [showCancelModal, setShowCancelModal] = useState(false)
    const [cancelModalType, setCancelModalType] = useState<'cancel' | 'edit'>('cancel')
    const [cancelBatchSpecExecution, { loading: isCancelLoading }] = useMutation<
        CancelBatchSpecExecutionResult,
        CancelBatchSpecExecutionVariables
    >(CANCEL_BATCH_SPEC_EXECUTION, {
        variables: { id: batchSpec.id },
        onError: setActionsError,
        onCompleted: () => {
            setShowCancelModal(false)
            history.push(`${url}/edit`)
        },
    })

    const [retryBatchSpecExecution, { loading: isRetryLoading }] = useMutation<
        RetryBatchSpecExecutionResult,
        RetryBatchSpecExecutionVariables
    >(RETRY_BATCH_SPEC_EXECUTION, { variables: { id: batchSpec.id }, onError: setActionsError })

    const onSelectEdit = useCallback(() => {
        if (isExecuting) {
            setCancelModalType('edit')
            setShowCancelModal(true)
        } else {
            history.push(`${url}/edit`)
        }
    }, [isExecuting, url, history])

    const onSelectCancel = useCallback(() => {
        setCancelModalType('cancel')
        setShowCancelModal(true)
    }, [])

    const showPreviewButton = !location.pathname.endsWith('preview') && !!batchSpec.applyURL

    // The actions menu button is wider than the "Preview" button, so to prevent layout
    // shift, we apply the width of the actions menu button to the "Preview" button
    // instead.
    const [menuReference, { width: menuWidth }] = useMeasure()

    return (
        <div className="relative">
            {showPreviewButton && (
                <Button
                    to={`${batchSpec.executionURL}/preview`}
                    variant="primary"
                    as={Link}
                    className={styles.previewButton}
                    style={{ width: menuWidth }}
                >
                    Preview
                </Button>
            )}
            <Menu>
                <div ref={menuReference} aria-hidden={showPreviewButton}>
                    <MenuButton variant="secondary" className={showPreviewButton ? styles.menuButtonHidden : undefined}>
                        Actions
                        <Icon as={ChevronDownIcon} className={styles.chevronIcon} />
                    </MenuButton>
                </div>
                <MenuList position={Position.bottomEnd}>
                    <MenuItem onSelect={onSelectEdit}>
                        <Icon as={PencilIcon} /> Edit spec{isExecuting ? '...' : ''}
                    </MenuItem>
                    {isExecuting && (
                        <MenuItem onSelect={onSelectCancel}>
                            <Icon as={CloseIcon} className={styles.cancelIcon} /> Cancel execution...
                        </MenuItem>
                    )}
                    {batchSpec.state !== BatchSpecState.COMPLETED && batchSpec.viewerCanRetry && (
                        <MenuItem onSelect={retryBatchSpecExecution} disabled={isRetryLoading}>
                            <Icon as={SyncIcon} /> Retry failed workspaces
                        </MenuItem>
                    )}
                </MenuList>
            </Menu>
            <CancelExecutionModal
                isOpen={showCancelModal}
                onCancel={() => setShowCancelModal(false)}
                onConfirm={cancelBatchSpecExecution}
                modalHeader={cancelModalType === 'cancel' ? 'Cancel execution' : 'The execution is still running'}
                modalBody={
                    <p>
                        {cancelModalType === 'cancel'
                            ? 'Are you sure you want to cancel the current execution?'
                            : 'You are unable to edit the spec when an execution is running.'}
                    </p>
                }
                isLoading={isCancelLoading}
            />
        </div>
    )
}
