import React from 'react'

import classNames from 'classnames'
import BrainIcon from 'mdi-react/BrainIcon'

import { Icon, Menu, MenuButton, MenuDivider, MenuHeader, MenuList } from '@sourcegraph/wildcard'

export const RepositoryMenu: React.FunctionComponent<{}> = () => (
    <Menu className="btn-icon">
        <>
            <MenuButton className={classNames('text-decoration-none test-user-nav-item-toggle')}>
                <Icon as={BrainIcon} />
            </MenuButton>
            <MenuList>
                <MenuHeader>Code intelligence</MenuHeader>
                <MenuDivider />
                <div className="px-2 py-1">
                    <div className="d-flex align-items-center">Some stuff here?</div>
                </div>
            </MenuList>
        </>
    </Menu>
)
