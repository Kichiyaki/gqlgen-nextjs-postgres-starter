export default class GraphQLError extends Error {
  constructor(message = '') {
    super(message);
    this.graphQLErrors = [{ message }];
  }
}
