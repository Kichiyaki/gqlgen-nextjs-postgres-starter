import gql from 'graphql-tag';

export const RESET_USER_PASSWORD_MUTATION = gql`
  query resetUserPasswordMutation($id: Int!, $token: String!) {
    resetUserPassword(id: $id, token: $token)
  }
`;
