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

export const GENERATE_NEW_ACTIVATION_TOKEN_FOR_CURRENT_USER = gql`
  mutation generateNewActivationTokenForCurrentUserMutation {
    generateNewActivationTokenForCurrentUser
  }
`;
