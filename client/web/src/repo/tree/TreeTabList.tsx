import React, { useMemo } from 'react'

import classNames from 'classnames'
import AccountIcon from 'mdi-react/AccountIcon'
import FileDocumentIcon from 'mdi-react/FileDocumentIcon'
import HistoryIcon from 'mdi-react/HistoryIcon'
import SourceBranchIcon from 'mdi-react/SourceBranchIcon'
import SourceCommitIcon from 'mdi-react/SourceCommitIcon'
import TagIcon from 'mdi-react/TagIcon'
import { useCallbackRef } from 'use-callback-ref'

import { TreeFields } from '@sourcegraph/shared/src/graphql-operations'
import { Icon, Link } from '@sourcegraph/wildcard'

interface TreeTabList {
    tree: TreeFields
    selectedTab: string
    setSelectedTab: (tab: string) => void
}

export const TreeTabList: React.FunctionComponent<React.PropsWithChildren<TreeTabList>> = ({
    tree,
    selectedTab,
    setSelectedTab,
}) => {
    type Tabs = { tab: string; title: string; isActive: boolean; logName: string; icon: JSX.Element; url: string }[]

    const tabs: Tabs = useMemo(
        () => [
            {
                tab: 'home',
                title: 'Home',
                isActive: selectedTab === 'home',
                logName: 'RepoHomeTab',
                icon: <Icon as={FileDocumentIcon} aria-hidden={true} />,
                url: `${tree.url}/`,
            },
            {
                tab: 'commits',
                title: 'Commits',
                isActive: selectedTab === 'commits',
                logName: 'RepoCommitsTab',
                icon: <Icon as={SourceCommitIcon} aria-hidden={true} />,
                url: `${tree.url}/-/commits/tab`,
            },
            {
                tab: 'branch',
                title: 'Branches',
                isActive: selectedTab === 'branch',
                logName: 'RepoBranchesTab',
                icon: <Icon as={SourceBranchIcon} aria-hidden={true} />,
                url: `${tree.url}/-/branch/tab`,
            },
            {
                tab: 'tags',
                title: 'Tags',
                isActive: selectedTab === 'tags',
                logName: 'RepoTagsTab',
                icon: <Icon as={TagIcon} aria-hidden={true} />,
                url: `${tree.url}/-/tag/tab`,
            },
            {
                tab: 'compare',
                title: 'Compare',
                isActive: selectedTab === 'compare',
                logName: 'RepoCompareTab',
                icon: <Icon as={HistoryIcon} aria-hidden={true} />,
                url: `${tree.url}/-/compare/tab`,
            },
            {
                tab: 'contributors',
                title: 'Contributors',
                isActive: selectedTab === 'contributors',
                logName: 'RepoContributorsTab',
                icon: <Icon as={AccountIcon} aria-hidden={true} />,
                url: `${tree.url}/-/contributors/tab`,
            },
        ],
        [selectedTab, tree.url]
    )

    const callbackReference = useCallbackRef<HTMLAnchorElement>(null, ref => ref?.focus())

    return (
        <nav className="d-flex mb-4">
            <ul className="nav nav-tabs w-100">
                {tabs.map(({ tab, title, isActive, icon, url }) => (
                    <li className="nav-item" key={`repo-${tab}-tab`}>
                        <Link
                            to={url}
                            className={classNames('nav-link text-content bg-transparent', isActive && 'active')}
                            onClick={() => setSelectedTab(tab)}
                            ref={selectedTab === tab ? callbackReference : null}
                        >
                            <div>
                                {icon}
                                <span className="d-inline-flex ml-1" data-tab-content={title}>
                                    {title}
                                </span>
                            </div>
                        </Link>
                    </li>
                ))}
            </ul>
        </nav>
    )
}
