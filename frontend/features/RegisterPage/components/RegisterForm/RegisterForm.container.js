import React from "react";
import { func } from "prop-types";
import { useMutation } from "@apollo/react-hooks";
import { Formik } from "formik";
import * as Yup from "yup";
import { omit } from "lodash";

import constants from "@config/constants";
import RegisterFormCmp from "./RegisterForm.component";
import { SIGNUP_MUTATION } from "../../mutations";
import { showErrorMessage, showSuccessMessage } from "@services/toastify";
import { FETCH_CURRENT_USER_QUERY } from "@graphql/queries/user.queries";
import registerPageConstants from "../../constants";

const RegisterForm = ({ t }) => {
  const [signupMutation] = useMutation(SIGNUP_MUTATION);

  const handleSubmit = async (values, { resetForm, setSubmitting }) => {
    try {
      await signupMutation({
        variables: {
          user: omit(values, [registerPageConstants.PASSWORD_CONFIRMATION])
        },
        refetchQueries: [{ query: FETCH_CURRENT_USER_QUERY }],
        awaitRefetchQueries: true
      });
      resetForm();
      showSuccessMessage(t("registerForm.success"));
    } catch (error) {
      if (error.graphQLErrors && error.graphQLErrors[0]) {
        showErrorMessage(error.graphQLErrors[0].message);
      } else {
        showErrorMessage(t("registerForm.errors.default"));
      }
    }
    setSubmitting(false);
  };

  return (
    <Formik
      initialValues={{
        login: "",
        email: "",
        password: "",
        passwordConfirmation: ""
      }}
      render={formikProps => <RegisterFormCmp {...formikProps} t={t} />}
      onSubmit={handleSubmit}
      validationSchema={Yup.object().shape({
        login: Yup.string()
          .trim()
          .min(
            constants.VALIDATION.minimumLengthOfLogin,
            t("registerForm.errors.validation.minimumLengthOfLogin", {
              count: constants.VALIDATION.minimumLengthOfLogin
            })
          )
          .max(
            constants.VALIDATION.maximumLengthOfLogin,
            t("registerForm.errors.validation.maximumLengthOfLogin", {
              count: constants.VALIDATION.maximumLengthOfLogin
            })
          )
          .required(t("registerForm.errors.validation.mustProvideLogin")),
        email: Yup.string()
          .trim()
          .email(t("registerForm.errors.validation.invalidEmail"))
          .required(t("registerForm.errors.validation.mustProvideEmail")),
        password: Yup.string()
          .trim()
          .required(t("registerForm.errors.validation.mustProvidePassword"))
          .min(
            constants.VALIDATION.minimumLengthOfPassword,
            t("registerForm.errors.validation.minimumLengthOfPassword", {
              count: constants.VALIDATION.minimumLengthOfPassword
            })
          )
          .max(
            constants.VALIDATION.maximumLengthOfPassword,
            t("registerForm.errors.validation.maximumLengthOfPassword", {
              count: constants.VALIDATION.maximumLengthOfPassword
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
          ),
        passwordConfirmation: Yup.string()
          .oneOf(
            [Yup.ref("password"), null],
            t("registerForm.errors.validation.passwordsAreNotTheSame")
          )
          .required(t("registerForm.errors.validation.mustProvidePassword"))
      })}
    />
  );
};

RegisterForm.propTypes = {
  t: func.isRequired
};

export default RegisterForm;
