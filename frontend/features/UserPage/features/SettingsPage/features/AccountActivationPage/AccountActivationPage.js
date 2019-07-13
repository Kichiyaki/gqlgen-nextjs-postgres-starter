import React from "react";
import Typography from "@material-ui/core/Typography";
import Divider from "@material-ui/core/Divider";
import PageLayout from "../../components/PageLayout/PageLayout";
import withCurrentUser from "@hocs/withCurrentUser";
import restrictionWrapper from "@hocs/restrictionWrapper";
import useI18N from "@lib/i18n/useI18N";
import ActivateAccountForm from "./components/ActivateAccountForm/ActivateAccountForm.container";
import GenerateNewActivationTokenForm from "./components/GenerateNewActivationTokenForm/GenerateNewActivationTokenForm";

const AccountActivationPage = () => {
  const translations = useI18N();
  const {
    USER_PAGE: {
      SETTINGS_PAGE: { ACCOUNT_ACTIVATION_PAGE }
    }
  } = translations;

  return (
    <PageLayout title={ACCOUNT_ACTIVATION_PAGE.title}>
      <Typography variant="h6" component="h3">
        {ACCOUNT_ACTIVATION_PAGE.activateAccount}
      </Typography>
      <div>
        <ActivateAccountForm translations={translations} />
      </div>
      <Divider />
      <div>
        <GenerateNewActivationTokenForm translations={translations} />
      </div>
    </PageLayout>
  );
};

export default withCurrentUser(
  restrictionWrapper({ needAuth: true, needDeactivatedAccount: false })(
    AccountActivationPage
  )
);
