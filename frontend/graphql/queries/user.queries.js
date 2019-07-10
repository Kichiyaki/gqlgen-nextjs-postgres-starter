import gql from "graphql-tag";

export const FETCH_CURRENT_USER_QUERY = gql`
  query fetchCurrentUserQuery {
    fetchCurrentUser {
      id
      login
      role
      activated
      email
      createdAt
    }
  }
`;
