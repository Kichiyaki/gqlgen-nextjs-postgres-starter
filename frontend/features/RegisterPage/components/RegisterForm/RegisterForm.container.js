import React from "react";
import { object } from "prop-types";
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

const RegisterForm = ({ translations }) => {
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
      showSuccessMessage(translations.REGISTER_PAGE.registerForm.success);
    } catch (error) {
      if (error.graphQLErrors && error.graphQLErrors[0]) {
        showErrorMessage(error.graphQLErrors[0].message);
      } else {
        showErrorMessage(
          translations.REGISTER_PAGE.registerForm.errors.default
        );
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
      render={formikProps => (
        <RegisterFormCmp {...formikProps} translations={translations} />
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
        email: Yup.string()
          .trim()
          .email(
            translations.REGISTER_PAGE.registerForm.errors.validation
              .invalidEmail
          )
          .required(
            translations.REGISTER_PAGE.registerForm.errors.validation
              .mustProvideEmail
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
          ),
        passwordConfirmation: Yup.string()
          .oneOf(
            [Yup.ref("password"), null],
            translations.REGISTER_PAGE.registerForm.errors.validation
              .passwordsAreNotTheSame
          )
          .required(
            translations.REGISTER_PAGE.registerForm.errors.validation
              .mustProvidePassword
          )
      })}
    />
  );
};

RegisterForm.propTypes = {
  translations: object.isRequired
};

export default RegisterForm;
