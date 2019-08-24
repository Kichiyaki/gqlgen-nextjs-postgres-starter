import React from "react";
import { useRouter } from "next/router";
import { makeStyles } from "@material-ui/core/styles";
import Typography from "@material-ui/core/Typography";
import Paper from "@material-ui/core/Paper";
import List from "@material-ui/core/List";
import Divider from "@material-ui/core/Divider";
import CheckCircleIcon from "@material-ui/icons/CheckCircle";
import FingerprintIcon from "@material-ui/icons/Fingerprint";
import useCurrentUser from "@hooks/useCurrentUser";
import constants from "@config/constants";
import NavLink from "./components/NavLink/NavLink";
import { withTranslation } from "@lib/i18n/i18n";

const useStyles = makeStyles(theme => ({
  root: {
    padding: theme.spacing(2)
  }
}));

const Navigation = ({ t }) => {
  const {
    data: { fetchCurrentUser }
  } = useCurrentUser();
  const classes = useStyles();
  const { route } = useRouter();

  return (
    <Paper className={classes.root}>
      <Typography align="center" variant="h5" component="h2">
        {t("settings")}
      </Typography>
      <Divider />
      <List component="nav">
        {fetchCurrentUser && !fetchCurrentUser.activated && (
          <NavLink
            text={t("accountActivation")}
            href={constants.ROUTES.userPage.settingsPage.accountActivation}
            selected={
              constants.ROUTES.userPage.settingsPage.accountActivation === route
            }
          >
            <CheckCircleIcon />
          </NavLink>
        )}
        <NavLink
          text={t("changePassword")}
          href={constants.ROUTES.userPage.settingsPage.changePassword}
          selected={
            constants.ROUTES.userPage.settingsPage.changePassword === route
          }
        >
          <FingerprintIcon />
        </NavLink>
      </List>
    </Paper>
  );
};

export default withTranslation("user-page/settings-page/navigation")(
  Navigation
);
