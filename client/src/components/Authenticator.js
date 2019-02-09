import React, { Component } from 'react'
import history from '../history'

const {Provider, Consumer} = React.createContext()

class AuthProvider extends Component {
    logout = () => {
        localStorage.removeItem('jwtToken')
        history.replace('/')
    }

    setSession = () => {
        let token = window.location.hash.substr(1)

        if (token === "") {
            history.replace('/')
            return
        }

        localStorage.setItem('jwtToken', window.location.hash.substr(1))
        history.replace('/dashboard')
    }

    isAuthenticated = () => {
        return !!localStorage.getItem('jwtToken')
    }

    render() {
        return (
            <Provider value={this}>
                {this.props.children}
            </Provider>
        )
    }
}

export { AuthProvider, Consumer }