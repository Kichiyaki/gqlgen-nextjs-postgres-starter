import { ApolloClient } from "apollo-client";
import { InMemoryCache } from "apollo-cache-inmemory";
import { setContext } from "apollo-link-context";
import { HttpLink } from "apollo-link-http";
import { onError } from "apollo-link-error";
import { ApolloLink } from "apollo-link";
import { BASE_URL } from "../config/api";
import withApollo from "../lib/next-with-apollo";

export default withApollo(
  ({ headers: reqHeaders = {}, initialState }) =>
    new ApolloClient({
      link: ApolloLink.from([
        onError(({ graphQLErrors, networkError }) => {
          if (process.env.NODE_ENV === "development") {
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
          uri: BASE_URL,
          credentials: "include"
        })
      ]),
      cache: new InMemoryCache().restore(initialState || {}),
      ssrMode: !process.browser
    })
);
