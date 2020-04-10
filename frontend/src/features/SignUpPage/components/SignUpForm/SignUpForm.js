import React from 'react';
import { pick } from 'lodash';
import { useMutation } from '@apollo/react-hooks';
import { func } from 'prop-types';
import { useFormik } from 'formik';
import * as Yup from 'yup';
import { useTranslation } from '@libs/i18n';
import usePrompt from '@libs/usePrompt';
import {
  MAXIMUM_LOGIN_LENGTH,
  MAXIMUM_PASSWORD_LENGTH,
  MINIMUM_LOGIN_LENGTH,
  MINIMUM_PASSWORD_LENGTH,
  CONTAIN_UPPERCASE,
  CONTAIN_DIGIT,
  CONTAIN_LOWERCASE,
} from '@config/sign-up-policy';
import { INPUT_IDS, PROPS_TO_SEND, SIGN_UP_MUTATION } from './constants';
import { SIGN_UP_PAGE } from '@config/namespaces';
import { SIGN_IN_PAGE } from '@config/routes';
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

export default function SignUpForm({ setMessage, setSeverity }) {
  const classes = useStyles();
  const { t } = useTranslation(SIGN_UP_PAGE);
  usePrompt('?');
  const [signUp] = useMutation(SIGN_UP_MUTATION, {
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
      confirmPassword: '',
      email: '',
    },
    onSubmit: async (user, { setSubmitting }) => {
      try {
        await signUp({ variables: { user: pick(user, PROPS_TO_SEND) } });
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
        .min(
          MINIMUM_LOGIN_LENGTH,
          t('signUpForm.errors.validation.minimumLoginLength', {
            count: MINIMUM_LOGIN_LENGTH,
          })
        )
        .max(
          MAXIMUM_LOGIN_LENGTH,
          t('signUpForm.errors.validation.maximumLoginLength', {
            count: MAXIMUM_LOGIN_LENGTH,
          })
        )
        .required(t('signUpForm.errors.validation.mustProvideLogin')),
      password: Yup.string()
        .trim()
        .required(t('signUpForm.errors.validation.mustProvidePassword'))
        .min(
          MINIMUM_PASSWORD_LENGTH,
          t('signUpForm.errors.validation.minimumPasswordLength', {
            count: MINIMUM_PASSWORD_LENGTH,
          })
        )
        .max(
          MAXIMUM_PASSWORD_LENGTH,
          t('signUpForm.errors.validation.maximumPasswordLength', {
            count: MAXIMUM_PASSWORD_LENGTH,
          })
        )
        .matches(
          CONTAIN_UPPERCASE,
          t(
            'signUpForm.errors.validation.passwordMustContainAtLeastOneUppercase'
          )
        )
        .matches(
          CONTAIN_LOWERCASE,
          t(
            'signUpForm.errors.validation.passwordMustContainAtLeastOneLowercase'
          )
        )
        .matches(
          CONTAIN_DIGIT,
          t('signUpForm.errors.validation.passwordMustContainAtLeastOneDigit')
        ),
      email: Yup.string()
        .trim()
        .email(t('signUpForm.errors.validation.invalidEmail'))
        .required(t('signUpForm.errors.validation.mustProvideEmail')),
      confirmPassword: Yup.string()
        .oneOf(
          [Yup.ref('password'), null],
          t('signUpForm.errors.validation.passwordsAreNotTheSame')
        )
        .required(t('signUpForm.errors.validation.mustProvidePassword')),
    }),
  });

  const defaultInputProps = {
    variant: 'outlined',
    margin: 'normal',
    required: true,
    fullWidth: true,
    onBlur: handleBlur,
    onChange: handleChange,
  };

  return (
    <Card>
      <CardHeader title={t('signUpForm.title')} />
      <CardContent>
        <TextField
          label={t('signUpForm.inputLabel.login')}
          name={INPUT_IDS.LOGIN}
          id={INPUT_IDS.LOGIN}
          value={values.login}
          error={touched.login && !!errors.login}
          helperText={touched.login && errors.login}
          {...defaultInputProps}
        />
        <TextField
          name={INPUT_IDS.EMAIL}
          id={INPUT_IDS.EMAIL}
          label={t('signUpForm.inputLabel.email')}
          value={values.email}
          error={touched.email && !!errors.email}
          helperText={touched.email && errors.email}
          {...defaultInputProps}
        />
        <TextField
          type="password"
          name={INPUT_IDS.PASSWORD}
          id={INPUT_IDS.PASSWORD}
          label={t('signUpForm.inputLabel.password')}
          value={values.password}
          error={touched.password && !!errors.password}
          helperText={touched.password && errors.password}
          {...defaultInputProps}
        />
        <TextField
          type="password"
          name={INPUT_IDS.CONFIRM_PASSWORD}
          id={INPUT_IDS.CONFIRM_PASSWORD}
          label={t('signUpForm.inputLabel.confirmPassword')}
          value={values.confirmPassword}
          error={touched.confirmPassword && !!errors.confirmPassword}
          helperText={touched.confirmPassword && errors.confirmPassword}
          {...defaultInputProps}
        />
        <Typography component="p">
          <Link href={SIGN_IN_PAGE}>
            {t('signUpForm.alreadyHaveAnAccount')}
          </Link>
        </Typography>
      </CardContent>
      <CardActions disableSpacing className={classes.actions}>
        <Button onClick={handleSubmit} disabled={isSubmitting}>
          {t('signUpForm.submitButton')}
        </Button>
      </CardActions>
    </Card>
  );
}
