import React from "react";
import { useMutation, useApolloClient } from "@apollo/react-hooks";
import { Formik } from "formik";
import * as Yup from "yup";
import { object } from "prop-types";
import isUUID from "validator/lib/isUUID";

import ActivateAccountFormCmp from "./ActivateAccountForm.component";
import { ACTIVATE_USER_ACCOUNT_MUTATION } from "../../mutations";
import { showErrorMessage, showSuccessMessage } from "@services/toastify";
import { FETCH_CURRENT_USER_QUERY } from "@graphql/queries/user.queries";

const ActivateAccountForm = ({ translations }) => {
  const [activateUserAccountMutation] = useMutation(
    ACTIVATE_USER_ACCOUNT_MUTATION
  );
  const client = useApolloClient();
  const {
    USER_PAGE: {
      SETTINGS_PAGE: { ACCOUNT_ACTIVATION_PAGE }
    }
  } = translations;

  const handleSubmit = async ({ token }, { resetForm, setSubmitting }) => {
    try {
      const { fetchCurrentUser } = client.readQuery({
        query: FETCH_CURRENT_USER_QUERY
      });
      await activateUserAccountMutation({
        variables: { token, id: fetchCurrentUser.id },
        update: (cache, { data: { activateUserAccount } }) => {
          cache.writeQuery({
            query: FETCH_CURRENT_USER_QUERY,
            data: {
              fetchCurrentUser: {
                ...fetchCurrentUser,
                ...activateUserAccount
              }
            }
          });
        }
      });
      resetForm();
      showSuccessMessage(ACCOUNT_ACTIVATION_PAGE.activateAccountForm.success);
    } catch (error) {
      if (error.graphQLErrors && error.graphQLErrors[0]) {
        showErrorMessage(error.graphQLErrors[0].message);
      } else {
        showErrorMessage(
          ACCOUNT_ACTIVATION_PAGE.activateAccountForm.errors.default
        );
      }
    }
    setSubmitting(false);
  };
  return (
    <Formik
      initialValues={{
        token: ""
      }}
      render={formikProps => (
        <ActivateAccountFormCmp {...formikProps} translations={translations} />
      )}
      onSubmit={handleSubmit}
      validationSchema={Yup.object().shape({
        token: Yup.string()
          .trim()
          .required(
            ACCOUNT_ACTIVATION_PAGE.activateAccountForm.errors.validation
              .mustProvideToken
          )
          .test(
            "is-uuid",
            ACCOUNT_ACTIVATION_PAGE.activateAccountForm.errors.validation
              .tokenIsInvalid,
            val => !val || isUUID(val)
          )
      })}
    />
  );
};

ActivateAccountForm.propTypes = {
  translations: object.isRequired
};

export default ActivateAccountForm;
