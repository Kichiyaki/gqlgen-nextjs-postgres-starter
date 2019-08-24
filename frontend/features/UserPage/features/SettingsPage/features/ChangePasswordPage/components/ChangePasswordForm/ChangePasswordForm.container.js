import React from "react";
import { useMutation } from "react-apollo";
import { Formik } from "formik";
import * as Yup from "yup";
import { func } from "prop-types";

import constants from "@config/constants";
import ChangePasswordFormCmp from "./ChangePasswordForm.component";
import { CHANGE_PASSWORD_MUTATION } from "../../mutations";
import { FETCH_CURRENT_USER_QUERY } from "@graphql/queries/user.queries";
import { showErrorMessage, showSuccessMessage } from "@services/toastify";
import pageConstants from "../../constants";

const passwordValidation = (t, type = "") =>
  Yup.string()
    .trim()
    .required(
      t(
        constants.NAMESPACES.registerPage +
          ":registerForm.errors.validation.mustProvidePassword"
      )
    )
    .min(
      constants.VALIDATION.minimumLengthOfPassword,
      t(
        constants.NAMESPACES.registerPage +
          ":registerForm.errors.validation.minimumLengthOfPassword",
        {
          count: constants.VALIDATION.minimumLengthOfPassword
        }
      )
    )
    .max(
      constants.VALIDATION.maximumLengthOfPassword,
      t(
        constants.NAMESPACES.registerPage +
          ":registerForm.errors.validation.maximumLengthOfPassword",
        {
          count: constants.VALIDATION.maximumLengthOfPassword
        }
      )
    )
    .matches(
      constants.REGEXES.containsUppercase,
      t(
        constants.NAMESPACES.registerPage +
          ":registerForm.errors.validation.passwordMustContainOneUppercase"
      )
    )
    .matches(
      constants.REGEXES.containsLowercase,
      t(
        constants.NAMESPACES.registerPage +
          ":registerForm.errors.validation.passwordMustContainOneLowercase"
      )
    )
    .notOneOf(
      type === pageConstants.NEW_PASSWORD
        ? [Yup.ref(pageConstants.CURRENT_PASSWORD)]
        : [],
      t("changePasswordForm.errors.validation.passwordsCannotBeTheSame")
    )
    .matches(
      constants.REGEXES.containsDigit,
      t(
        constants.NAMESPACES.registerPage +
          ":registerForm.errors.validation.passwordMustContainOneDigit"
      )
    );

const ChangePasswordForm = ({ t }) => {
  const [changePasswordMutation] = useMutation(CHANGE_PASSWORD_MUTATION);

  const handleSubmit = async (payload, { resetForm, setSubmitting }) => {
    try {
      await changePasswordMutation({
        variables: payload,
        refetchQueries: [{ query: FETCH_CURRENT_USER_QUERY }],
        awaitRefetchQueries: true
      });
      resetForm();
      showSuccessMessage(t("changePasswordForm.success"));
    } catch (error) {
      if (error.graphQLErrors && error.graphQLErrors[0]) {
        showErrorMessage(error.graphQLErrors[0].message);
      } else {
        showErrorMessage(t("changePasswordForm.errors.default"));
      }
    }
    setSubmitting(false);
  };
  return (
    <Formik
      initialValues={{
        currentPassword: "",
        newPassword: ""
      }}
      render={formikProps => <ChangePasswordFormCmp {...formikProps} t={t} />}
      onSubmit={handleSubmit}
      validationSchema={Yup.object().shape({
        currentPassword: passwordValidation(t),
        newPassword: passwordValidation(t, pageConstants.NEW_PASSWORD)
      })}
    />
  );
};

ChangePasswordForm.propTypes = {
  t: func.isRequired
};

export default ChangePasswordForm;
