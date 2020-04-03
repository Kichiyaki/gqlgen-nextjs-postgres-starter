import gql from 'graphql-tag';

export const INPUT_IDS = {
  LOGIN: 'login',
  PASSWORD: 'password'
};

export const SIGN_IN_MUTATION = gql`
  mutation signInMutation($login: String!, $password: String!) {
    signin(login: $login, password: $password) {
      id
    }
  }
`;
