import React, { Component } from 'react'
import { Router, Route, Redirect } from 'react-router-dom'
import { AuthProvider, Consumer } from '../components/Authenticator'
import history from '../history'
import http from '../httpClient'

const Landing = () => (
    <div>
        <p>Landing page</p>
        <a href={process.env.REACT_APP_DOMAIN + '/login'}>Login</a>
    </div>
)

class Dashboard extends Component {
    componentDidMount() {
        http.get('/accountID', {headers: {'Authorization': 'Bearer ' + localStorage.getItem('jwtToken')}})
            .then(data => console.log(data))
            .catch(err => console.error(err))
    }

    render() {
        return (
            <div>
                <Consumer>{context => (<button onClick={context.logout}>Logout</button>)}</Consumer>
                <p>Dashboard</p>
                <p>{localStorage.getItem('jwtToken')}</p>
            </div>
        )
    }
}

const Callback = () => (
    <Consumer>
        {context => {
            context.setSession()
            return (<div>Loading</div>)
        }}
    </Consumer>
)

const PrivateRoute = ({ component: Component, ...rest }) => (
    <Consumer>
        {context => {
            return <Route
            {...rest}
            render={props => (
                context.isAuthenticated() ? (
                    <Component {...props} />
                ) : (
                    <Redirect to={{pathname: "/", state: { from: props.location }}} />
                )
            )} />
        }}
    </Consumer>
)

class Index extends Component {
    render() {
        return (
            <Router history={history}>
                <AuthProvider>
                    <Route exact path="/" component={Landing} />
                    <PrivateRoute path="/dashboard" component={Dashboard} />
                    <Route path="/callback" component={Callback} />
                </AuthProvider>
            </Router>
        )
    }
}

export default Index