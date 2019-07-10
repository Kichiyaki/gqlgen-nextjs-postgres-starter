import React, { Fragment } from "react";
import { node, object, bool } from "prop-types";
import Grid from "@material-ui/core/Grid";
import CssBaseline from "@material-ui/core/CssBaseline";
import AppHeader from "./AppHeader/AppHeader";
import AppFooter from "./AppFooter/AppFooter";

const AppLayoutCmp = ({ children, gridProps, headerProps, showFooter }) => {
  return (
    <Fragment>
      <CssBaseline />
      <AppHeader {...headerProps} />
      <Grid container component="main" {...gridProps}>
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
