import { AUTH_LOGIN, AUTH_LOGOUT, AUTH_CHECK } from 'react-admin';

export default (fabricCli) => {
    const login = (email, password) => {
        const request = new Request('http://api.awesome.agency:8000/login', {
            method: 'POST',
            body: JSON.stringify({ email, password }),
            headers: new Headers({ 'Content-Type': 'application/json' }),
        });

        return fetch(request)
            .then(response => {
                if (response.status < 200 || response.status >= 300) {
                    throw new Error(response.statusText);
                }
                return response.json()
            }).then((json) => {
                localStorage.setItem('token', json.token);
            })
    };

    const register = (email, password) => {
        const request = new Request('http://api.awesome.agency:8000/register', {
            method: 'POST',
            body: JSON.stringify({ email, password }),
            headers: new Headers({ 'Content-Type': 'application/json' }),
        });

        return fetch(request)
            .then(response => {
                if (response.status < 200 || response.status >= 300) {
                    throw new Error(response.statusText);
                }

                return login(email, password);
            })
    };

    return (type, params) => {
        if (type === AUTH_LOGIN) {
            const { kind, username, password } = params;
            if (kind === "register") {
                return register(username, password)
            }

            login(username, password);
        } else if (type === AUTH_LOGOUT) {
            localStorage.removeItem('token');
            return Promise.resolve();
        } else if (type === AUTH_CHECK) {
            const token = localStorage.getItem('token');
            if(token === null || token === undefined || token === "") {
                throw new Error("not_logged_in");
            }
        }
        
        return Promise.resolve();
    }
}