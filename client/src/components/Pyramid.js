import React from 'react'
import { Row } from 'antd'

const recurb = (former, current, check) => {
    if (former + current >= check) {
        let place  = check - former - 1
        return [current - 1, place]
    }

    return recurb(former + current, current + 1, check)
}

const findIndex = num => recurb(0, 1, num)

// use the last place player to find the upper bound on the necessary tiers
const makeGrid = num => {
    let tiers = findIndex(num)[0] + 1
    return Array(tiers).fill(null).map((_, idx) => Array(idx + 1).fill(null))
}

const insertPlayers = (grid, players) => {
    players.forEach(p => {
        let coord = findIndex(p.place)
        grid[coord[0]][coord[1]] = p
    })

    return grid
}

const Box = props => (
    <div {...props} style={{border: "2px", "borderStyle": "solid", height: "2em", width: "5em", backgroundColor: "blue"}}>
        { props.children }
    </div>
)

export default props => {
    let lastPlace = Math.max.apply(null, props.players.map(x => x.place)),
        grid = makeGrid(lastPlace),
        filledGrid = insertPlayers(grid, props.players)

    return (
        <>
            {filledGrid.map((tier, idx) => (
                <Row key={idx} type="flex" justify="center">
                    { tier.map((t, i) => (<Box key={i}>{t.name}</Box>)) }
                </Row>
            ))}

        </>

    )
}