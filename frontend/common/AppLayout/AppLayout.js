import React, { Fragment } from "react";
import { node, object, bool } from "prop-types";
import { makeStyles } from "@material-ui/core/styles";
import Grid from "@material-ui/core/Grid";
import CssBaseline from "@material-ui/core/CssBaseline";
import AppHeader from "./AppHeader/AppHeader";
import AppFooter from "./AppFooter/AppFooter";

const useStyles = makeStyles(theme => ({
  main: {
    margin: theme.spacing(3, 0, 3, 0)
  }
}));

const AppLayoutCmp = ({ children, gridProps, headerProps, showFooter }) => {
  const classes = useStyles();
  return (
    <Fragment>
      <CssBaseline />
      <AppHeader {...headerProps} />
      <Grid
        container
        component="main"
        spacing={2}
        {...gridProps}
        classes={
          gridProps.classes
            ? { ...gridProps.classes, container: classes.main }
            : { container: classes.main }
        }
      >
        {children}
      </Grid>
      {showFooter && <AppFooter />}
    </Fragment>
  );
};

AppLayoutCmp.defaultProps = {
  headerProps: {},
  gridProps: {},
  showFooter: true
};

AppLayoutCmp.propTypes = {
  children: node.isRequired,
  gridProps: object,
  headerProps: object,
  showFooter: bool
};

export default AppLayoutCmp;
