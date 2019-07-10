import React, { useContext } from "react";
import { makeStyles } from "@material-ui/core/styles";

import Link from "./components/Link/Link";
import I18NContext from "@lib/i18n/context";
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

const AppFooter = () => {
  const classes = useStyles();
  const {
    APPLICATION: { footer }
  } = useContext(I18NContext);

  return (
    <footer className={classes.footer}>
      <div className={classes.navContainer}>
        <Link href={constants.ROUTES.root}>{footer.links.mainPage}</Link>
        <Link href={constants.ROUTES.register}>
          {footer.links.registration}
        </Link>
        <Link href={constants.ROUTES.login}>{footer.links.login}</Link>
        <Link href={constants.ROUTES.root}>{footer.links.rules}</Link>
        <Link href={constants.ROUTES.root}>{footer.links.aboutAuthor}</Link>
      </div>
      <p className={classes.copyright}>
        &copy; {footer.copyright(new Date().getFullYear())}
      </p>
    </footer>
  );
};

export default AppFooter;
