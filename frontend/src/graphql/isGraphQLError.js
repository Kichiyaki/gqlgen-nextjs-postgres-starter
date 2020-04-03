export default error => {
  return (
    Array.isArray(error.graphQLErrors) &&
    error.graphQLErrors.length > 0 &&
    error.graphQLErrors[0].message
  );
};
