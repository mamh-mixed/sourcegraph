import React, { useCallback, useContext, useMemo, useState } from 'react'

import AlertCircleIcon from 'mdi-react/AlertCircleIcon'

import { useQuery } from '@sourcegraph/http-client'
import { Settings, SettingsCascadeProps } from '@sourcegraph/shared/src/settings/settings'
import { ThemeProps } from '@sourcegraph/shared/src/theme'
import { Icon, LoadingSpinner } from '@sourcegraph/wildcard'

import { HeroPage } from '../../../../components/HeroPage'
import { GetBatchChangeToEditResult, GetBatchChangeToEditVariables, Scalars } from '../../../../graphql-operations'
// TODO: Move some of these to batch-spec/edit
import { GET_BATCH_CHANGE_TO_EDIT } from '../../create/backend'
import { ConfigurationForm } from '../../create/ConfigurationForm'
import { ExecutionOptions, ExecutionOptionsDropdown } from '../../create/ExecutionOptions'
import { InsightTemplatesBanner } from '../../create/InsightTemplatesBanner'
import { useInsightTemplates } from '../../create/useInsightTemplates'
import { BatchSpecContext, BatchSpecContextProvider } from '../BatchSpecContext'
import { ActionButtons } from '../header/ActionButtons'
import { BatchChangeHeader } from '../header/BatchChangeHeader'
import { TabBar, TabsConfig, TabName } from '../TabBar'

import { EditorFeedbackPanel } from './editor/EditorFeedbackPanel'
import { MonacoBatchSpecEditor } from './editor/MonacoBatchSpecEditor'
import { LibraryPane } from './library/LibraryPane'
import { WorkspacesPreviewPanel } from './workspaces-preview/WorkspacesPreviewPanel'

import layoutStyles from '../Layout.module.scss'
import styles from './EditBatchSpecPage.module.scss'

export interface EditBatchSpecPageProps extends SettingsCascadeProps<Settings>, ThemeProps {
    batchChange: { name: string; url: string; namespace: { id: Scalars['ID'] } }
}

export const EditBatchSpecPage: React.FunctionComponent<EditBatchSpecPageProps> = ({ batchChange, ...props }) => {
    const variables = useMemo(() => ({ namespace: batchChange.namespace.id, name: batchChange.name }), [
        batchChange.namespace.id,
        batchChange.name,
    ])

    const { data, error, loading, refetch } = useQuery<GetBatchChangeToEditResult, GetBatchChangeToEditVariables>(
        GET_BATCH_CHANGE_TO_EDIT,
        {
            variables,
            // Cache this data but always re-request it in the background when we revisit
            // this page to pick up newer changes.
            fetchPolicy: 'cache-and-network',
        }
    )

    const refetchBatchChange = useCallback(() => refetch(variables), [refetch, variables])

    if (loading && !data) {
        return (
            <div className="w-100 text-center">
                <Icon className="m-2" as={LoadingSpinner} />
            </div>
        )
    }

    if (!data?.batchChange || error) {
        return <HeroPage icon={AlertCircleIcon} title="Batch change not found" />
    }

    return (
        <BatchSpecContextProvider batchChange={data.batchChange} refetchBatchChange={refetchBatchChange}>
            <EditBatchSpecPageContent {...props} />
        </BatchSpecContextProvider>
    )
}

interface EditBatchSpecPageContentProps extends SettingsCascadeProps<Settings>, ThemeProps {}

const EditBatchSpecPageContent: React.FunctionComponent<React.PropsWithChildren<EditBatchSpecPageContentProps>> = ({
    settingsCascade,
    isLightTheme,
}) => {
    const { batchChange, editor, errors } = useContext(BatchSpecContext)

    const { insightTitle } = useInsightTemplates(settingsCascade)

    const [activeTabName, setActiveTabName] = useState<TabName>('batch spec')
    const tabsConfig = useMemo<TabsConfig[]>(
        () => [
            {
                name: 'configuration',
                isEnabled: true,
                handler: {
                    type: 'button',
                    onClick: () => setActiveTabName('configuration'),
                },
            },
            {
                name: 'batch spec',
                isEnabled: true,
                handler: {
                    type: 'button',
                    onClick: () => setActiveTabName('batch spec'),
                },
            },
        ],
        []
    )

    // TODO: Move to context??
    // NOTE: Technically there's only one option, and it's actually a preview option.
    const [executionOptions, setExecutionOptions] = useState<ExecutionOptions>({ runWithoutCache: false })

    return (
        <div className={layoutStyles.pageContainer}>
            {insightTitle && <InsightTemplatesBanner insightTitle={insightTitle} type="create" />}
            <div className={layoutStyles.headerContainer}>
                <BatchChangeHeader
                    namespace={{
                        to: `${batchChange.namespace.url}/batch-changes`,
                        text: batchChange.namespace.namespaceName,
                    }}
                    title={{ to: batchChange.url, text: batchChange.name }}
                    description={batchChange.description ?? undefined}
                />
                <ActionButtons>
                    <ExecutionOptionsDropdown
                        // execute={executeBatchSpec}
                        execute={() => alert('execute!')}
                        // isExecutionDisabled={isExecutionDisabled}
                        isExecutionDisabled={true}
                        // executionTooltip={executionTooltip}
                        executionTooltip="lol"
                        options={executionOptions}
                        onChangeOptions={setExecutionOptions}
                    />

                    {/* TODO: Come back to this after Adeola's PR is merged */}
                    {/* {downloadSpecModalDismissed ? (
                <BatchSpecDownloadLink name={batchChange.name} originalInput={code} isLightTheme={isLightTheme}>
                    or download for src-cli
                </BatchSpecDownloadLink>
            ) : (
                <Button className={styles.downloadLink} variant="link" onClick={() => setIsDownloadSpecModalOpen(true)}>
                    or download for src-cli
                </Button>
            )} */}
                </ActionButtons>
            </div>
            <TabBar activeTabName={activeTabName} tabsConfig={tabsConfig} />

            {activeTabName === 'configuration' ? (
                <ConfigurationForm isReadOnly={true} batchChange={batchChange} settingsCascade={settingsCascade} />
            ) : (
                <div className={styles.form}>
                    <LibraryPane name={batchChange.name} onReplaceItem={editor.handleCodeChange} />
                    <div className={styles.editorContainer}>
                        <h4 className={styles.header}>Batch spec</h4>
                        <MonacoBatchSpecEditor
                            batchChangeName={batchChange.name}
                            className={styles.editor}
                            isLightTheme={isLightTheme}
                            value={editor.code}
                            onChange={editor.handleCodeChange}
                        />
                        <EditorFeedbackPanel errors={errors} />
                    </div>
                    <WorkspacesPreviewPanel />
                </div>
            )}
        </div>
    )
}
