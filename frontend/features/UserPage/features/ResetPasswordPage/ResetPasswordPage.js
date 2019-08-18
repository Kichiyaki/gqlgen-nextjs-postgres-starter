import React, { useRef } from "react";
import Router, { useRouter } from "next/router";
import { Query } from "react-apollo";
import ClipLoader from "react-spinners/ClipLoader";
import isUUID from "validator/lib/isUUID";

import AppLayout from "@common/AppLayout/AppLayout";
import constants from "@config/constants";
import pageConstants from "./constants";
import { showErrorMessage, showSuccessMessage } from "@services/toastify";
import { RESET_PASSWORD_QUERY } from "./queries";
import { useTranslation } from "@lib/i18n/i18n";

const ResetPasswordPage = () => {
  const called = useRef(false);
  const { query, push } = useRouter();
  const { t } = useTranslation(pageConstants.NAMESPACE);

  const handleCompleted = () => {
    showSuccessMessage(t("success"));
    push(constants.ROUTES.root);
  };

  const handleError = ({ graphQLErrors }) => {
    if (called.current) return;
    called.current = true;

    if (graphQLErrors) {
      showErrorMessage(graphQLErrors[0].message);
    } else {
      showErrorMessage(t("errors.default"));
    }
    push(constants.ROUTES.root);
  };

  return (
    <Query
      query={RESET_PASSWORD_QUERY}
      variables={{ id: parseInt(query.id), token: query.token }}
      onCompleted={handleCompleted}
      onError={handleError}
      ssr={false}
      fetchPolicy="network-only"
    >
      {() => {
        return (
          <AppLayout gridProps={{ justify: "center", align: "center" }}>
            <ClipLoader size={250} />
          </AppLayout>
        );
      }}
    </Query>
  );
};

ResetPasswordPage.getInitialProps = ({ query, res }) => {
  const props = {
    namespacesRequired: [constants.NAMESPACES.common, pageConstants.NAMESPACE]
  };
  if (
    Object.keys(query).length === 0 ||
    isNaN(parseInt(query.id)) ||
    !isUUID(query.token)
  ) {
    if (res) {
      res.writeHead(302, {
        Location: constants.ROUTES.root
      });
      res.end();
      return props;
    } else {
      Router.push(constants.ROUTES.root);
    }
  }

  return props;
};

export default ResetPasswordPage;
