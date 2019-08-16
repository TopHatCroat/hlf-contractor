import React from 'react';
import { Admin } from 'react-admin';

import Login from './login/Login'
import fabricDataProvider from './data/provider';
import fabricAuthProvider from './data/auth';
import FabricClient from './fabric/client';

const fabricConfig = {};
// const fabricConfig = fs.readFileSync('../config.yaml', 'utf8');
const fabricCli = new FabricClient(fabricConfig);

const authProvider = fabricAuthProvider(fabricCli);
const dataProvider = fabricDataProvider(fabricCli, fabricAuthProvider);

const App = () => <Admin loginPage={Login} authProvider={authProvider} dataProvider={dataProvider} />;

export default App;
