import React from 'react';
import { useQuery, useMutation } from '@apollo/react-hooks';
import isGraphQLError from '@graphql/isGraphQLError';
import { ME } from '@graphql/queries/auth.queries';
import { useTranslation } from '@libs/i18n';
import useSnackbar, { SEVERITY } from '@libs/useSnackbar';
import restrictionWrapper from '@hocs/restrictionWrapper.hoc';
import { USER_PAGE, COMMON } from '@config/namespaces';
import { GENERATE_NEW_ACTIVATION_TOKEN_MUTATION } from './constants';

import {
  Container,
  Card,
  CardHeader,
  CardContent,
  Button,
  Snackbar,
} from '@material-ui/core';
import { Alert } from '@material-ui/lab';
import AppLayout from '@common/AppLayout/AppLayout';

function AccountPage() {
  const { data } = useQuery(ME, { fetchPolicy: 'cache-only' });
  const user = data && data.me ? data.me : {};
  const [
    generateNewActivationTokenMutation,
    { loading: generateNewActivationTokenLoading },
  ] = useMutation(GENERATE_NEW_ACTIVATION_TOKEN_MUTATION, {});
  const { t } = useTranslation(USER_PAGE.SETTINGS_PAGE.ACCOUNT_PAGE);
  const {
    snackbarProps,
    alertProps,
    message,
    setMessage,
    setSeverity,
  } = useSnackbar();

  const handleNewActivationTokenGeneration = async () => {
    try {
      await generateNewActivationTokenMutation();
      setMessage(t('generateNewActivationToken_success'));
      setSeverity(SEVERITY.SUCCESS);
    } catch (error) {
      if (isGraphQLError(error)) {
        setMessage(error.graphQLErrors[0].message);
      } else {
        setMessage(t('generateNewActivationToken_error'));
      }
      setSeverity(SEVERITY.ERROR);
    }
  };

  return (
    <AppLayout>
      <Container>
        <Card>
          <CardHeader
            title={
              user.activated
                ? t('accountActivated_true')
                : t('accountActivated_false')
            }
          />
          {!user.activated && (
            <CardContent>
              <Button
                onClick={handleNewActivationTokenGeneration}
                disabled={generateNewActivationTokenLoading}
              >
                {t('generateNewActivationToken')}
              </Button>
            </CardContent>
          )}
        </Card>
      </Container>
      <Snackbar {...snackbarProps}>
        <Alert {...alertProps}>{message}</Alert>
      </Snackbar>
    </AppLayout>
  );
}

AccountPage.getInitialProps = () => {
  return {
    namespacesRequired: [COMMON, USER_PAGE.SETTINGS_PAGE.ACCOUNT_PAGE],
  };
};

export default restrictionWrapper({ loggedIn: true })(AccountPage);
