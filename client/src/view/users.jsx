import React from 'react';
import {
    Edit,
    List,
    Datagrid,
    SimpleForm,
    EmailField,
    TextField
} from 'react-admin';
import MonetaryField from "../components/MonetaryField";

export const UserList = props => (
    <List {...props}>
        <Datagrid rowClick="edit">
            <EmailField source="email" />
            <TextField source="state" />
            <MonetaryField source="balance" />
        </Datagrid>
    </List>
);

export const UserEdit = props => (
    <Edit hasEdit={false} {...props}>
        <SimpleForm>
            <TextField source="contractor" />
        </SimpleForm>
    </Edit>
);
