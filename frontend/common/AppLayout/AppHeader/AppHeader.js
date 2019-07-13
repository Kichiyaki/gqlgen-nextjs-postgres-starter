import React, { useContext } from "react";
import { useMutation } from "@apollo/react-hooks";
import { isNil } from "lodash";
import { makeStyles } from "@material-ui/core/styles";
import AppBar from "@material-ui/core/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from "@material-ui/core/Typography";
import Button from "@material-ui/core/Button";
import I18NContext from "@lib/i18n/context";
import { showErrorMessage, showSuccessMessage } from "@services/toastify";
import useCurrentUser from "@hooks/useCurrentUser";
import { LOGOUT_USER_MUTATION } from "./mutations";
import constants from "./constants";

const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1
  },
  title: {
    flexGrow: 1
  }
}));

const AppHeader = () => {
  const classes = useStyles();
  const { APPLICATION } = useContext(I18NContext);

  const {
    data: { fetchCurrentUser }
  } = useCurrentUser();

  const [logoutUserMutation, { loading }] = useMutation(LOGOUT_USER_MUTATION);

  const logoutUser = async () => {
    try {
      await logoutUserMutation({
        refetchQueries: [{ query: FETCH_CURRENT_USER_QUERY }],
        awaitRefetchQueries: true
      });
      showSuccessMessage(APPLICATION.header.logout.success);
    } catch (error) {
      showErrorMessage(APPLICATION.header.logout.error);
    }
  };

  return (
    <div className={classes.root}>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" className={classes.title}>
            {APPLICATION.name}
          </Typography>
          <div>
            {!isNil(fetchCurrentUser) && (
              <Button
                onClick={logoutUser}
                data-testid={constants.LOGOUT_BUTTON}
                disabled={loading}
                color="inherit"
              >
                {APPLICATION.header.buttons.logout}
              </Button>
            )}
          </div>
        </Toolbar>
      </AppBar>
    </div>
  );
};

export default AppHeader;
