import React from 'react'

import BrainIcon from 'mdi-react/BrainIcon'

import { Badge, Icon, Link, Menu, MenuButton, MenuDivider, MenuHeader, MenuList, Position } from '@sourcegraph/wildcard'

import { useCodeIntelStatus } from './useCodeIntelStatus'

import styles from './RepositoryMenu.module.scss'

export interface RepositoryMenuProps {
    repoName: string
    revision: string
    filePath: string
    actionType: 'nav' | 'dropdown'
}

export const RepositoryMenu: React.FunctionComponent<RepositoryMenuProps> = ({
    repoName,
    revision,
    filePath,
    actionType,
}) => {
    const result = useCodeIntelStatus({ variables: { repository: repoName, commit: revision, path: filePath } })

    return actionType === 'dropdown' ? (
        <>TODO</>
    ) : (
        <Menu className="btn-icon">
            <>
                <MenuButton className="text-decoration-none">
                    <Icon as={BrainIcon} />
                </MenuButton>

                <MenuList position={Position.bottomEnd} className={styles.dropdownMenu}>
                    <MenuHeader>Code intelligence</MenuHeader>

                    <MenuDivider />

                    <div className="px-2 py-1">
                        <h2>Design template</h2>

                        <div className="d-flex align-items-center">
                            <div className="p-2 text-uppercase">
                                {/* <Badge>Unavailable</Badge> */}
                                {/* <Badge variant="primary">Available</Badge> */}
                                <Badge variant="success">Enabled</Badge>
                            </div>
                            <div className="p-2">
                                <span>Precise code intelligence</span>
                                <br />

                                <span className="text-muted">Last updated: 02/02/2022</span>
                                <br />

                                <span className="text-muted">80% Java supported</span>
                                <br />

                                <span className="text-muted">Index is failing</span>
                                <br />

                                <Link to="/">I want precise support!</Link>
                                <br />

                                <Link to="/">Enable precise code intelligence</Link>
                            </div>
                        </div>
                    </div>

                    <MenuDivider />

                    <div className="px-2 py-1">
                        <div className="d-flex align-items-center">
                            <div className="p-2">
                                <h2>Support</h2>

                                <ul>
                                    {result.data?.preciseSupport.map((preciseSupport, index) => (
                                        <li key={`precise-support-${index}`}>
                                            <span>
                                                <strong>Precise intelligence </strong> is available at level{' '}
                                                {preciseSupport.supportLevel} and confidence {preciseSupport.confidence}{' '}
                                                via{' '}
                                                {preciseSupport.indexers?.map((indexer, index) => (
                                                    <>
                                                        {index !== 0 ? ', ' : ''}
                                                        <span key={indexer.name}>
                                                            <Link to={indexer.url}>{indexer.name}</Link>
                                                        </span>
                                                    </>
                                                ))}
                                            </span>
                                            <br />
                                        </li>
                                    ))}

                                    {result.data?.searchBasedSupport.map((searchSupport, index) => (
                                        <li key={`search-support-${index}`}>
                                            <span>
                                                <strong>Search-based intelligence</strong> for language{' '}
                                                {searchSupport.language} is available at level{' '}
                                                {searchSupport.supportLevel}.
                                            </span>
                                            <br />
                                        </li>
                                    ))}
                                </ul>
                            </div>
                        </div>
                    </div>

                    <MenuDivider />

                    <div className="px-2 py-1">
                        <div className="d-flex align-items-center">
                            <div className="p-2">
                                <h2>Uploads</h2>

                                {result.data?.uploadIds.join(', ') || 'no lsif indexes covering this tree/blob'}
                            </div>
                        </div>
                    </div>

                    <MenuDivider />

                    <div className="px-2 py-1">
                        <div className="d-flex align-items-center">
                            <div className="p-2">
                                <h2>Indexes</h2>
                                BlKSDFU
                            </div>
                        </div>
                    </div>
                </MenuList>
            </>
        </Menu>
    )
}
