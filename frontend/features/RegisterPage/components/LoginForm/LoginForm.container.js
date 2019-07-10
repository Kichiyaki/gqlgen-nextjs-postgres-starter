import React from "react";
import { useMutation } from "@apollo/react-hooks";
import { Formik } from "formik";
import * as Yup from "yup";
import { object } from "prop-types";

import LoginFormCmp from "./LoginForm.component";
import { LOGIN_MUTATION } from "../../mutations";
import { showErrorMessage, showSuccessMessage } from "@services/toastify";
import { FETCH_CURRENT_USER_QUERY } from "@graphql/queries/user.queries";
import constants from "@config/constants";

const LoginForm = ({ translations }) => {
  const [loginMutation] = useMutation(LOGIN_MUTATION);

  const handleSubmit = async (values, { resetForm, setSubmitting }) => {
    try {
      await loginMutation({
        variables: { login: values.login, password: values.password },
        refetchQueries: [{ query: FETCH_CURRENT_USER_QUERY }],
        awaitRefetchQueries: true
      });
      resetForm();
      showSuccessMessage(translations.REGISTER_PAGE.loginForm.success);
    } catch (error) {
      if (error.graphQLErrors && error.graphQLErrors[0]) {
        showErrorMessage(error.graphQLErrors[0].message);
      } else {
        showErrorMessage(translations.REGISTER_PAGE.loginForm.errors.default);
      }
    }
    setSubmitting(false);
  };
  return (
    <Formik
      initialValues={{
        login: "",
        password: ""
      }}
      render={formikProps => (
        <LoginFormCmp {...formikProps} translations={translations} />
      )}
      onSubmit={handleSubmit}
      validationSchema={Yup.object().shape({
        login: Yup.string()
          .trim()
          .min(
            constants.VALIDATION.minimumLengthOfLogin,
            translations.REGISTER_PAGE.registerForm.errors.validation
              .minimumLengthOfLogin
          )
          .max(
            constants.VALIDATION.maximumLengthOfLogin,
            translations.REGISTER_PAGE.registerForm.errors.validation
              .maximumLengthOfLogin
          )
          .required(
            translations.REGISTER_PAGE.registerForm.errors.validation
              .mustProvideLogin
          ),
        password: Yup.string()
          .trim()
          .required(
            translations.REGISTER_PAGE.registerForm.errors.validation
              .mustProvidePassword
          )
          .min(
            constants.VALIDATION.minimumLengthOfPassword,
            translations.REGISTER_PAGE.registerForm.errors.validation
              .minimumLengthOfPassword
          )
          .max(
            constants.VALIDATION.maximumLengthOfPassword,
            translations.REGISTER_PAGE.registerForm.errors.validation
              .maximumLengthOfPassword
          )
          .matches(
            constants.REGEXES.containsUppercase,
            translations.REGISTER_PAGE.registerForm.errors.validation
              .passwordMustContainsOneUppercase
          )
          .matches(
            constants.REGEXES.containsLowercase,
            translations.REGISTER_PAGE.registerForm.errors.validation
              .passwordMustContainsOneLowercase
          )
          .matches(
            constants.REGEXES.containsDigit,
            translations.REGISTER_PAGE.registerForm.errors.validation
              .passwordMustContainsOneDigit
          )
      })}
    />
  );
};

LoginForm.propTypes = {
  translations: object.isRequired
};

export default LoginForm;
