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
  MINIMUM_PASSWORD_LENGTH
} from '@config/sign-up-policy';
import { INPUT_IDS, SIGN_IN_MUTATION } from './constants';
import SignInForm from './SignInForm';

let showDialog, setMessage, setSeverity, push;
const credentials = {
  login: 'Logineszkowys',
  password: '123sssDD22sd'
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
    client
  };
};

describe('features > SignUpPage > components > SignInForm', () => {
  test('login and password are required', async () => {
    const { getAllByDisplayValue, queryByText } = renderCmp();

    getAllByDisplayValue('').forEach(el => {
      fireEvent.change(el, { target: { value: '' } });
      fireEvent.blur(el);
    });

    await waitFor(() => {
      [
        t('signInForm.errors.validation.mustProvideLogin'),
        t('signInForm.errors.validation.mustProvidePassword')
      ].forEach(text => {
        expect(queryByText(text)).toBeInTheDocument();
      });
    });
  });

  test(`login length should be between ${MINIMUM_LOGIN_LENGTH} and ${MAXIMUM_LOGIN_LENGTH} characters`, async () => {
    const {
      getAllByDisplayValue,
      getByDisplayValue,
      queryByText
    } = renderCmp();
    let value = 'a';

    getAllByDisplayValue('').forEach(el => {
      if (el.id === INPUT_IDS.LOGIN) {
        fireEvent.change(el, { target: { value } });
        fireEvent.blur(el);
      }
    });

    await waitFor(() =>
      expect(
        queryByText(
          t('signInForm.errors.validation.minimumLoginLength', {
            count: MINIMUM_LOGIN_LENGTH
          })
        )
      ).toBeInTheDocument()
    );

    const el = getByDisplayValue(value);
    for (let k = 1; k <= MAXIMUM_LOGIN_LENGTH + 5; k++) {
      value += 'a';
    }

    fireEvent.change(el, {
      target: {
        value
      }
    });
    fireEvent.blur(el);

    await waitFor(() =>
      expect(
        queryByText(
          t('signInForm.errors.validation.maximumLoginLength', {
            count: MAXIMUM_LOGIN_LENGTH
          })
        )
      ).toBeInTheDocument()
    );
  });

  test(`password length should be between ${MINIMUM_PASSWORD_LENGTH} and ${MAXIMUM_PASSWORD_LENGTH} characters`, async () => {
    const {
      getAllByDisplayValue,
      queryByText,
      getByDisplayValue
    } = renderCmp();
    let value = 'asasd';

    getAllByDisplayValue('').forEach(el => {
      if (el.id === INPUT_IDS.PASSWORD) {
        fireEvent.change(el, { target: { value } });
        fireEvent.blur(el);
      }
    });

    await waitFor(() =>
      expect(
        queryByText(
          t('signInForm.errors.validation.minimumPasswordLength', {
            count: MINIMUM_PASSWORD_LENGTH
          })
        )
      ).toBeInTheDocument()
    );

    const el = getByDisplayValue(value);
    for (let k = 1; k <= MAXIMUM_PASSWORD_LENGTH + 5; k++) {
      value += 'a';
    }

    fireEvent.change(el, {
      target: {
        value
      }
    });
    fireEvent.blur(el);

    await waitFor(() =>
      expect(
        queryByText(
          t('signInForm.errors.validation.maximumPasswordLength', {
            count: MAXIMUM_PASSWORD_LENGTH
          })
        )
      ).toBeInTheDocument()
    );
  });

  test('password must contain at least 1 lowercase', async () => {
    const { getAllByDisplayValue, queryByText } = renderCmp();

    getAllByDisplayValue('').forEach(el => {
      if (el.id === INPUT_IDS.PASSWORD) {
        fireEvent.change(el, { target: { value: 'ASDASDAASDASDASDA' } });
        fireEvent.blur(el);
      }
    });

    await waitFor(() =>
      expect(
        queryByText(
          t(
            'signInForm.errors.validation.passwordMustContainAtLeastOneLowercase'
          )
        )
      ).toBeInTheDocument()
    );
  });

  test('password must contain at least 1 uppercase', async () => {
    const { getAllByDisplayValue, queryByText } = renderCmp();

    getAllByDisplayValue('').forEach(el => {
      if (el.id === INPUT_IDS.PASSWORD) {
        fireEvent.change(el, { target: { value: 'asdasdasdasdaadsa' } });
        fireEvent.blur(el);
      }
    });

    await waitFor(() =>
      expect(
        queryByText(
          t(
            'signInForm.errors.validation.passwordMustContainAtLeastOneUppercase'
          )
        )
      ).toBeInTheDocument()
    );
  });

  test('password must contain at least 1 digit', async () => {
    const { getAllByDisplayValue, queryByText } = renderCmp();

    getAllByDisplayValue('').forEach(el => {
      if (el.id === INPUT_IDS.PASSWORD) {
        fireEvent.change(el, { target: { value: 'asdasdasdasASDdaadsa' } });
        fireEvent.blur(el);
      }
    });

    await waitFor(() =>
      expect(
        queryByText(
          t('signInForm.errors.validation.passwordMustContainAtLeastOneDigit')
        )
      ).toBeInTheDocument()
    );
  });

  test('successful submit', async () => {
    const mocks = [
      {
        request: {
          query: SIGN_IN_MUTATION,
          variables: credentials
        },
        result: {
          data: {
            signin: {
              id: 1
            }
          }
        }
      },
      {
        request: {
          query: ME
        },
        result: {
          data: {
            me: {
              id: 1
            }
          }
        }
      }
    ];

    const { getAllByDisplayValue, queryByText, getByDisplayValue } = renderCmp(
      mocks
    );

    getAllByDisplayValue('').forEach(el => {
      fireEvent.change(el, {
        target: { value: credentials[el.id] }
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
