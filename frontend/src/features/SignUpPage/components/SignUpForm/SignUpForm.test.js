import { render, fireEvent, waitFor } from '@testing-library/react';
import i18n from 'i18next';
import { pick } from 'lodash';
import { ME } from '@graphql/queries/auth.queries';
import TestLayout from '@utils/test_utils/TestLayout';
import createApolloClient from '@utils/test_utils/createApolloClient';
import { SIGN_UP_PAGE } from '@config/namespaces';
import {
  MAXIMUM_LOGIN_LENGTH,
  MAXIMUM_PASSWORD_LENGTH,
  MINIMUM_LOGIN_LENGTH,
  MINIMUM_PASSWORD_LENGTH
} from '@config/sign-up-policy';
import { INPUT_IDS, SIGN_UP_MUTATION, PROPS_TO_SEND } from './constants';
import SignUpForm from './SignUpForm';

let showDialog, setMessage, setSeverity, push;
const user = {
  login: 'Logineszkowys',
  password: '123sssDD22sd',
  confirmPassword: '123sssDD22sd',
  email: 'email@email.com'
};
const t = i18n.getFixedT(null, SIGN_UP_PAGE);

const renderCmp = (mocks = []) => {
  showDialog = jest.fn();
  setMessage = jest.fn();
  setSeverity = jest.fn();
  push = jest.fn();
  const client = createApolloClient(mocks);

  return {
    ...render(
      <TestLayout routerProps={{ push }}>
        <SignUpForm
          showDialog={showDialog}
          setMessage={setMessage}
          setSeverity={setSeverity}
        />
      </TestLayout>
    ),
    client
  };
};

describe('features > SignUpPage > components > SignUpForm', () => {
  test('login, email and password are required', async () => {
    const { getAllByDisplayValue, queryAllByText, queryByText } = renderCmp();

    getAllByDisplayValue('').forEach(el => {
      fireEvent.change(el, { target: { value: '' } });
      fireEvent.blur(el);
    });

    await waitFor(() => {
      [
        t('signUpForm.errors.validation.mustProvideLogin'),
        t('signUpForm.errors.validation.mustProvideEmail'),
        [t('signUpForm.errors.validation.mustProvidePassword'), 2]
      ].forEach(test => {
        if (Array.isArray(test)) {
          expect(queryAllByText(test[0])).toHaveLength(test[1]);
        } else {
          expect(queryByText(test)).toBeInTheDocument();
        }
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
          t('signUpForm.errors.validation.minimumLoginLength', {
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
          t('signUpForm.errors.validation.maximumLoginLength', {
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
          t('signUpForm.errors.validation.minimumPasswordLength', {
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
          t('signUpForm.errors.validation.maximumPasswordLength', {
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
            'signUpForm.errors.validation.passwordMustContainAtLeastOneLowercase'
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
            'signUpForm.errors.validation.passwordMustContainAtLeastOneUppercase'
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
          t('signUpForm.errors.validation.passwordMustContainAtLeastOneDigit')
        )
      ).toBeInTheDocument()
    );
  });

  test('passwords must be the same', async () => {
    const { getAllByDisplayValue, queryByText } = renderCmp();

    getAllByDisplayValue('').forEach(el => {
      if (el.id === INPUT_IDS.PASSWORD) {
        fireEvent.change(el, { target: { value: 'asdasdasdasASDdaadsa' } });
        fireEvent.blur(el);
      }
      if (el.id === INPUT_IDS.CONFIRM_PASSWORD) {
        fireEvent.change(el, { target: { value: 'sdsdsds' } });
        fireEvent.blur(el);
      }
    });

    await waitFor(() =>
      expect(
        queryByText(t('signUpForm.errors.validation.passwordsAreNotTheSame'))
      ).toBeInTheDocument()
    );
  });

  test('should not allow entering the invalid email address', async () => {
    const { getAllByDisplayValue, queryByText, queryAllByText } = renderCmp();
    const email = 'asdfffsasd';

    getAllByDisplayValue('').forEach(el => {
      if (el.id === INPUT_IDS.EMAIL) {
        fireEvent.change(el, { target: { value: email } });
        fireEvent.blur(el);
      }
    });

    await waitFor(() => {
      [t('signUpForm.errors.validation.invalidEmail')].forEach(text => {
        expect(queryByText(text)).toBeInTheDocument();
      });
    });
  });

  test('successful submit', async () => {
    const mocks = [
      {
        request: {
          query: SIGN_UP_MUTATION,
          variables: pick(user, PROPS_TO_SEND)
        },
        result: {
          data: {
            signup: {
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
        target: { value: user[el.id] }
      });
      fireEvent.blur(el);
    });

    await waitFor(() => {
      expect(getByDisplayValue(user.login)).toBeInTheDocument();
    });

    fireEvent.click(queryByText(t('signUpForm.submitButton')));
    await waitFor(() => {
      expect(setMessage).not.toHaveBeenCalled();
      expect(setSeverity).not.toHaveBeenCalled();
    });
  });

  test('should correctly call showDialog', async () => {
    const { getByText } = renderCmp();
    fireEvent.click(getByText(t('signUpForm.forgotPassword')));
    await waitFor(() => {
      expect(showDialog).toHaveBeenCalledTimes(1);
    });
  });
});
