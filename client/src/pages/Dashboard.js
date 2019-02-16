import React, { Component } from 'react'
import http from '../httpClient'
import { Consumer } from '../components/Authenticator'
import Pyramid from '../components/Pyramid'
import "antd/dist/antd.css"
import { Row, Col, Tabs } from 'antd'

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

const TabPane = Tabs.TabPane

const Challenge = props => (
    <>
        <Col span={4}>
            { props.current.displayName }
            <hr />
            {props.players.map((n, idx) => (
                <Row key={idx}><Col span={6} onClick={props.select(n)}>{n.displayName}</Col></Row>
            ))}
        </Col>
        <Col span={20}>
                { <Pyramid mode="placement" players={props.players} rankings={props.rankings} click={props.click} /> }
        </Col>
    </>
)


// setup tabs with modes
// pyramid based on logged in users, pyramid based on rank
export default class Dashboard extends Component {
    state = {
        profile: null,
        current: { displayName: "nothing" },
        players: [],
        rankings: [],
        mode: null
    }

    selectCurrent = player => () => this.setState({current: player})

    setPlayerRank = rank => () => {
        if (this.state.current.displayName === "nothing") return

        http.post(
            '/rankings', 
            {rank, player: this.state.current}, 
            {headers: {'Authorization': 'Bearer ' + localStorage.getItem('jwtToken')}}
        ).then(({data}) => { 
            console.log("ranking post:\n", data) 
            if (this.state.rankings.find(r => r.rank === rank)) {
                this.setState({rankings: this.state.rankings.map(r => rank === r.rank ? data : r)})
            } else {
                this.setState({rankings: this.state.rankings.concat(data)})
            }
        }).catch(err => console.log(err))
    }


    componentDidMount() {
        http.get('/player', {headers: {'Authorization': 'Bearer ' + localStorage.getItem('jwtToken')}})
            .then(({data}) => { console.log("player data:\n",data); this.setState({profile: data }) })
            .catch(err => console.error(err))

        http.get('/players', {headers: {'Authorization': 'Bearer ' + localStorage.getItem('jwtToken')}})
            .then(({data}) => { console.log("players data:\n",data); this.setState({players: data }) })
            .catch(err => console.error(err))
        
        http.get("/rankings", {headers: {'Authorization': 'Bearer ' + localStorage.getItem('jwtToken')}})
            .then(({data}) => { console.log("rankings data:\n",data); this.setState({rankings: data }) })
            .catch(err => console.log(err))
    }

    render() {
        const { profile, rankings, current, players } = this.state
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
                    <Tabs defaultActiveKey="place">
                        <TabPane tab="challenge" key="challenge">Challenge</TabPane>
                        <TabPane tab="admin place" key="place">
                            <Challenge current={current} players={players} rankings={rankings} select={this.selectCurrent} click={this.setPlayerRank} />
                        </TabPane>
                        <TabPane tab="hello" key="hello">hello</TabPane>
                    </Tabs>
                </Row>
            </div>
        )
    }
}