import React, { useRef } from "react";
import { RouterContext } from "next-server/dist/lib/router-context";
import Head from "next/head";
import { getDataFromTree } from "react-apollo";
import { object } from "prop-types";
import initApollo from "./initApollo";
import getDisplayName from "./getDisplayName";

const ssrMode = !process.browser;

export default (client, options = { getDataFromTree: "ssr" }) => App => {
  const Apollo = props => {
    const ref = useRef(
      initApollo(client, {
        initialState: props.apolloState.data
      })
    );
    return <App {...props} apollo={ref.current} />;
  };

  Apollo.displayName = `withApollo(${getDisplayName(App)})`;

  Apollo.propTypes = {
    apolloState: object,
    apollo: object
  };

  Apollo.getInitialProps = async appCtx => {
    const { Component, router, ctx } = appCtx;
    const headers = ctx.req ? ctx.req.headers : {};
    const apollo = initApollo(client, { ctx, headers });
    const apolloState = {};
    const getInitialProps = App.getInitialProps;

    let appProps = { pageProps: {} };

    if (getInitialProps) {
      ctx.apolloClient = apollo;
      appProps = await getInitialProps(appCtx);
    }

    if (ctx.res && (ctx.res.headersSent || ctx.res.finished)) {
      return {};
    }

    if (
      options.getDataFromTree === "always" ||
      (options.getDataFromTree === "ssr" && ssrMode)
    ) {
      try {
        await getDataFromTree(
          <RouterContext.Provider
            value={{
              route: router.route,
              pathname: router.pathname,
              query: router.query,
              asPath: router.asPath
            }}
          >
            <App
              {...appProps}
              Component={Component}
              router={router}
              apolloState={apolloState}
              apollo={apollo}
            />
          </RouterContext.Provider>
        );
      } catch (error) {
        // Prevent Apollo Client GraphQL errors from crashing SSR.
        if (process.env.NODE_ENV !== "production") {
          // tslint:disable-next-line no-console This is a necessary debugging log
          console.error("GraphQL error occurred [getDataFromTree]", error);
        }
      }

      if (ssrMode) {
        // getDataFromTree does not call componentWillUnmount
        // head side effect therefore need to be cleared manually
        Head.rewind();
      }

      apolloState.data = apollo.cache.extract();
    }

    return {
      ...appProps,
      apolloState
    };
  };

  return Apollo;
};
