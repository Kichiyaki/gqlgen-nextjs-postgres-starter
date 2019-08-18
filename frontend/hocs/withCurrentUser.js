import React from "react";
import { FETCH_CURRENT_USER_QUERY } from "../graphql/queries/user.queries";

const withCurrentUser = WrappedComponent => {
  const WithCurrentUser = props => {
    return <WrappedComponent {...props} />;
  };

  WithCurrentUser.getInitialProps = async ctx => {
    let pageProps = {};

    await ctx.apolloClient.query({
      query: FETCH_CURRENT_USER_QUERY,
      fetchPolicy: "network-only"
    });

    if (WrappedComponent.getInitialProps) {
      pageProps = await WrappedComponent.getInitialProps(ctx);
    }

    return pageProps;
  };

  return WithCurrentUser;
};

export default withCurrentUser;
