import React, { useRef } from 'react';
import fetch from 'node-fetch';
import { ApolloClient } from '@apollo/client';
import { InMemoryCache } from 'apollo-cache-inmemory';
import { setContext } from 'apollo-link-context';
import { HttpLink } from 'apollo-link-http';
import { onError } from 'apollo-link-error';
import { ApolloLink } from 'apollo-link';
import { URL } from '../config/api';

const withApollo = getClient => Component => {
  const WithApollo = props => {
    const ref = useRef(
      props.apollo ||
        getClient({
          initialState: props.apolloState
        })
    );
    return <Component {...props} apollo={ref.current} />;
  };

  WithApollo.getInitialProps = async appCtx => {
    const { ctx } = appCtx;
    const headers = ctx.req ? ctx.req.headers : {};
    const apollo = getClient({ ctx, headers });
    let apolloState = {};

    let props = { pageProps: {} };

    if (Component.getInitialProps) {
      ctx.apolloClient = apollo;
      props = await Component.getInitialProps(appCtx);
    }

    if (ctx.res && (ctx.res.headersSent || ctx.res.finished)) {
      return {};
    }

    apolloState = apollo.cache.extract();
    apollo.toJSON = () => null;

    return {
      ...props,
      apolloState,
      apollo
    };
  };

  return WithApollo;
};

export default withApollo(({ headers: reqHeaders = {}, initialState = {} }) => {
  return new ApolloClient({
    link: ApolloLink.from([
      onError(({ graphQLErrors, networkError }) => {
        if (process.env.NODE_ENV === 'development') {
          if (graphQLErrors)
            graphQLErrors.map(({ message, locations, path }) =>
              console.log(
                `[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`
              )
            );
          if (networkError) console.log(`[Network error]: ${networkError}`);
        }
      }),
      setContext((_, { headers = {} }) => {
        return {
          headers: {
            ...headers,
            cookie: reqHeaders.cookie || headers.cookie
          }
        };
      }),
      new HttpLink({
        uri: URL,
        credentials: 'include',
        fetch
      })
    ]),
    cache: new InMemoryCache().restore(initialState),
    ssrMode: !process.browser
  });
});
