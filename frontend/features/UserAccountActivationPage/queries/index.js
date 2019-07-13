import gql from "graphql-tag";

export const ACTIVATE_USER_ACCOUNT_QUERY = gql`
  query activateUserAccountQuery($id: Int!, $token: String!) {
    activateUserAccount(id: $id, token: $token) {
      id
      login
    }
  }
`;
