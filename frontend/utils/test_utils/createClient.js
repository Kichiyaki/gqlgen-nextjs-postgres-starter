import { ApolloClient } from "apollo-client";
import { MockLink } from "react-apollo/test-utils";
import { InMemoryCache } from "apollo-cache-inmemory";
import { FETCH_CURRENT_USER_QUERY } from "../../graphql/queries/user.queries";

export default ({
  mocks = [],
  cache = new InMemoryCache({ addTypename: false }),
  user = undefined
} = {}) => {
  if (user) {
    cache.writeQuery({
      query: FETCH_CURRENT_USER_QUERY,
      data: {
        fetchCurrentUser: user
      }
    });
  }
  return new ApolloClient({
    cache,
    link: new MockLink(mocks)
  });
};
