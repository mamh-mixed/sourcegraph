import React from 'react'

import BrainIcon from 'mdi-react/BrainIcon'

import { Badge, Icon, Link, Menu, MenuButton, MenuDivider, MenuHeader, MenuList, Position } from '@sourcegraph/wildcard'

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
}) =>
    actionType === 'dropdown' ? (
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
                        <div className="d-flex align-items-center">
                            <div className="p-2 text-uppercase">
                                {/* <Badge>Unavailable</Badge>
                            <Badge variant="primary">Available</Badge> */}
                                <Badge variant="success">Enabled</Badge>
                            </div>
                            <div className="p-2">
                                <span>Precise code intelligence</span>
                                <br />

                                <span>
                                    `{repoName}@{revision}:{filePath}`
                                </span>
                                <br />

                                <span className="text-muted">Last updated: 02/02/2022</span>
                                <br />

                                {/* <span className="text-muted">80% Java supported</span><br /> */}
                                {/* <span className="text-muted">Index is failing</span><br /> */}
                                {/* <Link to="sdf">I want precise support!</Link><br /> */}

                                <Link to="sdf">Enable precise code intelligence</Link>
                            </div>
                        </div>
                    </div>
                </MenuList>
            </>
        </Menu>
    )
