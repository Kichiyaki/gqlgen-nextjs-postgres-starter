import React, { useState } from 'react';
import { bool, func } from 'prop-types';
import { useMutation } from '@apollo/react-hooks';
import isEmail from 'validator/lib/isEmail';
import { useTranslation } from '@libs/i18n';
import isGraphQLError from '@graphql/isGraphQLError';
import { COMMON } from '@config/namespaces';
import { GENERATE_NEW_RESET_PASSWORD_TOKEN_MUTATION } from './constants';

import {
  Dialog,
  Button,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField
} from '@material-ui/core';

const ResetPasswordDialog = ({ open, onClose, setMessage, setSeverity }) => {
  const [email, setEmail] = useState('');
  const { t } = useTranslation(COMMON);
  const [generateNewResetPasswordTokenMutation, { loading }] = useMutation(
    GENERATE_NEW_RESET_PASSWORD_TOKEN_MUTATION,
    {
      ignoreResults: true
    }
  );

  const handleChange = e => {
    setEmail(e.target.value);
  };

  const handleSubmit = async () => {
    if (email.trim().length === 0) {
      setMessage(t('resetPasswordDialog.errors.validation.mustProvideEmail'));
      setSeverity('error');
      return;
    }
    if (!isEmail(email)) {
      setMessage(t('resetPasswordDialog.errors.validation.invalidEmail'));
      setSeverity('error');
      return;
    }

    try {
      await generateNewResetPasswordTokenMutation({ variables: { email } });
      setMessage(t('resetPasswordDialog.success'));
      setSeverity('success');
      onClose();
    } catch (error) {
      if (isGraphQLError(error)) {
        setMessage(error.graphQLErrors[0].message);
      } else {
        setMessage(t('resetPasswordDialog.errors.default'));
      }
      setSeverity('error');
    }
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>{t('resetPasswordDialog.title')}</DialogTitle>
      <DialogContent>
        <TextField
          autoFocus
          margin="dense"
          label={t('resetPasswordDialog.inputLabel.email')}
          type="email"
          fullWidth
          value={email}
          onChange={handleChange}
          required
        />
      </DialogContent>
      <DialogActions>
        <Button disabled={loading} onClick={handleSubmit} color="secondary">
          {t('resetPasswordDialog.submitButton')}
        </Button>
        <Button disabled={loading} onClick={onClose} color="secondary">
          {t('resetPasswordDialog.cancelButton')}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

ResetPasswordDialog.propTypes = {
  open: bool.isRequired,
  onClose: func.isRequired,
  setMessage: func.isRequired,
  setSeverity: func.isRequired
};

export default ResetPasswordDialog;
