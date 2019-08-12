import React from "react";
import { useMutation, useApolloClient } from "@apollo/react-hooks";
import { Formik } from "formik";
import * as Yup from "yup";
import { func } from "prop-types";
import isUUID from "validator/lib/isUUID";

import ActivateAccountFormCmp from "./ActivateAccountForm.component";
import { ACTIVATE_USER_ACCOUNT_MUTATION } from "../../mutations";
import { showErrorMessage, showSuccessMessage } from "@services/toastify";
import { FETCH_CURRENT_USER_QUERY } from "@graphql/queries/user.queries";

const ActivateAccountForm = ({ t }) => {
  const [activateUserAccountMutation] = useMutation(
    ACTIVATE_USER_ACCOUNT_MUTATION
  );
  const client = useApolloClient();

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
      showSuccessMessage(t("activateAccountForm.success"));
    } catch (error) {
      if (error.graphQLErrors && error.graphQLErrors[0]) {
        showErrorMessage(error.graphQLErrors[0].message);
      } else {
        showErrorMessage(t("activateAccountForm.errors.default"));
      }
    }
    setSubmitting(false);
  };
  return (
    <Formik
      initialValues={{
        token: ""
      }}
      render={formikProps => <ActivateAccountFormCmp {...formikProps} t={t} />}
      onSubmit={handleSubmit}
      validationSchema={Yup.object().shape({
        token: Yup.string()
          .trim()
          .required(t("activateAccountForm.errors.validation.mustProvideToken"))
          .test(
            "is-uuid",
            t("activateAccountForm.errors.validation.tokenIsInvalid"),
            val => !val || isUUID(val)
          )
      })}
    />
  );
};

ActivateAccountForm.propTypes = {
  t: func.isRequired
};

export default ActivateAccountForm;
