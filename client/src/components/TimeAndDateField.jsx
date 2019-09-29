import React from 'react';
import PropTypes from 'prop-types';

const TimeAndDateField = ({ source, record = {} }) => {
    if(record[source] === undefined || record[source] === null) {
        return <div>TBA</div>
    }

    const dateAndTime = new Date(record[source]);
    return <>
        <div>{dateAndTime.toLocaleDateString("hr-HR", { weekday: 'short', year: 'numeric', month: 'short', day: 'numeric' })}</div>
        <div style={{fontWeight: 600}}>{dateAndTime.toLocaleTimeString("hr-HR")}</div>
    </>;
};

TimeAndDateField.propTypes = {
    label: PropTypes.string,
    record: PropTypes.object,
    source: PropTypes.string.isRequired,
};

export default TimeAndDateField;
