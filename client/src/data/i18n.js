import englishMessages from 'ra-language-english';

const extendedEnglishMessages = {
    ...englishMessages,
    app: {
        auth: {
            register: "Register",
            register_failed: "Registration failed",
            login: "Log in",
            login_failed: "Login failed",
        }
    },
    "not authorized": "Login failed"
};

const messages = {
    en: extendedEnglishMessages,
};

const i18nProvider = locale => messages[locale];

export default i18nProvider;
