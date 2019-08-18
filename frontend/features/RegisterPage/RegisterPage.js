import React, { useState } from "react";
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
import ResetPasswordModal from "./components/ResetPasswordModal/ResetPasswordModal";
import withCurrentUser from "@hocs/withCurrentUser";
import restrictionWrapper from "@hocs/restrictionWrapper";
import constants from "@config/constants";
import pageConstants from "./constants";
import { useTranslation } from "@lib/i18n/i18n";

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
  const [isOpen, setIsOpen] = useState(false);

  const handleOpen = () => {
    setIsOpen(true);
  };

  const handleClose = () => {
    setIsOpen(false);
  };

  const { t } = useTranslation(pageConstants.NAMESPACE);

  return (
    <Container component="main" maxWidth="xs">
      <CssBaseline />
      <div className={classes.paper}>
        <Avatar className={classes.avatar}>
          <LockOutlinedIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          {route === constants.ROUTES.register ? t("signup") : t("signin")}
        </Typography>
        {route === constants.ROUTES.register ? (
          <RegisterForm t={t} />
        ) : (
          <LoginForm t={t} />
        )}
        <Grid container>
          <Grid item xs>
            <Link variant="body2" onClick={handleOpen}>
              {t("forgottenPassword")}
            </Link>
          </Grid>
          <Grid item>
            {route === constants.ROUTES.register ? (
              <NextLink href={constants.ROUTES.login}>
                <Link variant="body2">{t("alreadyHaveAnAccount")}</Link>
              </NextLink>
            ) : (
              <NextLink href={constants.ROUTES.register}>
                <Link variant="body2">{t("dontHaveAnAccount")}</Link>
              </NextLink>
            )}
          </Grid>
        </Grid>
      </div>
      {isOpen && (
        <ResetPasswordModal open={isOpen} t={t} handleClose={handleClose} />
      )}
    </Container>
  );
};

RegisterPage.getInitialProps = () => {
  return {
    namespacesRequired: [pageConstants.NAMESPACE, constants.NAMESPACES.common]
  };
};

export default withCurrentUser(
  restrictionWrapper({ mustBeLoggedOut: true })(RegisterPage)
);
