import React from "react";
import { useMutation } from "react-apollo";
import { isNil } from "lodash";

import { makeStyles } from "@material-ui/core/styles";
import AppBar from "@material-ui/core/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from "@material-ui/core/Typography";
import Button from "@material-ui/core/Button";

import Link from "@common/Link/Link";
import { showErrorMessage, showSuccessMessage } from "@services/toastify";
import useCurrentUser from "@hooks/useCurrentUser";
import { FETCH_CURRENT_USER_QUERY } from "@graphql/queries/user.queries";
import { useTranslation } from "@lib/i18n/i18n";
import { LOGOUT_USER_MUTATION } from "./mutations";
import globalConstants from "@config/constants";
import constants from "./constants";

const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1
  },
  title: {
    flexGrow: 1
  },
  buttons__container: {
    "& > *:not(:last-child)": {
      marginRight: theme.spacing(1)
    }
  }
}));

const AppHeader = () => {
  const classes = useStyles();

  const {
    data: { fetchCurrentUser }
  } = useCurrentUser();

  const { t } = useTranslation(globalConstants.NAMESPACES.common);

  const [logoutUserMutation, { loading }] = useMutation(LOGOUT_USER_MUTATION);

  const logoutUser = async () => {
    try {
      await logoutUserMutation({
        refetchQueries: [{ query: FETCH_CURRENT_USER_QUERY }],
        awaitRefetchQueries: true
      });
      showSuccessMessage(t("HEADER.logout.success"));
    } catch (error) {
      showErrorMessage(t("HEADER.logout.error"));
    }
  };

  return (
    <div className={classes.root}>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" className={classes.title}>
            {t("APPLICATION.name")}
          </Typography>
          {!isNil(fetchCurrentUser) && (
            <div className={classes.buttons__container}>
              {!fetchCurrentUser.activated && (
                <Link
                  href={
                    globalConstants.ROUTES.userPage.settingsPage
                      .accountActivation
                  }
                  linkProps={{ color: "secondary" }}
                >
                  {t("HEADER.links.activateAccount")}
                </Link>
              )}
              <Button
                onClick={logoutUser}
                data-testid={constants.LOGOUT_BUTTON}
                disabled={loading}
                color="inherit"
              >
                {t("HEADER.buttons.logout")}
              </Button>
            </div>
          )}
        </Toolbar>
      </AppBar>
    </div>
  );
};

export default AppHeader;
