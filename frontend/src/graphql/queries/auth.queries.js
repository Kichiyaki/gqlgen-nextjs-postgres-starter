import gql from 'graphql-tag';

export const ME = gql`
  query currentUser {
    me {
      id
      login
      role
      activated
      email
      createdAt
      updatedAt
    }
  }
`;
