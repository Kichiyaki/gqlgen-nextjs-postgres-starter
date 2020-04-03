import React, { useState } from 'react';
import { oneOf, bool } from 'prop-types';
import { useQuery, useMutation } from '@apollo/react-hooks';
import { useRouter } from 'next/router';
import { useTranslation } from '@libs/i18n';
import { COMMON } from '@config/namespaces';
import { NAME } from '@config/application';
import { ME } from '@graphql/queries/auth.queries';
import { SIGN_OUT_MUTATION } from './constants';
import {
  MAIN_PAGE,
  SIGN_IN_PAGE,
  SIGN_UP_PAGE,
  USER_PAGE
} from '@config/routes';

import { makeStyles } from '@material-ui/core/styles';
import {
  AppBar,
  Toolbar,
  Typography,
  Container,
  IconButton,
  Menu,
  Hidden,
  Button,
  MenuItem
} from '@material-ui/core';
import { Menu as MenuIcon } from '@material-ui/icons';
import Link from '@common/Link/Link';

const useStyles = makeStyles(theme => ({
  titleWrapper: {
    flexGrow: 1,
    display: 'flex',
    alignItems: 'center'
  },
  whiteBgOnHover: {
    '&:hover': {
      backgroundColor: 'rgba(255, 255, 255, 0.1)'
    }
  },
  navbar: props => {
    const defaults = {
      '&  a': {
        borderRadius: '4px',
        transition: 'all .2s',
        cursor: 'pointer',
        padding: theme.spacing(0.75, 1)
      }
    };
    return props.transparent
      ? {
          ...defaults,
          backgroundColor: 'transparent !important',
          boxShadow: 'none'
        }
      : defaults;
  },
  title: {
    marginRight: theme.spacing(4)
  }
}));

const navLinks = [
  {
    name: 'navbar.link.mainPage',
    href: MAIN_PAGE
  }
];

export default function Navbar({ transparent, position }) {
  const classes = useStyles({ transparent });
  const { t } = useTranslation(COMMON);
  const { data } = useQuery(ME, { fetchPolicy: 'cache-only' });
  const [anchorEl, setAnchorEl] = useState(null);
  const [signOut, { loading }] = useMutation(SIGN_OUT_MUTATION, {
    ignoreResults: true,
    awaitRefetchQueries: true,
    refetchQueries: [{ query: ME }]
  });
  const open = Boolean(anchorEl);
  const logged = data && data.me;

  const handleMenu = event => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleSignOut = () => signOut();

  return (
    <AppBar position={position} component="nav" className={classes.navbar}>
      <Container>
        <Toolbar disableGutters>
          <div className={classes.titleWrapper}>
            <Typography className={classes.title} variant="h4" component="h3">
              {NAME}
            </Typography>
            <Hidden smDown>
              {navLinks.map(l => (
                <Link
                  key={l.name}
                  href={l.href}
                  color="inherit"
                  underline="none"
                  className={classes.whiteBgOnHover}
                >
                  {t(l.name)}
                </Link>
              ))}
            </Hidden>
          </div>
          <div>
            <Hidden smDown>
              {logged ? (
                <>
                  <Link
                    href={USER_PAGE.SETTINGS_PAGE.ACCOUNT_PAGE}
                    color="inherit"
                    underline="none"
                    className={classes.whiteBgOnHover}
                  >
                    {t('navbar.link.settings')}
                  </Link>
                  <Button onClick={handleSignOut} disabled={loading}>
                    {t('navbar.signout')}
                  </Button>
                </>
              ) : (
                <>
                  <Link
                    href={SIGN_IN_PAGE}
                    color="inherit"
                    underline="none"
                    className={classes.whiteBgOnHover}
                  >
                    {t('navbar.link.signin')}
                  </Link>
                  <Link
                    href={SIGN_UP_PAGE}
                    color="inherit"
                    underline="none"
                    className={classes.whiteBgOnHover}
                  >
                    {t('navbar.link.signup')}
                  </Link>
                </>
              )}
            </Hidden>
            <Hidden mdUp>
              <IconButton
                edge="start"
                color="inherit"
                aria-label="menu"
                onClick={handleMenu}
              >
                <MenuIcon />
              </IconButton>
              <Menu
                anchorEl={anchorEl}
                anchorOrigin={{
                  vertical: 'top',
                  horizontal: 'right'
                }}
                keepMounted
                transformOrigin={{
                  vertical: 'top',
                  horizontal: 'right'
                }}
                open={open}
                onClose={handleClose}
                disableScrollLock
              >
                {navLinks.map(l => (
                  <MenuItem key={l.name}>
                    <Link href={l.href} color="inherit" underline="none">
                      {t(l.name)}
                    </Link>
                  </MenuItem>
                ))}
                {logged && (
                  <MenuItem>
                    <Link
                      href={USER_PAGE.SETTINGS_PAGE.ACCOUNT_PAGE}
                      color="inherit"
                      underline="none"
                    >
                      {t('navbar.link.settings')}
                    </Link>
                  </MenuItem>
                )}
                {logged && (
                  <MenuItem onClick={handleSignOut}>
                    {t('navbar.signout')}
                  </MenuItem>
                )}
                {!logged && (
                  <MenuItem>
                    <Link href={SIGN_IN_PAGE} color="inherit" underline="none">
                      {t('navbar.link.signin')}
                    </Link>
                  </MenuItem>
                )}
                {!logged && (
                  <MenuItem>
                    <Link href={SIGN_UP_PAGE} color="inherit" underline="none">
                      {t('navbar.link.signup')}
                    </Link>
                  </MenuItem>
                )}
              </Menu>
            </Hidden>
          </div>
        </Toolbar>
      </Container>
    </AppBar>
  );
}

Navbar.defaultProps = {
  transparent: false,
  position: 'static'
};

Navbar.propTypes = {
  position: oneOf(['relative', 'static', 'fixed', 'sticky', 'absolute'])
    .isRequired,
  transparent: bool.isRequired
};
