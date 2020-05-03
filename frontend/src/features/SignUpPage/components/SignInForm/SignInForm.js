import React from 'react';
import { func } from 'prop-types';
import { useMutation } from '@apollo/react-hooks';
import { useFormik } from 'formik';
import * as Yup from 'yup';
import { useTranslation } from '@libs/i18n';
import usePrompt from '@libs/usePrompt';
import { INPUT_IDS, SIGN_IN_MUTATION } from './constants';
import { SIGN_IN_PAGE } from '@config/namespaces';
import { SIGN_UP_PAGE } from '@config/routes';
import { ME } from '@graphql/queries/auth.queries';
import isGraphQLError from '@graphql/isGraphQLError';

import { makeStyles } from '@material-ui/core/styles';
import {
  Card,
  CardHeader,
  CardContent,
  CardActions,
  TextField,
  Button,
  Typography,
} from '@material-ui/core';
import Link from '@common/Link/Link';

const useStyles = makeStyles(() => ({
  actions: {
    justifyContent: 'flex-end',
  },
}));

export default function SignInForm({ showDialog, setMessage, setSeverity }) {
  const classes = useStyles();
  const { t } = useTranslation(SIGN_IN_PAGE);
  usePrompt('?');
  const [signIn] = useMutation(SIGN_IN_MUTATION, {
    ignoreResults: true,
    awaitRefetchQueries: true,
    refetchQueries: [{ query: ME }],
  });
  const {
    values,
    handleChange,
    handleBlur,
    touched,
    errors,
    handleSubmit,
    isSubmitting,
  } = useFormik({
    initialValues: {
      login: '',
      password: '',
    },
    onSubmit: async (credentials, { setSubmitting }) => {
      try {
        await signIn({ variables: credentials });
      } catch (error) {
        if (isGraphQLError(error)) {
          setMessage(error.graphQLErrors[0].message);
        } else {
          setMessage(t('signInForm.errors.default'));
        }
        setSeverity('error');
        setSubmitting(false);
      }
    },
    validationSchema: Yup.object().shape({
      login: Yup.string()
        .trim()
        .required(t('signInForm.errors.validation.mustProvideLogin')),
      password: Yup.string()
        .trim()
        .required(t('signInForm.errors.validation.mustProvidePassword')),
    }),
  });

  return (
    <Card>
      <CardHeader title={t('signInForm.title')} />
      <CardContent>
        <TextField
          variant="outlined"
          margin="normal"
          required
          fullWidth
          label={t('signInForm.inputLabel.login')}
          name={INPUT_IDS.LOGIN}
          id={INPUT_IDS.LOGIN}
          value={values.login}
          onBlur={handleBlur}
          onChange={handleChange}
          error={touched.login && !!errors.login}
          helperText={touched.login && errors.login}
        />
        <TextField
          type="password"
          variant="outlined"
          margin="normal"
          required
          fullWidth
          name={INPUT_IDS.PASSWORD}
          id={INPUT_IDS.PASSWORD}
          label={t('signInForm.inputLabel.password')}
          value={values.password}
          onBlur={handleBlur}
          onChange={handleChange}
          error={touched.password && !!errors.password}
          helperText={touched.password && errors.password}
        />
        <Typography component="p">
          <Link href={SIGN_UP_PAGE}>{t('signInForm.dontHaveAnAccount')}</Link>
        </Typography>
        <Typography
          style={{ cursor: 'pointer' }}
          onClick={showDialog}
          component="p"
          color="secondary"
        >
          {t('signInForm.forgotPassword')}
        </Typography>
      </CardContent>
      <CardActions disableSpacing className={classes.actions}>
        <Button onClick={handleSubmit} disabled={isSubmitting}>
          {t('signInForm.submitButton')}
        </Button>
      </CardActions>
    </Card>
  );
}

SignInForm.propTypes = {
  showDialog: func.isRequired,
};
