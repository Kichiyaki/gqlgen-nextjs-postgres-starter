import React, { Fragment } from 'react';
import { bool, node, string, object } from 'prop-types';
import classnames from 'classnames';

import { makeStyles } from '@material-ui/core/styles';
import { CssBaseline } from '@material-ui/core';

import Navbar from './Navbar/Navbar';
import Footer from './Footer/Footer';

const useStyles = makeStyles(theme => ({
  '@global': {
    body: {
      position: 'relative'
    }
  },
  container: {
    minHeight: `calc(100vh - 128px)`,
    padding: theme.spacing(3, 0),
    [theme.breakpoints.down('xs')]: {
      minHeight: `calc(100vh - 112px)`
    },
    position: 'relative'
  }
}));

export default function AppLayout({
  children,
  navProps,
  className,
  showFooter
}) {
  const classes = useStyles();
  return (
    <Fragment>
      <Navbar {...navProps} />
      <main className={classnames(classes.container, className)}>
        {children}
      </main>
      {showFooter && <Footer />}
      <CssBaseline />
    </Fragment>
  );
}

AppLayout.defaultProps = {
  navProps: {},
  showFooter: true
};

AppLayout.propTypes = {
  showFooter: bool.isRequired,
  navProps: object.isRequired,
  children: node,
  className: string
};
