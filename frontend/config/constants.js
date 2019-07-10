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
    register: "/register"
  },
  AUTHOR_FULL_NAME: "Dawid Wysokiński"
};
