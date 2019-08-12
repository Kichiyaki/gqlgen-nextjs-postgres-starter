import React from "react";
import { func } from "prop-types";
import { makeStyles } from "@material-ui/core/styles";

import Link from "./components/Link/Link";
import { withTranslation } from "@lib/i18n/i18n";
import constants from "@config/constants";

const useStyles = makeStyles(theme => ({
  footer: {
    width: "100%",
    paddingTop: 30,
    paddingBottom: 30,
    backgroundColor: "#000",
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
    color: "#808080"
  },
  navContainer: {
    padding: "0 30px",
    display: "flex",
    flexDirection: "row",
    flexWrap: "wrap",
    marginBottom: theme.spacing(1)
  },
  copyright: {
    fontSize: 12,
    padding: 0,
    margin: 0
  }
}));

const AppFooter = ({ t }) => {
  const classes = useStyles();

  return (
    <footer className={classes.footer}>
      <div className={classes.navContainer}>
        <Link href={constants.ROUTES.root}>{t("FOOTER.links.mainPage")}</Link>
        <Link href={constants.ROUTES.register}>
          {t("FOOTER.links.registration")}
        </Link>
        <Link href={constants.ROUTES.login}>{t("FOOTER.links.login")}</Link>
        <Link href={constants.ROUTES.root}>{t("FOOTER.links.termsOfUse")}</Link>
        <Link href={constants.ROUTES.root}>
          {t("FOOTER.links.aboutAuthor")}
        </Link>
      </div>
      <p className={classes.copyright}>
        &copy;{" "}
        {t("FOOTER.copyright", {
          year: new Date().getFullYear(),
          fullName: constants.AUTHOR_FULL_NAME
        })}
      </p>
    </footer>
  );
};

AppFooter.propTypes = {
  t: func.isRequired
};

export default withTranslation(constants.NAMESPACES.common)(AppFooter);
