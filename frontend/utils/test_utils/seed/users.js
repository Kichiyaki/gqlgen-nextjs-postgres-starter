import constants from "../../../config/constants";

const newUser = ({
  login = "",
  role = constants.ROLES.defaultRole,
  createdAt = new Date(0),
  email = "",
  activated = false,
  slug = "",
  id = 0
} = {}) => ({
  login,
  role,
  createdAt,
  email,
  activated,
  slug,
  id
});

export const users = [
  newUser({
    login: "Kiszkowaty",
    role: constants.ROLES.administrativeRole,
    email: "kicha@gmail.com",
    activated: true,
    slug: "1-kiszkowaty",
    id: 1
  }),
  newUser({
    login: "Kichowaty",
    email: "kichowaty@gmail.com",
    slug: "1-kichowaty",
    id: 2
  })
];
