import React from 'react';
import { List, Datagrid, TextField, EmailField, UrlField } from 'react-admin';
import CustomUrl from './custom-url';

export const UserList = props => (
    <List {...props}>
        <Datagrid rowClick="edit">
            <TextField source="id" />
            <TextField source="name" />
            <EmailField source="email" />
            <TextField source="phone" />
            <CustomUrl source="website" />
            <TextField source="company.name" />
        </Datagrid>
    </List>
);