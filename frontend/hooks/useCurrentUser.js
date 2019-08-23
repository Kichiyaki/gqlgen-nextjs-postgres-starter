import { useQuery } from "react-apollo";
import { FETCH_CURRENT_USER_QUERY } from "../graphql/queries/user.queries";

export default (options = { fetchPolicy: "cache-only" }) => {
  return useQuery(FETCH_CURRENT_USER_QUERY, options);
};
