import React from 'react';
import {
    Create,
    Edit,
    SimpleForm,
    List,
    Datagrid,
    TextField,
    ReferenceField,
    SelectInput,
} from 'react-admin';
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
            <TimeAndDateField source="stop_date" />
        </Datagrid>
    </List>
);

export const ChargeTransactionCreate = props => (
    <Create {...props}>
        <SimpleForm>
            <SelectInput source="contractor" choices={[
                { id: 'Pharmatic', name: 'Pharmatic' },
            ]} />
        </SimpleForm>
    </Create>
);

export const ChargeTransactionEdit = props => (
    <Edit hasEdit={false} toolbar={<></>} {...props}>
        <SimpleForm>
            <TextField source="contractor" />
            <TextField source="state" />
            <TimeAndDateField source="start_date" />
        </SimpleForm>
    </Edit>
);
