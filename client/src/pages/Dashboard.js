import React, { Component } from 'react'
import http from '../httpClient'
import { Consumer } from '../components/Authenticator'
import "antd/dist/antd.css"
import { Row, Col } from 'antd'

const Box = props => (
    <div style={{border: "2px", "border-style": "solid", height: "2em", width: "5em", backgroundColor: "blue"}}>
        { props.children }
    </div>
)
export default class Dashboard extends Component {
    state = {
        profile: null,
        current: "nothing",
        pyramid: [
            Array(1).fill(null),
            Array(2).fill(null),
            Array(3).fill(null),
            Array(4).fill(null),
            Array(5).fill(null),
            Array(6).fill(null),
            Array(7).fill(null),
            Array(8).fill(null),
        ]
    }

    componentDidMount() {
        http.get('/accountID', {headers: {'Authorization': 'Bearer ' + localStorage.getItem('jwtToken')}})
            .then(({data}) => { console.log(data); this.setState({profile: data }) })
            .catch(err => console.error(err))
    }

    render() {
        const { profile, pyramid, current } = this.state
        return (
            <div>
                <Row>
                    <Col span={8}>
                        { profile && profile.displayName}
                    </Col>
                    <Col span={8} offset={8}>
                        <Consumer>{context => (<button onClick={context.logout}>Logout</button>)}</Consumer>
                    </Col>
                </Row>
                <Row>
                    <Col span={4}>
                        { current }
                        <hr />
                        {["James", "Chris", "Jane", "Mike"].map((n, idx) => (
                            <Row key={idx}><Col span={6} onClick={e => this.setState({current: n})}>{n}</Col></Row>
                        ))}
                    </Col>
                    <Col span={20}>
                        {pyramid.map((tier, idx) => (
                            <Row key={idx} type="flex" justify="center">
                                { tier.map((t, i) => (<Box key={i}>{t}</Box>)) }
                            </Row>
                        ))}
                    </Col>
                </Row>
            </div>
        )
    }
}