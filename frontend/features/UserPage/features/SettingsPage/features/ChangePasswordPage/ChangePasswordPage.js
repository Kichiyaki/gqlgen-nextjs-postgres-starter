import React from "react";
import withCurrentUser from "@hocs/withCurrentUser";
import restrictionWrapper from "@hocs/restrictionWrapper";
import { useTranslation } from "@lib/i18n/i18n";
import constants from "@config/constants";
import ChangePasswordForm from "./components/ChangePasswordForm/ChangePasswordForm.container";
import pageConstants from "./constants";
import PageLayout from "../../components/PageLayout/PageLayout";

const ChangePasswordPage = () => {
  const { t } = useTranslation([
    pageConstants.NAMESPACE,
    constants.NAMESPACES.registerPage
  ]);

  return (
    <PageLayout title={t("title")}>
      <ChangePasswordForm t={t} />
    </PageLayout>
  );
};

ChangePasswordPage.getInitialProps = () => {
  return {
    namespacesRequired: [
      constants.NAMESPACES.common,
      constants.NAMESPACES.userPage.settingsPage.navigation,
      constants.NAMESPACES.registerPage,
      pageConstants.NAMESPACE
    ]
  };
};

export default withCurrentUser(
  restrictionWrapper({ needAuth: true })(ChangePasswordPage)
);
