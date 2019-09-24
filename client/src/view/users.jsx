import React from 'react';
import { List, Datagrid, EmailField, TextField } from 'react-admin';
import MonetaryField from "../components/MonetaryField";

export const UserList = props => (
    <List {...props}>
        <Datagrid rowClick="edit">
            <EmailField source="email" />
            <TextField source="status" />
            <MonetaryField source="balance" />
        </Datagrid>
    </List>
);