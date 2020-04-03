import React, { useState } from 'react';
import restrictionWrapper from '@hocs/restrictionWrapper.hoc';
import useSnackbar from '@libs/useSnackbar';
import { COMMON, SIGN_UP_PAGE, SIGN_IN_PAGE } from '@config/namespaces';
import * as routes from '@config/routes';
import { Container, Snackbar } from '@material-ui/core';
import { Alert } from '@material-ui/lab';
import AppLayout from '@common/AppLayout/AppLayout';
import ResetPasswordDialog from '@common/Dialogs/ResetPasswordDialog/ResetPasswordDialog';
import SignInForm from './components/SignInForm/SignInForm';
import SignUpForm from './components/SignUpForm/SignUpForm';

function SignUpPage({ pathname }) {
  const [isOpen, setIsOpen] = useState(false);
  const { snackbarProps, alertProps, message, ...snackbarBag } = useSnackbar();

  const showDialog = () => setIsOpen(true);
  const hideDialog = () => setIsOpen(false);

  return (
    <AppLayout>
      <Container maxWidth="sm">
        {routes.SIGN_IN_PAGE === pathname ? (
          <SignInForm {...snackbarBag} showDialog={showDialog} />
        ) : (
          <SignUpForm {...snackbarBag} showDialog={showDialog} />
        )}
      </Container>
      {isOpen && (
        <ResetPasswordDialog
          open={isOpen}
          onClose={hideDialog}
          {...snackbarBag}
        />
      )}
      <Snackbar {...snackbarProps}>
        <Alert {...alertProps}>{message}</Alert>
      </Snackbar>
    </AppLayout>
  );
}

SignUpPage.getInitialProps = ({ pathname }) => {
  const props = {
    namespacesRequired: [COMMON],
    pathname
  };
  if (pathname === routes.SIGN_IN_PAGE) {
    props.namespacesRequired.push(SIGN_IN_PAGE);
  } else {
    props.namespacesRequired.push(SIGN_UP_PAGE);
  }
  return props;
};

export default restrictionWrapper({ loggedOut: true })(SignUpPage);
