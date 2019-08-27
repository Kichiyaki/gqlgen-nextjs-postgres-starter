import gql from "graphql-tag";

export const ACTIVATE_USER_ACCOUNT_MUTATION = gql`
  mutation activateUserAccountMutation($id: Int!, $token: String!) {
    activateUserAccount(id: $id, token: $token) {
      id
      login
      activated
    }
  }
`;

export const CHANGE_PASSWORD_MUTATION = gql`
  mutation changePasswordMutation(
    $currentPassword: String!
    $newPassword: String!
  ) {
    changePassword(currentPassword: $currentPassword, newPassword: $newPassword)
  }
`;
