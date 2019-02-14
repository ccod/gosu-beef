import React, { Component } from 'react'
import http from '../httpClient'
import { Consumer } from '../components/Authenticator'
//import Pyramid from '../components/Pyramid'
import "antd/dist/antd.css"
import { Row, Col } from 'antd'

const exampleLadder = [
    { place: 1, name: "scrubbles" },
    { place: 2, name: "Lost Jinjo" },
    { place: 3, name: "Dan" },
    { place: 4, name: "Jamary" },
    { place: 5, name: "Kirby" },
    { place: 6, name: "Pastry" },
    { place: 7, name: "Guildlin" },
    { place: 8, name: "Thor" },
    { place: 9, name: "Mac" },
    { place: 10, name: "Silverknight" },
    { place: 11, name: "Hypno" },
    { place: 12, name: "Nonickname" },
    { place: 13, name: "Water" },
    { place: 14, name: "PicolaHitman" },
    { place: 15, name: "Bodichitla" },
    { place: 16, name: "SimplePickup" },
    { place: 17, name: "Noahfecks" },
    { place: 18, name: "Bum" },
    { place: 19, name: "Blitz" },
    { place: 20, name: "Katharsis" },
    { place: 21, name: "ArcticKinger" },
]

export default class Dashboard extends Component {
    state = {
        profile: null,
        current: "nothing",
        players: exampleLadder
    }

    componentDidMount() {
        http.get('/player', {headers: {'Authorization': 'Bearer ' + localStorage.getItem('jwtToken')}})
            .then(({data}) => { console.log(data); this.setState({profile: data }) })
            .catch(err => console.error(err))

        http.get('/players', {headers: {'Authorization': 'Bearer ' + localStorage.getItem('jwtToken')}})
            .then(({data}) => { console.log(data); this.setState({players: data }) })
            .catch(err => console.error(err))
    }

    render() {
        const { profile, current, players } = this.state
        return (
            <div>
                <Row>
                    <Col span={8}>
                        { profile && profile.displayName}
                        { profile && profile.profileId }
                    </Col>
                    <Col span={8} offset={8}>
                        <Consumer>{context => (<button onClick={context.logout}>Logout</button>)}</Consumer>
                    </Col>
                </Row>
                <Row>
                    <Col span={4}>
                        { current.displayName }
                        <hr />
                        {players.map((n, idx) => (
                            <Row key={idx}><Col span={6} onClick={e => this.setState({current: n})}>{n.displayName}</Col></Row>
                        ))}
                    </Col>
                    <Col span={20}>
                            {/* <Pyramid players={players} /> */}
                    </Col>
                </Row>
            </div>
        )
    }
}