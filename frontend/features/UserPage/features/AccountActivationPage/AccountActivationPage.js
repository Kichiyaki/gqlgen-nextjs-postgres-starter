import React, { useEffect } from "react";
import Router, { useRouter } from "next/router";
import { useApolloClient } from "react-apollo";
import ClipLoader from "react-spinners/ClipLoader";
import isUUID from "validator/lib/isUUID";

import AppLayout from "@common/AppLayout/AppLayout";
import constants from "@config/constants";
import { showErrorMessage, showSuccessMessage } from "@services/toastify";
import { ACTIVATE_USER_ACCOUNT_QUERY } from "./queries";
import { useTranslation } from "@lib/i18n/i18n";
import pageConstants from "./constants";

const AccountActivationPage = () => {
  const { query, push } = useRouter();
  const { t } = useTranslation(pageConstants.NAMESPACE);
  const client = useApolloClient();

  useEffect(() => {
    if (window) {
      fetchData();
    }
  }, [query]);

  const fetchData = async () => {
    try {
      const {
        data: { activateUserAccount }
      } = await client.query({
        query: ACTIVATE_USER_ACCOUNT_QUERY,
        variables: { id: parseInt(query.id), token: query.token },
        fetchPolicy: "no-cache"
      });
      showSuccessMessage(t("success", { login: activateUserAccount.login }));
      push(constants.ROUTES.root);
    } catch ({ graphQLErrors }) {
      if (graphQLErrors) {
        showErrorMessage(graphQLErrors[0].message);
      } else {
        showErrorMessage(t("errors.default"));
      }
      push(constants.ROUTES.root);
    }
  };

  return (
    <AppLayout gridProps={{ justify: "center", align: "center" }}>
      <ClipLoader size={250} />
    </AppLayout>
  );
};

AccountActivationPage.getInitialProps = ({ query, res }) => {
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

export default AccountActivationPage;
