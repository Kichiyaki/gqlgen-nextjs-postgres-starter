import gql from 'graphql-tag';

export const INPUT_IDS = {
  LOGIN: 'login',
  PASSWORD: 'password',
  CONFIRM_PASSWORD: 'confirmPassword',
  EMAIL: 'email'
};

export const PROPS_TO_SEND = [
  INPUT_IDS.LOGIN,
  INPUT_IDS.PASSWORD,
  INPUT_IDS.EMAIL
];

export const SIGN_UP_MUTATION = gql`
  mutation signupMutation($user: UserInput!) {
    signup(user: $user) {
      id
    }
  }
`;
