import React from 'react';
import { List, Datagrid, TextField, DateField } from 'react-admin';
import MonetaryField from "../components/MonetaryField";

export const ChargeTransactionList = props => (
    <List {...props}>
        <Datagrid rowClick="edit">
            <TextField source="status" />
            <MonetaryField source="price" />
            <DateField source="startTime" />
            <DateField source="endTime" />
        </Datagrid>
    </List>
);