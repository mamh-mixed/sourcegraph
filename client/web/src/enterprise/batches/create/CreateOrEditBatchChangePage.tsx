import React, { useCallback, useEffect, useMemo, useState } from 'react'

import { ApolloQueryResult } from '@apollo/client'
import classNames from 'classnames'
import { compact, noop } from 'lodash'
import AlertCircleIcon from 'mdi-react/AlertCircleIcon'
import InfoCircleOutlineIcon from 'mdi-react/InfoCircleOutlineIcon'
import LockIcon from 'mdi-react/LockIcon'
import { useHistory, useLocation } from 'react-router'

import { ErrorAlert } from '@sourcegraph/branded/src/components/alerts'
import { Form } from '@sourcegraph/branded/src/components/Form'
import { useMutation, useQuery } from '@sourcegraph/http-client'
import { Settings } from '@sourcegraph/shared/src/schema/settings.schema'
import {
    SettingsCascadeProps,
    SettingsOrgSubject,
    SettingsUserSubject,
} from '@sourcegraph/shared/src/settings/settings'
import { ThemeProps } from '@sourcegraph/shared/src/theme'
import { HeroPage } from '@sourcegraph/web/src/components/HeroPage'
import {
    PageHeader,
    Button,
    Container,
    Input,
    LoadingSpinner,
    FeedbackBadge,
    RadioButton,
    Icon,
    Panel,
} from '@sourcegraph/wildcard'

import { BatchChangesIcon } from '../../../batches/icons'
import { PageTitle } from '../../../components/PageTitle'
import {
    BatchChangeFields,
    EditBatchChangeFields,
    GetBatchChangeToEditResult,
    GetBatchChangeToEditVariables,
    CreateEmptyBatchChangeVariables,
    CreateEmptyBatchChangeResult,
    Scalars,
    BatchSpecWorkspaceResolutionState,
    CreateBatchSpecFromRawVariables,
    CreateBatchSpecFromRawResult,
} from '../../../graphql-operations'
import { BatchSpecDownloadLink } from '../BatchSpec'

import { GET_BATCH_CHANGE_TO_EDIT, CREATE_EMPTY_BATCH_CHANGE, CREATE_BATCH_SPEC_FROM_RAW } from './backend'
import { EditorFeedbackPanel } from './editor/EditorFeedbackPanel'
import { MonacoBatchSpecEditor } from './editor/MonacoBatchSpecEditor'
import { ExecutionOptions, ExecutionOptionsDropdown } from './ExecutionOptions'
import { LibraryPane } from './library/LibraryPane'
import { NamespaceSelector } from './NamespaceSelector'
import { useBatchSpecCode } from './useBatchSpecCode'
import { useExecuteBatchSpec } from './useExecuteBatchSpec'
import { useInitialBatchSpec } from './useInitialBatchSpec'
import { useNamespaces } from './useNamespaces'
import { useWorkspacesPreview } from './useWorkspacesPreview'
import { useImportingChangesets } from './workspaces-preview/useImportingChangesets'
import { useWorkspaces, WorkspacePreviewFilters } from './workspaces-preview/useWorkspaces'
import { WorkspacesPreview } from './workspaces-preview/WorkspacesPreview'

import styles from './CreateOrEditBatchChangePage.module.scss'

export interface CreateOrEditBatchChangePageProps extends ThemeProps, SettingsCascadeProps<Settings> {
    /**
     * The id for the namespace that the batch change should be created in, or that it
     * already belongs to, if it already exists.
     */
    initialNamespaceID?: Scalars['ID']
    /** The batch change name, if it already exists. */
    batchChangeName?: BatchChangeFields['name']
}

/**
 * CreateOrEditBatchChangePage is the new SSBC-oriented page for creating a new batch change
 * or editing and re-executing a new batch spec for an existing one.
 */
export const CreateOrEditBatchChangePage: React.FunctionComponent<CreateOrEditBatchChangePageProps> = ({
    initialNamespaceID,
    batchChangeName,
    ...props
}) => {
    const { data, error, loading, refetch } = useQuery<GetBatchChangeToEditResult, GetBatchChangeToEditVariables>(
        GET_BATCH_CHANGE_TO_EDIT,
        {
            // If we don't have the batch change name or namespace, the user hasn't created a
            // batch change yet, so skip the request.
            skip: !initialNamespaceID || !batchChangeName,
            variables: {
                namespace: initialNamespaceID as Scalars['ID'],
                name: batchChangeName as BatchChangeFields['name'],
            },
            // Cache this data but always re-request it in the background when we revisit
            // this page to pick up newer changes.
            fetchPolicy: 'cache-and-network',
        }
    )

    const refetchBatchChange = useCallback(
        () =>
            refetch({
                namespace: initialNamespaceID as Scalars['ID'],
                name: batchChangeName as BatchChangeFields['name'],
            }),
        [initialNamespaceID, batchChangeName, refetch]
    )

    if (!batchChangeName) {
        return <CreatePage namespaceID={initialNamespaceID} {...props} />
    }

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

    return <EditPage batchChange={data.batchChange} refetchBatchChange={refetchBatchChange} {...props} />
}

interface CreatePageProps extends SettingsCascadeProps<Settings> {
    /**
     * The namespace the batch change should be created in. If none is provided, it will
     * default to the user's own namespace.
     */
    namespaceID?: Scalars['ID']
}

const CreatePage: React.FunctionComponent<CreatePageProps> = ({ namespaceID, settingsCascade }) => {
    const [template, setTemplate] = useState<string | undefined>()
    const location = useLocation()

    const [createEmptyBatchChange, { loading: batchChangeLoading, error: batchChangeError }] = useMutation<
        CreateEmptyBatchChangeResult,
        CreateEmptyBatchChangeVariables
    >(CREATE_EMPTY_BATCH_CHANGE)
    const [createBatchSpecFromRaw, { loading: batchSpecLoading, error: batchSpecError }] = useMutation<
        CreateBatchSpecFromRawResult,
        CreateBatchSpecFromRawVariables
    >(CREATE_BATCH_SPEC_FROM_RAW)

    const loading = batchChangeLoading || batchSpecLoading
    const error = batchChangeError || batchSpecError

    const { namespaces, defaultSelectedNamespace } = useNamespaces(settingsCascade, namespaceID)

    // The namespace selected for creating the new batch change under.
    const [selectedNamespace, setSelectedNamespace] = useState<SettingsUserSubject | SettingsOrgSubject>(
        defaultSelectedNamespace
    )

    const [nameInput, setNameInput] = useState('')
    const [isNameValid, setIsNameValid] = useState<boolean>()

    useEffect(() => {
        const parameters = new URLSearchParams(location.search)
        if (parameters.has('kind')) {
            switch (parameters.get('kind')) {
                case 'replaceSymbol':
                    setTemplate(renameSymbolTemplate(nameInput))
                case 'goCheckerS1003':
                    setTemplate(goCheckerS1003Template(nameInput))
            }
        }
    }, [location.search, nameInput])

    const onNameChange = useCallback<React.ChangeEventHandler<HTMLInputElement>>(event => {
        setNameInput(event.target.value)
        setIsNameValid(NAME_PATTERN.test(event.target.value))
    }, [])

    const history = useHistory()
    const handleCancel = (): void => history.goBack()
    const handleCreate = (): void => {
        createEmptyBatchChange({
            variables: { namespace: selectedNamespace.id, name: nameInput },
        })
            .then(args =>
                args.data?.createEmptyBatchChange.id && template
                    ? createBatchSpecFromRaw({
                          variables: { namespace: selectedNamespace.id, spec: template, noCache: false },
                      }).then(() => Promise.resolve(args))
                    : Promise.resolve(args)
            )
            .then(({ data }) => (data ? history.push(`${data.createEmptyBatchChange.url}/edit`) : noop()))
            // We destructure and surface the error from `useMutation` instead.
            .catch(noop)
    }

    return (
        <div className="container">
            <div className="container col-8 my-4">
                <PageTitle title="Create new batch change" />
                <PageHeader
                    path={[{ icon: BatchChangesIcon, to: '.' }, { text: 'Create batch change' }]}
                    className="flex-1 pb-2"
                    description="Run custom code over hundreds of repositories and manage the resulting changesets."
                    annotation={
                        <FeedbackBadge status="experimental" feedback={{ mailto: 'support@sourcegraph.com' }} />
                    }
                />
                <Form className="my-4 pb-5" onSubmit={handleCreate}>
                    <Container className="mb-4">
                        {error && <ErrorAlert error={error} />}
                        <NamespaceSelector
                            namespaces={namespaces}
                            selectedNamespace={selectedNamespace.id}
                            onSelect={setSelectedNamespace}
                        />
                        <Input
                            label="Batch change name"
                            value={nameInput}
                            onChange={onNameChange}
                            // pattern={String(NAME_PATTERN)}
                            required={true}
                            status={isNameValid === undefined ? undefined : isNameValid ? 'valid' : 'error'}
                        />
                        <small className="text-muted">
                            Give it a short, descriptive name to reference the batch change on Sourcegraph. Do not
                            include confidential information.{' '}
                            <span className={classNames(isNameValid === false && 'text-danger')}>
                                Only regular characters, _ and - are allowed.
                            </span>
                        </small>
                        <hr className="my-3" />
                        <h3 className="text-muted">
                            Visibility <Icon data-tooltip="Coming soon" as={InfoCircleOutlineIcon} />
                        </h3>
                        <div className="form-group mb-1">
                            <RadioButton
                                name="visibility"
                                value="public"
                                className="mr-2"
                                checked={true}
                                disabled={true}
                                label="Public"
                                aria-label="Public"
                            />
                        </div>
                        <div className="form-group mb-0">
                            <RadioButton
                                name="visibility"
                                value="private"
                                className="mr-2 mb-0"
                                disabled={true}
                                label={
                                    <>
                                        Private <Icon className="text-warning" aria-hidden={true} as={LockIcon} />
                                    </>
                                }
                                aria-label="Private"
                            />
                        </div>
                    </Container>
                    <div>
                        <Button
                            variant="primary"
                            type="submit"
                            onClick={handleCreate}
                            disabled={loading || nameInput === '' || !isNameValid}
                            className="mr-2"
                        >
                            Create batch change
                        </Button>
                        <Button variant="secondary" type="button" outline={true} onClick={handleCancel}>
                            Cancel
                        </Button>
                    </div>
                </Form>
            </div>
        </div>
    )
}

function goCheckerS1003Template(name: string): string{
    return `name: ${name}
    description: |
      This batch change uses [Comby](https://comby.dev) to replace calls to strings.Index with strings.Contains.
    on:
      - repositoriesMatchingQuery: strings.Index(:[1], :[2]) < 0  or strings.Index(:[1], :[2]) == -1 or strings.Index(:[1], :[2]) != -1 or strings.Index(:[1], :[2]) >= 0 or strings.Index(:[1], :[2]) > -1 or strings.IndexAny(:[1], :[2]) < 0 or strings.IndexAny(:[1], :[2]) == -1 or strings.IndexAny(:[1], :[2]) != -1 or strings.IndexAny(:[1], :[2]) >= 0 or strings.IndexAny(:[1], :[2]) > -1 or strings.IndexRune(:[1], :[2]) < 0 or strings.IndexRune(:[1], :[2]) == -1 or strings.IndexRune(:[1], :[2]) != -1 or strings.IndexRune(:[1], :[2]) >= 0 or strings.IndexRune(:[1], :[2]) > -1 patternType:structural archived:no count:1

    steps:
      - run: comby -config /tmp/rule.toml -f .go -i -exclude-dir vendor,.
        container: comby/comby
        files:
          /tmp/rule.toml: |
            [S1003_01]
            match='strings.IndexRune(:[1], :[2]) > -1'
            rewrite='strings.ContainsRune(:[1], :[2])'

            [S1003_02]
            match='strings.IndexRune(:[1], :[2]) >= 0'
            rewrite='strings.ContainsRune(:[1], :[2])'

            [S1003_03]
            match='strings.IndexRune(:[1], :[2]) != -1'
            rewrite='strings.ContainsRune(:[1], :[2])'

            [S1003_04]
            match='strings.IndexRune(:[1], :[2]) == -1'
            rewrite='!strings.ContainsRune(:[1], :[2])'

            [S1003_05]
            match='strings.IndexRune(:[1], :[2]) < 0'
            rewrite='!strings.ContainsRune(:[1], :[2])'

            [S1003_06]
            match='strings.IndexAny(:[1], :[2]) > -1'
            rewrite='strings.ContainsAny(:[1], :[2])'

            [S1003_07]
            match='strings.IndexAny(:[1], :[2]) >= 0'
            rewrite='strings.ContainsAny(:[1], :[2])'

            [S1003_08]
            match='strings.IndexAny(:[1], :[2]) != -1'
            rewrite='strings.ContainsAny(:[1], :[2])'

            [S1003_09]
            match='strings.IndexAny(:[1], :[2]) == -1'
            rewrite='!strings.ContainsAny(:[1], :[2])'

            [S1003_10]
            match='strings.IndexAny(:[1], :[2]) < 0'
            rewrite='!strings.ContainsAny(:[1], :[2])'

            [S1003_11]
            match='strings.Index(:[1], :[2]) > -1'
            rewrite='strings.Contains(:[1], :[2])'

            [S1003_12]
            match='strings.Index(:[1], :[2]) >= 0'
            rewrite='strings.Contains(:[1], :[2])'

            [S1003_13]
            match='strings.Index(:[1], :[2]) != -1'
            rewrite='strings.Contains(:[1], :[2])'

            [S1003_14]
            match='strings.Index(:[1], :[2]) == -1'
            rewrite='!strings.Contains(:[1], :[2])'

            [S1003_15]
            match='strings.Index(:[1], :[2]) < 0'
            rewrite='!strings.Contains(:[1], :[2])'
      - run: |
          results="\${{ join repository.search_result_paths "\n" }}"
          for result in $results;
          do
            echo "Rewriting \${result}"
            gofmt -w "\${result}"
          done;
        container: golang

    changesetTemplate:
      title: Replace calls to strings.Index with strings.Contains.
      body: This batch change uses [Comby](https://comby.dev) to replace calls to strings.Index with strings.Contains.
      branch: batches/\${{batch_change.name}}
      commit:
        message: Replace calls to strings.Index with strings.Contains.

`
}

function renameSymbolTemplate(name: string): string {
    return `name: ${name}
description: Renames the symbol "foo" to "bar"

on:
    - repository: github.com/hashicorp/errwrap

steps:
    - run: renamer -repoPath=/work -repoName=github.com/hashicorp/errwrap -rev=7b00e5db719c64d14dd0caaacbd13e76254d02c0 -filePath=errwrap.go -line=162 -character=6 -replacement=BatchYeah
      container: eseliger/renamer:0.0.1
      env:
        SRC_ENDPOINT: https://k8s.sgdev.org
        SRC_ACCESS_TOKEN: 494bf58148bc74c0370c685e3bd033b6152653f4

changesetTemplate:
    title: Hello World
    body: My first batch change!
    branch: hello-world
    commit:
        message: Append Hello World to all README.md files
`
}

const INVALID_BATCH_SPEC_TOOLTIP = "There's a problem with your batch spec."
const WORKSPACES_PREVIEW_SIZE = 'batch-changes.ssbc-workspaces-preview-size'

interface EditPageProps extends ThemeProps {
    batchChange: EditBatchChangeFields
    refetchBatchChange: () => Promise<ApolloQueryResult<GetBatchChangeToEditResult>>
}

const EditPage: React.FunctionComponent<EditPageProps> = ({ batchChange, refetchBatchChange, isLightTheme }) => {
    // Get the latest batch spec for the batch change.
    const { batchSpec, isApplied: isLatestBatchSpecApplied, initialCode: initialBatchSpecCode } = useInitialBatchSpec(
        batchChange
    )

    // Manage the batch spec input YAML code that's being edited.
    const { code, debouncedCode, isValid, handleCodeChange, excludeRepo, errors: codeErrors } = useBatchSpecCode(
        initialBatchSpecCode,
        batchChange.name
    )

    const [filters, setFilters] = useState<WorkspacePreviewFilters>()
    const workspacesConnection = useWorkspaces(batchSpec.id, filters)
    const importingChangesetsConnection = useImportingChangesets(batchSpec.id)

    // When we successfully submit the latest batch spec code to the backend for a new
    // workspaces preview, we follow up by refetching the batch change to get the latest
    // batch spec ID.
    const onComplete = useCallback(() => {
        // We handle any error here higher up the chain, so we can ignore it.
        refetchBatchChange().then(noop).catch(noop)
    }, [refetchBatchChange])

    // NOTE: Technically there's only one option, and it's actually a preview option.
    const [executionOptions, setExecutionOptions] = useState<ExecutionOptions>({ runWithoutCache: false })

    // Manage the batch spec that was last submitted to the backend for the workspaces preview.
    const {
        preview: previewBatchSpec,
        isInProgress: isWorkspacesPreviewInProgress,
        error: previewError,
        clearError: clearPreviewError,
        hasPreviewed,
        cancel,
        resolutionState,
    } = useWorkspacesPreview(batchSpec.id, {
        isBatchSpecApplied: isLatestBatchSpecApplied,
        namespaceID: batchChange.namespace.id,
        noCache: executionOptions.runWithoutCache,
        onComplete,
        filters,
    })

    const clearErrorsAndHandleCodeChange = useCallback(
        (newCode: string) => {
            clearPreviewError()
            handleCodeChange(newCode)
        },
        [handleCodeChange, clearPreviewError]
    )

    // Disable the preview button if the batch spec code is invalid or the on: statement
    // is missing, or if we're already processing a preview.
    const previewDisabled = useMemo(
        () => (isValid !== true ? INVALID_BATCH_SPEC_TOOLTIP : isWorkspacesPreviewInProgress),
        [isValid, isWorkspacesPreviewInProgress]
    )

    // The batch spec YAML code is considered stale if any part of it changes. This is
    // because of a current limitation of the backend where we need to re-submit the batch
    // spec code and wait for the new workspaces preview to finish resolving before we can
    // execute, or else the execution will use an older batch spec. We will address this
    // when we implement the "auto-saving" feature and decouple previewing workspaces from
    // updating the batch spec code.
    const isBatchSpecStale = useMemo(() => initialBatchSpecCode !== debouncedCode, [
        initialBatchSpecCode,
        debouncedCode,
    ])

    // Manage submitting a batch spec for execution.
    const { executeBatchSpec, isLoading: isExecuting, error: executeError } = useExecuteBatchSpec(batchSpec.id)

    // Disable the execute button if any of the following are true:
    // - The batch spec code is invalid.
    // - There was an error with the preview.
    // - We're in the middle of previewing or executing the batch spec.
    // - We haven't yet submitted the batch spec to the backend yet for a preview.
    // - The batch spec on the backend is stale.
    // - The current workspaces evaluation is not complete.
    const [isExecutionDisabled, executionTooltip] = useMemo(() => {
        const isExecutionDisabled = Boolean(
            isValid !== true ||
                previewError ||
                isWorkspacesPreviewInProgress ||
                isExecuting ||
                !hasPreviewed ||
                isBatchSpecStale ||
                resolutionState !== BatchSpecWorkspaceResolutionState.COMPLETED
        )
        // The execution tooltip only shows if the execute button is disabled, and explains why.
        const executionTooltip =
            isValid === false || previewError
                ? INVALID_BATCH_SPEC_TOOLTIP
                : !hasPreviewed
                ? 'Preview workspaces first before you run.'
                : isBatchSpecStale
                ? 'Update your workspaces preview before you run.'
                : isWorkspacesPreviewInProgress || resolutionState !== BatchSpecWorkspaceResolutionState.COMPLETED
                ? 'Wait for the preview to finish first.'
                : undefined

        return [isExecutionDisabled, executionTooltip]
    }, [
        hasPreviewed,
        isValid,
        previewError,
        isWorkspacesPreviewInProgress,
        isExecuting,
        isBatchSpecStale,
        resolutionState,
    ])

    const actionButtons = (
        <>
            <ExecutionOptionsDropdown
                execute={executeBatchSpec}
                isExecutionDisabled={isExecutionDisabled}
                executionTooltip={executionTooltip}
                options={executionOptions}
                onChangeOptions={setExecutionOptions}
            />
            <BatchSpecDownloadLink name={batchChange.name} originalInput={code} isLightTheme={isLightTheme}>
                or download for src-cli
            </BatchSpecDownloadLink>
        </>
    )

    return (
        <BatchChangePage
            namespace={batchChange.namespace}
            title={batchChange.name}
            description={batchChange.description}
            actionButtons={actionButtons}
        >
            <div className={classNames(styles.editorLayoutContainer, 'd-flex flex-1 mt-2')}>
                <LibraryPane name={batchChange.name} onReplaceItem={clearErrorsAndHandleCodeChange} />
                <div className={styles.editorContainer}>
                    <h4 className={styles.header}>Batch spec</h4>
                    <MonacoBatchSpecEditor
                        batchChangeName={batchChange.name}
                        className={styles.editor}
                        isLightTheme={isLightTheme}
                        value={code}
                        onChange={clearErrorsAndHandleCodeChange}
                    />
                    <EditorFeedbackPanel
                        errors={compact([codeErrors.update, codeErrors.validation, previewError, executeError])}
                    />
                </div>
                <Panel
                    defaultSize={500}
                    minSize={405}
                    maxSize={1400}
                    position="right"
                    storageKey={WORKSPACES_PREVIEW_SIZE}
                >
                    <div className={styles.workspacesPreviewContainer}>
                        <WorkspacesPreview
                            previewDisabled={previewDisabled}
                            preview={() => previewBatchSpec(debouncedCode)}
                            batchSpecStale={
                                isBatchSpecStale || isWorkspacesPreviewInProgress || resolutionState === 'CANCELED'
                            }
                            hasPreviewed={hasPreviewed}
                            excludeRepo={excludeRepo}
                            cancel={cancel}
                            isWorkspacesPreviewInProgress={isWorkspacesPreviewInProgress}
                            resolutionState={resolutionState}
                            workspacesConnection={workspacesConnection}
                            importingChangesetsConnection={importingChangesetsConnection}
                            setFilters={setFilters}
                        />
                    </div>
                </Panel>
            </div>
        </BatchChangePage>
    )
}

const getNamespaceDisplayName = (namespace: SettingsUserSubject | SettingsOrgSubject): string => {
    switch (namespace.__typename) {
        case 'User':
            return namespace.displayName ?? namespace.username
        case 'Org':
            return namespace.displayName ?? namespace.name
    }
}

/** TODO: This duplicates the URL field from the org/user resolvers on the backend, but we
 * don't have access to that from the settings cascade presently. Can we get it included
 * in the cascade instead somehow? */
const getNamespaceBatchChangesURL = (namespace: SettingsUserSubject | SettingsOrgSubject): string => {
    switch (namespace.__typename) {
        case 'User':
            return '/users/' + namespace.username + '/batch-changes'
        case 'Org':
            return '/organizations/' + namespace.name + '/batch-changes'
    }
}

interface BatchChangePageProps {
    /** The namespace that should appear in the topmost `PageHeader`. */
    namespace: SettingsUserSubject | SettingsOrgSubject
    /** The title to use in the topmost `PageHeader`, alongside the `namespaceName`. */
    title: string
    /** The description to use in the topmost `PageHeader` beneath the titles. */
    description?: string | null
    /** Optionally, any action buttons that should appear in the top left of the page. */
    actionButtons?: JSX.Element
}

/**
 * BatchChangePage is a page layout component that renders a consistent header for
 * SSBC-style batch change pages and should wrap the other content contained on the page.
 */
const BatchChangePage: React.FunctionComponent<BatchChangePageProps> = ({
    children,
    namespace,
    title,
    description,
    actionButtons,
}) => (
    <div className="d-flex flex-column p-4 w-100 h-100">
        <div className="d-flex flex-0 justify-content-between align-items-start">
            <PageHeader
                path={[
                    { icon: BatchChangesIcon },
                    {
                        to: getNamespaceBatchChangesURL(namespace),
                        text: getNamespaceDisplayName(namespace),
                    },
                    { text: title },
                ]}
                className="flex-1 pb-2"
                description={
                    description || 'Run custom code over hundreds of repositories and manage the resulting changesets.'
                }
                annotation={<FeedbackBadge status="experimental" feedback={{ mailto: 'support@sourcegraph.com' }} />}
            />
            {actionButtons && (
                <div className="d-flex flex-column flex-0 align-items-center justify-content-center">
                    {actionButtons}
                </div>
            )}
        </div>
        {children}
    </div>
)
/* Regex pattern for a valid batch change name. Needs to match what's defined in the BatchSpec JSON schema. */
const NAME_PATTERN = /^[\w.-]+$/
