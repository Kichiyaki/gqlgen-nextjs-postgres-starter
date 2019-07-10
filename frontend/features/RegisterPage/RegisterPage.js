import React, { useContext } from "react";
import NextLink from "next/link";
import { useRouter } from "next/router";
import Avatar from "@material-ui/core/Avatar";
import CssBaseline from "@material-ui/core/CssBaseline";
import Link from "@material-ui/core/Link";
import Grid from "@material-ui/core/Grid";
import LockOutlinedIcon from "@material-ui/icons/LockOutlined";
import Typography from "@material-ui/core/Typography";
import { makeStyles } from "@material-ui/core/styles";
import Container from "@material-ui/core/Container";
import RegisterForm from "./components/RegisterForm/RegisterForm.container";
import LoginForm from "./components/LoginForm/LoginForm.container";
import withCurrentUser from "@hocs/withCurrentUser";
import restrictionWrapper from "@hocs/restrictionWrapper";
import constants from "@config/constants";
import I18NContext from "@lib/i18n/context";

const useStyles = makeStyles(theme => ({
  "@global": {
    body: {
      backgroundColor: theme.palette.common.white
    }
  },
  paper: {
    marginTop: theme.spacing(8),
    display: "flex",
    flexDirection: "column",
    alignItems: "center"
  },
  avatar: {
    margin: theme.spacing(1),
    backgroundColor: theme.palette.secondary.main
  }
}));

const RegisterPage = () => {
  const classes = useStyles();
  const { route } = useRouter();
  const translations = useContext(I18NContext);

  return (
    <Container component="main" maxWidth="xs">
      <CssBaseline />
      <div className={classes.paper}>
        <Avatar className={classes.avatar}>
          <LockOutlinedIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          {route === constants.ROUTES.register
            ? translations.REGISTER_PAGE.signup
            : translations.REGISTER_PAGE.signin}
        </Typography>
        {route === constants.ROUTES.register ? (
          <RegisterForm translations={translations} />
        ) : (
          <LoginForm translations={translations} />
        )}
        <Grid container>
          <Grid item xs>
            <NextLink href="#">
              <Link variant="body2">
                {translations.REGISTER_PAGE.forgottenPassword}
              </Link>
            </NextLink>
          </Grid>
          <Grid item>
            {route === constants.ROUTES.register ? (
              <NextLink href={constants.ROUTES.login}>
                <Link variant="body2">
                  {translations.REGISTER_PAGE.alreadyHaveAnAccount}
                </Link>
              </NextLink>
            ) : (
              <NextLink href={constants.ROUTES.register}>
                <Link variant="body2">
                  {translations.REGISTER_PAGE.dontHaveAnAccount}
                </Link>
              </NextLink>
            )}
          </Grid>
        </Grid>
      </div>
    </Container>
  );
};

export default withCurrentUser(
  restrictionWrapper({ mustBeLoggedOut: true })(RegisterPage)
);
