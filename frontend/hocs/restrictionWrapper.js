import React, { useEffect } from "react";
import { useQuery } from "@apollo/react-hooks";
import Router, { useRouter } from "next/router";
import { isNil } from "lodash";
import constants from "../config/constants";
import { FETCH_CURRENT_USER_QUERY } from "../graphql/queries/user.queries";

const buildRedirectObj = (redirect = false, location = "") => ({
  redirect,
  location
});

const restrictionWrapper = ({
  needAuth = false,
  needActivatedAccount = false,
  needDeactivatedAccount = false,
  needAdministrativePrivileges = false,
  mustBeLoggedOut = false
} = {}) => WrappedComponent => {
  const checkUserPrivileges = user => {
    const isLoggedIn = !isNil(user);

    if (needAuth && !isLoggedIn) {
      return buildRedirectObj(true, constants.ROUTES.login);
    } else if (mustBeLoggedOut && isLoggedIn) {
      return buildRedirectObj(true, constants.ROUTES.root);
    } else if (isLoggedIn && !user.activated && needActivatedAccount) {
      return buildRedirectObj(true, constants.ROUTES.root);
    } else if (isLoggedIn && user.activated && needDeactivatedAccount) {
      return buildRedirectObj(true, constants.ROUTES.root);
    } else if (
      isLoggedIn &&
      user.role !== constants.ROLES.administrativeRole &&
      needAdministrativePrivileges
    ) {
      return buildRedirectObj(true, constants.ROUTES.root);
    }

    return buildRedirectObj(false);
  };

  const RestrictionWrapper = () => {
    const {
      data: { fetchCurrentUser }
    } = useQuery(FETCH_CURRENT_USER_QUERY, { fetchPolicy: "cache-only" });
    const router = useRouter();

    useEffect(() => {
      const { redirect, location } = checkUserPrivileges(fetchCurrentUser);
      if (redirect) {
        router.push(location);
      }
    }, [fetchCurrentUser]);

    return <WrappedComponent />;
  };

  RestrictionWrapper.getInitialProps = async ctx => {
    const { apolloClient, res } = ctx;
    let pageProps = {};

    const { fetchCurrentUser } = apolloClient.readQuery({
      query: FETCH_CURRENT_USER_QUERY
    });

    if (res) {
      const { redirect, location } = checkUserPrivileges(fetchCurrentUser);
      if (redirect) {
        res.writeHead(302, {
          Location: location
        });
        res.end();
        return { pageProps };
      }
    } else {
      const { redirect, location } = checkUserPrivileges(fetchCurrentUser);
      if (redirect) {
        Router.push(location);
        return { pageProps };
      }
    }

    if (WrappedComponent.getInitialProps) {
      pageProps = await WrappedComponent.getInitialProps(ctx);
    }

    return { pageProps };
  };

  return RestrictionWrapper;
};

export default restrictionWrapper;
