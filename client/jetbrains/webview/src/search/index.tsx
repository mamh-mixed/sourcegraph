import { render } from 'react-dom'

import { ContentMatch } from '@sourcegraph/shared/src/search/stream'
import { AnchorLink, setLinkComponent } from '@sourcegraph/wildcard'

import { App } from './App'
import { callJava } from './mockJavaInterface'

setLinkComponent(AnchorLink)

let isDarkTheme = false
let instanceURL = 'https://sourcegraph.com'

type RequestToJavaAction = 'getConfig' | 'getTheme' | 'saveLastSearch' | 'loadLastSearch'

export interface RequestToJava {
    action: RequestToJavaAction
    arguments: object
}

export interface Theme {
    isDarkTheme: boolean
    buttonColor: string
}

export interface PluginConfig {
    instanceURL: string
}

/* Add global functions to global window object */
declare global {
    interface Window {
        initializeSourcegraph: () => void
        callJava: (request: RequestToJava) => Promise<object>
    }
}

async function onOpen(match: ContentMatch, lineIndex: number): Promise<void> {
    console.log('open', await loadContent(match), match.lineMatches[lineIndex])
}

async function onPreviewChange(match: ContentMatch, lineIndex: number): Promise<void> {
    console.log('preview', await loadContent(match), match.lineMatches[lineIndex])
}

function onPreviewClear(): void {
    console.log('clear preview')
}

function renderReactApp(): void {
    const node = document.querySelector('#main') as HTMLDivElement
    render(
        <App
            isDarkTheme={isDarkTheme}
            instanceURL={instanceURL}
            onOpen={onOpen}
            onPreviewChange={onPreviewChange}
            onPreviewClear={onPreviewClear}
        />,
        node
    )
}

async function getConfig(): Promise<PluginConfig> {
    try {
        return (await window.callJava({ action: 'getConfig', arguments: {} })) as PluginConfig
    } catch (error) {
        console.error(`Failed to get config: ${(error as Error).message}`)
        return { instanceURL: 'https://sourcegraph.com' }
    }
}

function applyConfig(config: PluginConfig): void {
    instanceURL = config.instanceURL
}

async function getTheme(): Promise<Theme> {
    try {
        return (await window.callJava({ action: 'getTheme', arguments: {} })) as Theme
    } catch (error) {
        console.error(`Failed to get theme: ${(error as Error).message}`)
        return {
            isDarkTheme: true,
            buttonColor: '#0078d4',
        }
    }
}

function applyTheme(theme: Theme): void {
    // Dark/light theme
    document.documentElement.classList.add('theme')
    document.documentElement.classList.remove(theme.isDarkTheme ? 'theme-light' : 'theme-dark')
    document.documentElement.classList.add(theme.isDarkTheme ? 'theme-dark' : 'theme-light')
    isDarkTheme = theme.isDarkTheme

    // Button color (test)
    const buttonColor = theme.buttonColor
    const root = document.querySelector(':root') as HTMLElement
    if (buttonColor) {
        root.style.setProperty('--button-color', buttonColor)
    }
    root.style.setProperty('--primary', buttonColor)
}

window.initializeSourcegraph = async () => {
    const [theme, config] = await Promise.all([getTheme(), getConfig()])
    applyConfig(config)
    applyTheme(theme)
    renderReactApp()
}

/* Initialize app for standalone server */
if (window.location.search.includes('standalone=true')) {
    window.callJava = callJava
    window.initializeSourcegraph()
}
