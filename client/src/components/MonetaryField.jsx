import React from 'react';
import PropTypes from 'prop-types';

const MonetaryField = ({ source, record = {} }) => {
    if(record[source] === undefined || record[source] === null) {
        return <div>TBA</div>
    }

    const amount = new Intl.NumberFormat().format(record[source] / 100);
    return <div>{`${amount} €}`}</div>
};

MonetaryField.propTypes = {
    label: PropTypes.string,
    record: PropTypes.object,
    source: PropTypes.string.isRequired,
};

export default MonetaryField;
