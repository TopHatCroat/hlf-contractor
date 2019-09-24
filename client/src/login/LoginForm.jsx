import React from 'react';
import PropTypes from 'prop-types';
import { Field, reduxForm } from 'redux-form';
import { connect } from 'react-redux';
import compose from 'recompose/compose';
import { CardActions, DialogTitle, Button, TextField } from "@material-ui/core";
import CircularProgress from '@material-ui/core/CircularProgress';
import { withStyles, createStyles } from '@material-ui/core/styles';
import { withTranslate, userLogin } from 'react-admin';

const styles = ({ spacing }) =>
    createStyles({
        form: {
            padding: '0 1em 1em 1em',
        },
        input: {
            marginTop: '1em',
        },
        button: {
            width: '100%',
        },
        icon: {
            marginRight: spacing.unit,
        },
    });

// see http://redux-form.com/6.4.3/examples/material-ui/
const renderInput = ({
                         meta: { touched, error } = { touched: false, error: '' }, // eslint-disable-line react/prop-types
                         input: { ...inputProps }, // eslint-disable-line react/prop-types
                         ...props
                     }) => (
    <TextField
        error={!!(touched && error)}
        helperText={touched && error}
        {...inputProps}
        {...props}
        fullWidth
    />
);

const login = (auth, dispatch, { redirectTo }) =>
    dispatch(userLogin(auth, redirectTo));

class LoginForm extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            isRegister: false
        }
    }

    toggleRegister = () => {
        const { isRegister } = this.state;
        const { change } = this.props;
        this.setState({
            isRegister: !isRegister
        });
        change("kind", !isRegister ? "register" : "login");
    };

    render() {
        const { classes, isLoading, translate, handleSubmit } = this.props;
        const { isRegister } = this.state;

        return (
            <form onSubmit={handleSubmit(login)}>
                <DialogTitle>
                    {isRegister ? translate('app.auth.register') : translate('app.auth.login')}
                </DialogTitle>
                <div className={classes.form}>
                    <div className={classes.input}>
                        <Field
                            autoFocus
                            id="username"
                            name="username"
                            component={renderInput}
                            label={translate('ra.auth.username')}
                            disabled={isLoading}
                        />
                    </div>
                    <div className={classes.input}>
                        <Field
                            id="password"
                            name="password"
                            component={renderInput}
                            label={translate('ra.auth.password')}
                            type="password"
                            disabled={isLoading}
                        />
                    </div>
                    <div>
                        <Field
                            id="isRegister"
                            component={renderInput}
                            name="kind"
                            type="hidden"
                            style={{ height: 0 }}
                        />
                    </div>
                </div>
                <CardActions>
                    <Button
                        type="button"
                        variant="text"
                        color="secondary"
                        hidden={!isRegister}
                        className={classes.button}
                        onClick={this.toggleRegister}
                    >
                        {!isRegister ? translate('app.auth.register') : translate('app.auth.login')}
                    </Button>
                    <Button
                        variant="outlined"
                        type="submit"
                        color="primary"
                        disabled={isLoading}
                        className={classes.button}
                    >
                        {isLoading && (
                            <CircularProgress
                                className={classes.icon}
                                size={18}
                                thickness={2}
                            />
                        )}
                        {isRegister ? translate('app.auth.register') : translate('app.auth.login')}
                    </Button>
                </CardActions>
            </form>
        )
    }
}

const mapStateToProps = (state) => ({
    isLoading: state.admin.loading > 0,
});

const enhance = compose(
        withStyles(styles),
        withTranslate,
        connect(mapStateToProps),
        reduxForm({
            form: 'signIn',
            validate: (values, props) => {
                const errors = { username: '', password: '' };
                const { translate } = props;
                if (!values.username) {
                    errors.username = translate('ra.validation.required');
                }
                if (!values.password) {
                    errors.password = translate('ra.validation.required');
                }
                return errors;
            },
        })
);

const EnhancedLoginForm = enhance(LoginForm);

EnhancedLoginForm.propTypes = {
    redirectTo: PropTypes.string,
};

export default EnhancedLoginForm;