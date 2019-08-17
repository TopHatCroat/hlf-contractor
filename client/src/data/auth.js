import { AUTH_LOGIN, AUTH_LOGOUT, AUTH_CHECK, showNotification } from 'react-admin';

export default (apiUrl, fabricCli) => {
    const login = (email, password) => {
        const request = new Request(`${apiUrl}/login`, {
            method: 'POST',
            body: JSON.stringify({ email, password }),
            headers: new Headers({ 'Content-Type': 'application/json' }),
        });

        return fetch(request)
            .then(response => {
                if (response.status < 200 || response.status >= 300) {
                    throw new Error('app.auth.login_error');
                }
                return response.json()
            }).then((json) => {
                localStorage.setItem('token', json.token);
            })
    };

    const register = (email, password) => {
        const request = new Request(`${apiUrl}/register`, {
            method: 'POST',
            body: JSON.stringify({ email, password }),
            headers: new Headers({ 'Content-Type': 'application/json' }),
        });

        return fetch(request)
            .then(response => {
                if (response.status < 200 || response.status >= 300) {
                    throw new Error('app.auth.register_error');
                }

                return login(email, password);
            })
    };

    return (type, params) => {
        if (type === AUTH_LOGIN) {
            const { kind, username, password } = params;
            if (kind === "register") {
                return register(username, password)
                    .catch((error) => {
                        showNotification(error.message, 'error');
                    })
            }


            login(username, password)
                .catch((error) => {
                    showNotification(error.message, 'error');
                });
        } else if (type === AUTH_LOGOUT) {
            localStorage.removeItem('token');
            return Promise.resolve();
        } else if (type === AUTH_CHECK) {
            return localStorage.getItem('token') ? Promise.resolve() : Promise.reject();
        }
        
        return Promise.resolve();
    }
}