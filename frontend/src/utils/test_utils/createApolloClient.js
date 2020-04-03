import { ApolloClient } from '@apollo/client';
import { MockLink } from '@apollo/react-testing';
import { InMemoryCache } from 'apollo-cache-inmemory';
import { ME } from '@graphql/queries/auth.queries';

export default ({
  mocks = [],
  cache = new InMemoryCache({ addTypename: false }),
  user = undefined
} = {}) => {
  if (user) {
    cache.writeQuery({
      query: ME,
      data: {
        me: user
      }
    });
  }
  return new ApolloClient({
    cache,
    link: new MockLink(mocks)
  });
};
