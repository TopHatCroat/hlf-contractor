import React from 'react';
import PropTypes from 'prop-types';
import { DateField } from 'react-admin';

const UserField = ({ source, record = {} }) => {
    const data = record[source].split(":");
    return <>
        <DateField re source={source} options={{ weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }} />
        <DateField style={{fontWeight: 600}}>{data[1]}</DateField>
    </>;
};

UserField.propTypes = {
    label: PropTypes.string,
    record: PropTypes.object,
    source: PropTypes.string.isRequired,
};

export default UserField;
