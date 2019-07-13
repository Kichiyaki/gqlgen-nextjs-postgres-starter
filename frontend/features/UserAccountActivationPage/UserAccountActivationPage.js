import React, { useRef } from "react";
import { isNil } from "lodash";
import Router, { useRouter } from "next/router";
import { Query } from "react-apollo";
import ClipLoader from "react-spinners/ClipLoader";
import isUUID from "validator/lib/isUUID";
import { makeStyles } from "@material-ui/core/styles";

import AppLayout from "@common/AppLayout/AppLayout";
import constants from "@config/constants";
import { showErrorMessage, showSuccessMessage } from "@services/toastify";
import useI18N from "@lib/i18n/useI18N";
import { ACTIVATE_USER_ACCOUNT_QUERY } from "./queries";

const useStyles = makeStyles(() => ({
  container: {
    minHeight: "50vh",
    display: "flex",
    alignItems: "center",
    justifyContent: "center"
  }
}));

const UserAccountActivationPage = () => {
  const called = useRef(false);
  const classes = useStyles();
  const { query, push } = useRouter();
  const translations = useI18N();

  const handleCompleted = ({ activateUserAccount }) => {
    if (!isNil(activateUserAccount)) {
      showSuccessMessage(
        translations.USER_ACCOUNT_ACTIVATION_PAGE.success(
          activateUserAccount.login
        )
      );
    } else {
      showErrorMessage(
        translations.USER_ACCOUNT_ACTIVATION_PAGE.errors.default
      );
    }
    push(constants.ROUTES.root);
  };

  const handleError = ({ graphQLErrors }) => {
    if (called.current) return;
    called.current = true;

    if (graphQLErrors) {
      showErrorMessage(graphQLErrors[0].message);
    } else {
      showErrorMessage(
        translations.USER_ACCOUNT_ACTIVATION_PAGE.errors.default
      );
    }
    push(constants.ROUTES.root);
  };

  return (
    <Query
      query={ACTIVATE_USER_ACCOUNT_QUERY}
      variables={{ id: parseInt(query.id), token: query.token }}
      onCompleted={handleCompleted}
      onError={handleError}
      ssr={false}
      fetchPolicy="network-only"
    >
      {() => {
        return (
          <AppLayout gridProps={{ classes }}>
            <ClipLoader size={250} />
          </AppLayout>
        );
      }}
    </Query>
  );
};

UserAccountActivationPage.getInitialProps = ({ query, res }) => {
  if (
    Object.keys(query).length === 0 ||
    isNaN(parseInt(query.id)) ||
    !isUUID(query.token)
  ) {
    if (res) {
      res.writeHead(302, {
        Location: "/"
      });
      res.end();
      return {};
    } else {
      Router.push("/");
    }
  }

  return {};
};

export default UserAccountActivationPage;
