import gql from 'graphql-tag';

export const GENERATE_NEW_RESET_PASSWORD_TOKEN_MUTATION = gql`
  mutation generateNewResetPasswordTokenMutation($email: String!) {
    generateNewResetPasswordToken(email: $email)
  }
`;
