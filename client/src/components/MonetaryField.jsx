import React from 'react';
import PropTypes from 'prop-types';
import TextField from "@material-ui/core/TextField";

const MonetaryField = ({ source, record = {} }) => {
    const amount = new Intl.NumberFormat().format(record[source] / 100);
    return (<TextField>{`${amount} €}`}</TextField>)
};

MonetaryField.propTypes = {
    label: PropTypes.string,
    record: PropTypes.object,
    source: PropTypes.string.isRequired,
};

export default MonetaryField;
