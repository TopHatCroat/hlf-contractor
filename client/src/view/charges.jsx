import React from 'react';
import { List, Datagrid, TextField, ReferenceField } from 'react-admin';
import MonetaryField from "../components/MonetaryField";
import UserField from "../components/UserField";
import TimeAndDateField from "../components/TimeAndDateField";


export const ChargeTransactionList = props => (
    <List {...props}>
        <Datagrid rowClick="edit">
            <ReferenceField source="user_email" reference="users">
                <UserField  source="id" />
            </ReferenceField>
            <TextField source="contractor" />
            <MonetaryField source="price" />
            <TextField source="state" />
            <TimeAndDateField source="start_date" />
            <TimeAndDateField source="stop_date" options={{ weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }} />
        </Datagrid>
    </List>
);

