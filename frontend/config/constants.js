export default {
  ROLES: {
    defaultRole: "Użytkownik",
    administrativeRole: "Administrator"
  },
  VALIDATION: {
    minimumLengthOfPassword: 8,
    maximumLengthOfPassword: 128,
    minimumLengthOfLogin: 3,
    maximumLengthOfLogin: 72
  },
  REGEXES: {
    containsUppercase: /[A-ZŻŹĆĄŚĘŁÓŃ]+/,
    containsLowercase: /[a-zzżźćńółęąś]+/,
    containsDigit: /\d+/
  },
  ROUTES: {
    root: "/",
    login: "/login",
    register: "/register",
    userPage: {
      settingsPage: {
        accountActivation: "/user/settings/activate",
        changePassword: "/user/settings/change-password"
      }
    }
  },
  AUTHOR_FULL_NAME: "Dawid Wysokiński",
  NAMESPACES: {
    common: "common",
    userPage: {
      settingsPage: {
        navigation: "user-page/settings-page/navigation"
      }
    }
  }
};
