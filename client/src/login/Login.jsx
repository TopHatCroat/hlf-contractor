import React, { Component } from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import Card from '@material-ui/core/Card';
import Avatar from '@material-ui/core/Avatar';
import {
    MuiThemeProvider,
    createMuiTheme,
    withStyles,
    createStyles,
} from '@material-ui/core/styles';
import LockIcon from '@material-ui/icons/Lock';
import { defaultTheme, Notification } from 'react-admin';
import LoginForm from './LoginForm';

const styles = (theme) =>
    createStyles({
        main: {
            display: 'flex',
            flexDirection: 'column',
            minHeight: '100vh',
            height: '1px',
            alignItems: 'center',
            justifyContent: 'flex-start',
            backgroundRepeat: 'no-repeat',
            backgroundSize: 'cover',
        },
        card: {
            minWidth: 300,
            marginTop: '6em',
        },
        avatar: {
            margin: '1em',
            display: 'flex',
            justifyContent: 'center',
        },
        icon: {
            backgroundColor: theme.palette.secondary[500],
        },
    });

class Login extends Component {
    theme = createMuiTheme(this.props.theme);
    containerRef = React.createRef();
    backgroundImageLoaded = false;

    updateBackgroundImage = () => {
        if (!this.backgroundImageLoaded && this.containerRef.current) {
            const { backgroundImage } = this.props;
            this.containerRef.current.style.backgroundImage = `url(${backgroundImage})`;
            this.backgroundImageLoaded = true;
        }
    };

    // Load background image asynchronously to speed up time to interactive
    lazyLoadBackgroundImage() {
        const { backgroundImage } = this.props;

        if (backgroundImage) {
            const img = new Image();
            img.onload = this.updateBackgroundImage;
            img.src = backgroundImage;
        }
    }

    componentDidMount() {
        this.lazyLoadBackgroundImage();
    }

    componentDidUpdate() {
        if (!this.backgroundImageLoaded) {
            this.lazyLoadBackgroundImage();
        }
    }

    render() {
        const {
            backgroundImage,
            classes,
            className,
            loginForm,
            staticContext,
            ...rest
        } = this.props;

        return (
            <MuiThemeProvider theme={this.theme}>
                <div
                    className={classnames(classes.main, className)}
                    {...rest}
                    ref={this.containerRef}
                >
                    <Card className={classes.card}>
                        <div className={classes.avatar}>
                            <Avatar className={classes.icon}>
                                <LockIcon />
                            </Avatar>
                        </div>
                        {loginForm}
                    </Card>
                    <Notification />
                </div>
            </MuiThemeProvider>
        );
    }
}

const EnhancedLogin = withStyles(styles)(Login);

EnhancedLogin.propTypes = {
    backgroundImage: PropTypes.string,
    loginForm: PropTypes.element,
    theme: PropTypes.object,
    staticContext: PropTypes.object,
};

EnhancedLogin.defaultProps = {
    backgroundImage: 'https://source.unsplash.com/random/1600x900/daily',
    theme: defaultTheme,
    loginForm: <LoginForm />,
};
export default EnhancedLogin;