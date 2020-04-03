import React from 'react';
import { AppBar, Toolbar, Typography, Container } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { NAME } from '@config/application';

const useStyles = makeStyles(() => ({
  title: {
    flexGrow: 1
  }
}));

export default function Footer() {
  const classes = useStyles();
  return (
    <AppBar position="static" component="footer">
      <Toolbar disableGutters>
        <Container>
          <Typography component="p">
            &copy; {new Date().getFullYear()} {NAME}
          </Typography>
        </Container>
      </Toolbar>
    </AppBar>
  );
}
