import gql from 'graphql-tag';

export const GENERATE_NEW_ACTIVATION_TOKEN_MUTATION = gql`
  mutation generateNewActivationTokenMutation {
    generateNewActivationTokenForMe
  }
`;
