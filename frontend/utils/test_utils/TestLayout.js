import React from "react";
import { ToastContainer } from "react-toastify";

import MockRouter from "./MockRouter";
import MockI18nextProvider from "./MockI18nextProvider";
import ApolloProvider from "@common/ApolloProvider/ApolloProvider";
import createClient from "./createClient";

const TestLayout = ({ client, routerProps, children }) => {
  return (
    <MockRouter {...routerProps}>
      <ApolloProvider client={client}>
        <MockI18nextProvider>
          {children}
          <ToastContainer />
        </MockI18nextProvider>
      </ApolloProvider>
    </MockRouter>
  );
};

TestLayout.defaultProps = {
  client: createClient(),
  routerProps: {}
};

export default TestLayout;
