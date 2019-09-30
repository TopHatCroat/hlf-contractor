import React from 'react';
import { Admin, Resource } from 'react-admin';
import i18nProvider from './data/i18n';
import fabricDataProvider from './data/provider';
import fabricAuthProvider from './data/auth';
import FabricClient from './fabric/client';
import Login from './login/Login'
import { UserList } from './view/users';
import {
    ChargeTransactionCreate,
    ChargeTransactionList,
    ChargeTransactionEdit,
} from './view/charges';

const fabricConfig = {};
// const fabricConfig = fs.readFileSync('../config.yaml', 'utf8');
const apiUrl = "http://api.awesome.agency:8000";

const fabricCli = new FabricClient(fabricConfig);
const authProvider = fabricAuthProvider(apiUrl, fabricCli);
const dataProvider = fabricDataProvider(apiUrl);


const App = () => <Admin
        authProvider={authProvider}
        dataProvider={dataProvider}
        i18nProvider={i18nProvider}
        loginPage={Login}
        >
    <Resource
        name="charges"
        list={ChargeTransactionList}
        create={ChargeTransactionCreate}
        edit={ChargeTransactionEdit}
    />
    <Resource name="users" list={UserList} />
</Admin>;

export default App;
