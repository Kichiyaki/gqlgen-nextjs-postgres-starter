import React from "react";
import { useMutation } from "@apollo/react-hooks";
import { Formik } from "formik";
import * as Yup from "yup";
import { func } from "prop-types";

import LoginFormCmp from "./LoginForm.component";
import { LOGIN_MUTATION } from "../../mutations";
import { showErrorMessage, showSuccessMessage } from "@services/toastify";
import { FETCH_CURRENT_USER_QUERY } from "@graphql/queries/user.queries";
import constants from "@config/constants";

const LoginForm = ({ t }) => {
  const [loginMutation] = useMutation(LOGIN_MUTATION);

  const handleSubmit = async (values, { resetForm, setSubmitting }) => {
    try {
      await loginMutation({
        variables: { login: values.login, password: values.password },
        refetchQueries: [{ query: FETCH_CURRENT_USER_QUERY }],
        awaitRefetchQueries: true
      });
      resetForm();
      showSuccessMessage(t("loginForm.success"));
    } catch (error) {
      if (error.graphQLErrors && error.graphQLErrors[0]) {
        showErrorMessage(error.graphQLErrors[0].message);
      } else {
        showErrorMessage(t("loginForm.errors.default"));
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
      render={formikProps => <LoginFormCmp {...formikProps} t={t} />}
      onSubmit={handleSubmit}
      validationSchema={Yup.object().shape({
        login: Yup.string()
          .trim()
          .min(
            constants.VALIDATION.minimumLengthOfLogin,
            t("registerForm.errors.validation.minimumLengthOfLogin", {
              characters: constants.VALIDATION.minimumLengthOfLogin
            })
          )
          .max(
            constants.VALIDATION.maximumLengthOfLogin,
            t("registerForm.errors.validation.maximumLengthOfLogin", {
              characters: constants.VALIDATION.maximumLengthOfLogin
            })
          )
          .required(t("registerForm.errors.validation.mustProvideLogin")),
        password: Yup.string()
          .trim()
          .required(t("registerForm.errors.validation.mustProvidePassword"))
          .min(
            constants.VALIDATION.minimumLengthOfPassword,
            t("registerForm.errors.validation.minimumLengthOfPassword", {
              characters: constants.VALIDATION.minimumLengthOfPassword
            })
          )
          .max(
            constants.VALIDATION.maximumLengthOfPassword,
            t("registerForm.errors.validation.maximumLengthOfPassword", {
              characters: constants.VALIDATION.maximumLengthOfPassword
            })
          )
          .matches(
            constants.REGEXES.containsUppercase,
            t("registerForm.errors.validation.passwordMustContainOneUppercase")
          )
          .matches(
            constants.REGEXES.containsLowercase,
            t("registerForm.errors.validation.passwordMustContainOneLowercase")
          )
          .matches(
            constants.REGEXES.containsDigit,
            t("registerForm.errors.validation.passwordMustContainOneDigit")
          )
      })}
    />
  );
};

LoginForm.propTypes = {
  t: func.isRequired
};

export default LoginForm;
