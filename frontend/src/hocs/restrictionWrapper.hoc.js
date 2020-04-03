import React, { useEffect } from 'react';
import URLSearchParams from '@ungap/url-search-params';
import { useQuery } from '@apollo/react-hooks';
import Router, { useRouter } from 'next/router';
import { ME } from '@graphql/queries/auth.queries';
import { COMMON } from '@config/namespaces';
import { SIGN_IN_PAGE, MAIN_PAGE } from '@config/routes';

const restrictionWrapper = ({
  loggedIn = false,
  loggedOut = false,
  activated = false,
  deactivated = false
} = {}) => WrappedCmp => {
  const getMatchRoute = (
    pathname = '',
    query = {},
    logged = false,
    a = false
  ) => {
    if (!logged && loggedIn) {
      return [SIGN_IN_PAGE, SIGN_IN_PAGE, true];
    } else if (logged && loggedOut) {
      if (pathname === SIGN_IN_PAGE && query.asPath && query.pathname) {
        return [query.pathname, query.asPath, false];
      }
      return [MAIN_PAGE, MAIN_PAGE, false];
    } else if (a && deactivated) {
      return [MAIN_PAGE, MAIN_PAGE, false];
    } else if (!a && activated) {
      return [MAIN_PAGE, MAIN_PAGE, false];
    }
    return ['', '', false];
  };

  const RestrictionWrapper = props => {
    const router = useRouter();
    const { data } = useQuery(ME, {
      fetchPolicy: 'cache-only'
    });
    const logged = data && data.me;
    const id = data && data.me && data.me.id ? data.me.id : 0;
    const role = data && data.me && data.me.role ? data.me.role : 0;
    const activated = data && data.me ? data.me.activated : false;

    useEffect(() => {
      const params = new URLSearchParams(window.location.search);
      const [route, asPath, shouldAddReturnURL] = getMatchRoute(
        router.pathname,
        {
          asPath: params.get('asPath'),
          pathname: params.get('pathname')
        },
        logged,
        data && data.me ? data.me.activated : false
      );
      if (route) {
        const query = shouldAddReturnURL
          ? `?asPath=${router.asPath}&pathname=${router.pathname}`
          : '';
        router.push(route, asPath + query);
      }
    }, [logged, id, role, activated]);

    return <WrappedCmp {...props} />;
  };
  RestrictionWrapper.getInitialProps = async ctx => {
    const { res, apolloClient, pathname, query } = ctx;
    let pageProps = { namespacesRequired: [COMMON] };
    const { me: user } = await apolloClient.readQuery({ query: ME });
    const logged = !!user;

    const [route, asPath, shouldAddReturnURL] = getMatchRoute(
      pathname,
      query,
      logged,
      user ? user.activated : false
    );

    if (route) {
      const query = shouldAddReturnURL
        ? `?returnUrl=${asPath}&pathname=${pathname}`
        : '';
      if (res) {
        res.writeHead(302, {
          Location: asPath + query
        });
        res.end();
        return pageProps;
      } else {
        await Router.push(route, asPath + query);
        return pageProps;
      }
    }

    if (WrappedCmp.getInitialProps) {
      pageProps = await WrappedCmp.getInitialProps(ctx);
    }
    return pageProps;
  };

  return RestrictionWrapper;
};

export default restrictionWrapper;
