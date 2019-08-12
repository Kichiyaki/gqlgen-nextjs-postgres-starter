import gql from "graphql-tag";

export const RESET_PASSWORD_QUERY = gql`
  query resetPasswordQuery($id: Int!, $token: String!) {
    resetPassword(id: $id, token: $token)
  }
`;
