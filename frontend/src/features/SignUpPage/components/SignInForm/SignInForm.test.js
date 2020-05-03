import { render, fireEvent, waitFor } from '@testing-library/react';
import i18n from 'i18next';
import { ME } from '@graphql/queries/auth.queries';
import TestLayout from '@utils/test_utils/TestLayout';
import createApolloClient from '@utils/test_utils/createApolloClient';
import { SIGN_IN_PAGE } from '@config/namespaces';
import {
  MAXIMUM_LOGIN_LENGTH,
  MAXIMUM_PASSWORD_LENGTH,
  MINIMUM_LOGIN_LENGTH,
  MINIMUM_PASSWORD_LENGTH,
} from '@config/sign-up-policy';
import { INPUT_IDS, SIGN_IN_MUTATION } from './constants';
import SignInForm from './SignInForm';

let showDialog, setMessage, setSeverity, push;
const credentials = {
  login: 'Logineszkowys',
  password: '123sssDD22sd',
};
const t = i18n.getFixedT(null, SIGN_IN_PAGE);

const renderCmp = (mocks = []) => {
  showDialog = jest.fn();
  setMessage = jest.fn();
  setSeverity = jest.fn();
  push = jest.fn();
  const client = createApolloClient(mocks);

  return {
    ...render(
      <TestLayout routerProps={{ push }}>
        <SignInForm
          showDialog={showDialog}
          setMessage={setMessage}
          setSeverity={setSeverity}
        />
      </TestLayout>
    ),
    client,
  };
};

describe('features > SignUpPage > components > SignInForm', () => {
  test('login and password are required', async () => {
    const { getAllByDisplayValue, queryByText } = renderCmp();

    getAllByDisplayValue('').forEach((el) => {
      fireEvent.change(el, { target: { value: '' } });
      fireEvent.blur(el);
    });

    await waitFor(() => {
      [
        t('signInForm.errors.validation.mustProvideLogin'),
        t('signInForm.errors.validation.mustProvidePassword'),
      ].forEach((text) => {
        expect(queryByText(text)).toBeInTheDocument();
      });
    });
  });

  test('successful submit', async () => {
    const mocks = [
      {
        request: {
          query: SIGN_IN_MUTATION,
          variables: credentials,
        },
        result: {
          data: {
            signin: {
              id: 1,
            },
          },
        },
      },
      {
        request: {
          query: ME,
        },
        result: {
          data: {
            me: {
              id: 1,
            },
          },
        },
      },
    ];

    const { getAllByDisplayValue, queryByText, getByDisplayValue } = renderCmp(
      mocks
    );

    getAllByDisplayValue('').forEach((el) => {
      fireEvent.change(el, {
        target: { value: credentials[el.id] },
      });
      fireEvent.blur(el);
    });

    await waitFor(() => {
      expect(getByDisplayValue(credentials.login)).toBeInTheDocument();
    });

    fireEvent.click(queryByText(t('signInForm.submitButton')));
    await waitFor(() => {
      expect(setMessage).not.toHaveBeenCalled();
      expect(setSeverity).not.toHaveBeenCalled();
    });
  });

  test('should correctly call showDialog', async () => {
    const { getByText } = renderCmp();
    fireEvent.click(getByText(t('signInForm.forgotPassword')));
    await waitFor(() => {
      expect(showDialog).toHaveBeenCalledTimes(1);
    });
  });
});
