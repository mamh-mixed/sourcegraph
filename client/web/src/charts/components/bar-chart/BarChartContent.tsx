import { ReactElement, SVGProps, useRef, useState } from 'react'

import { Group } from '@visx/group'
import classNames from 'classnames'
import { ScaleBand, ScaleLinear } from 'd3-scale'

import { Tooltip } from '../../core'

import { GroupedBars } from './components/GroupedBars'
import { StackedBars } from './components/StackedBars'
import { BarTooltipContent } from './components/TooltipContent'
import { Category } from './utils/get-grouped-categories'

import styles from './BarChartContent.module.scss'

interface ActiveSegment<Datum> {
    category: Category<Datum>
    datum: Datum
}

interface BarChartContentProps<Datum> extends SVGProps<SVGGElement> {
    stacked: boolean

    top: number
    left: number

    xScale: ScaleBand<string>
    yScale: ScaleLinear<number, number>
    categories: Category<Datum>[]

    getDatumName: (datum: Datum) => string
    getDatumValue: (datum: Datum) => number
    getDatumColor: (datum: Datum) => string | undefined
    getDatumLink: (datum: Datum) => string | undefined | null
    onBarClick: (datum: Datum) => void
}

export function BarChartContent<Datum>(props: BarChartContentProps<Datum>): ReactElement {
    const {
        xScale,
        yScale,
        categories,
        stacked,
        top,
        left,
        width = 0,
        height = 0,
        getDatumName,
        getDatumValue,
        getDatumColor,
        getDatumLink,
        onBarClick,
        ...attributes
    } = props

    const rootRef = useRef<SVGGElement>(null)
    const [activeSegment, setActiveSegment] = useState<ActiveSegment<Datum> | null>(null)

    const withActiveLink = activeSegment?.datum ? getDatumLink(activeSegment?.datum) : null

    return (
        <Group
            {...attributes}
            innerRef={rootRef}
            className={classNames(styles.root, { [styles.rootWithHoveredLinkPoint]: withActiveLink })}
        >
            {stacked ? (
                <StackedBars
                    categories={categories}
                    xScale={xScale}
                    yScale={yScale}
                    getDatumName={getDatumName}
                    getDatumValue={getDatumValue}
                    getDatumColor={getDatumColor}
                    height={+height}
                    onBarHover={(datum, category) => setActiveSegment({ datum, category })}
                    onBarLeave={() => setActiveSegment(null)}
                    onBarClick={onBarClick}
                />
            ) : (
                <GroupedBars
                    categories={categories}
                    xScale={xScale}
                    yScale={yScale}
                    getDatumName={getDatumName}
                    getDatumValue={getDatumValue}
                    getDatumColor={getDatumColor}
                    getDatumLink={getDatumLink}
                    height={+height}
                    width={+width}
                    onBarHover={(datum, category) => setActiveSegment({ datum, category })}
                    onBarLeave={() => setActiveSegment(null)}
                    onBarClick={onBarClick}
                />
            )}

            {activeSegment && rootRef.current && (
                <Tooltip containerElement={rootRef.current as unknown as HTMLElement}>
                    <BarTooltipContent
                        category={activeSegment.category}
                        activeBar={activeSegment.datum}
                        getDatumColor={getDatumColor}
                        getDatumValue={getDatumValue}
                        getDatumName={getDatumName}
                    />
                </Tooltip>
            )}
        </Group>
    )
}
