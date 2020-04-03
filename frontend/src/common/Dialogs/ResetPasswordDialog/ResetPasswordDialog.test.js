import React from 'react';
import { render, fireEvent, waitFor } from '@testing-library/react';
import i18n from 'i18next';
import { ApolloProvider } from 'react-apollo';

import { GENERATE_NEW_RESET_PASSWORD_TOKEN_MUTATION } from './constants';
import ResetPasswordDialog from './ResetPasswordDialog';
import createApolloClient from '@utils/test_utils/createApolloClient';
import MockRouter from '@utils/test_utils/MockRouter';
import { COMMON } from '@config/namespaces';

let onClose, setMessage, setSeverity;
const t = i18n.getFixedT(null, COMMON);

const renderResetPasswordDialog = (mocks = []) => {
  const client = createApolloClient({ mocks });
  onClose = jest.fn();
  setMessage = jest.fn();
  setSeverity = jest.fn();
  return {
    ...render(
      <MockRouter>
        <ApolloProvider client={client}>
          <ResetPasswordDialog
            onClose={onClose}
            open
            setMessage={setMessage}
            setSeverity={setSeverity}
          />
        </ApolloProvider>
      </MockRouter>
    ),
    client
  };
};

describe('common > Dialogs > ResetPasswordDialog', () => {
  test('email is required', async () => {
    const { getByText } = renderResetPasswordDialog();
    fireEvent.click(getByText(t('resetPasswordDialog.submitButton')));

    await waitFor(() =>
      expect(setMessage).toHaveBeenCalledWith(
        t('resetPasswordDialog.errors.validation.mustProvideEmail')
      )
    );
  });

  test('email address must be valid', async () => {
    const { getByText, getByDisplayValue } = renderResetPasswordDialog();
    const email = 'asdf';
    fireEvent.change(getByDisplayValue(''), { target: { value: email } });

    await waitFor(() => getByDisplayValue(email));

    fireEvent.click(getByText(t('resetPasswordDialog.submitButton')));

    await waitFor(() =>
      expect(setMessage).toHaveBeenCalledWith(
        t('resetPasswordDialog.errors.validation.invalidEmail', { email })
      )
    );
  });

  test('success', async () => {
    const email = 'test@test.com';
    const mocks = [
      {
        request: {
          query: GENERATE_NEW_RESET_PASSWORD_TOKEN_MUTATION,
          variables: { email }
        },
        result: {
          data: {
            generateNewResetPasswordToken: 'Success'
          }
        }
      }
    ];

    const { getByText, getByDisplayValue } = renderResetPasswordDialog(mocks);
    fireEvent.change(getByDisplayValue(''), { target: { value: email } });

    await waitFor(() => getByDisplayValue(email));

    fireEvent.click(getByText(t('resetPasswordDialog.submitButton')));

    await waitFor(() => {
      expect(setMessage).toHaveBeenLastCalledWith(
        t('resetPasswordDialog.success')
      );
      expect(setSeverity).toHaveBeenLastCalledWith('success');
      expect(onClose).toHaveBeenCalledTimes(1);
    });
  });
});
