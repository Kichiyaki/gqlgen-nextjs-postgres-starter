import React from "react";
import Typography from "@material-ui/core/Typography";
import Divider from "@material-ui/core/Divider";
import PageLayout from "../../components/PageLayout/PageLayout";
import withCurrentUser from "@hocs/withCurrentUser";
import restrictionWrapper from "@hocs/restrictionWrapper";
import { useTranslation } from "@lib/i18n/i18n";
import ActivateAccountForm from "./components/ActivateAccountForm/ActivateAccountForm.container";
import GenerateNewActivationTokenForm from "./components/GenerateNewActivationTokenForm/GenerateNewActivationTokenForm";
import constants from "@config/constants";
import pageConstants from "./constants";

const AccountActivationPage = () => {
  const { t } = useTranslation(pageConstants.NAMESPACE);

  return (
    <PageLayout title={t("title")}>
      <Typography variant="h6" component="h3">
        {t("activateAccount")}
      </Typography>
      <div>
        <ActivateAccountForm t={t} />
      </div>
      <Divider />
      <div>
        <GenerateNewActivationTokenForm t={t} />
      </div>
    </PageLayout>
  );
};

AccountActivationPage.getInitialProps = () => {
  return {
    namespacesRequired: [
      constants.NAMESPACES.common,
      constants.NAMESPACES.userPage.settingsPage.navigation,
      pageConstants.NAMESPACE
    ]
  };
};

export default withCurrentUser(
  restrictionWrapper({ needAuth: true, needDeactivatedAccount: true })(
    AccountActivationPage
  )
);
