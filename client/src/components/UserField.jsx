import React from 'react';
import PropTypes from 'prop-types';

const UserField = ({ source, record = {} }) => {
    const data = record[source].split(":");
    return <>
        <div>{data[0]}</div>
        <div style={{fontWeight: 600}}>{data[1]}</div>
    </>;
};

UserField.propTypes = {
    label: PropTypes.string,
    record: PropTypes.object,
    source: PropTypes.string.isRequired,
};

export default UserField;
