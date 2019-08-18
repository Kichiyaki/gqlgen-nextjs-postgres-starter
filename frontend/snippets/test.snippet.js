import React from "react";
import { render } from "@testing-library/react";

import Page from "./Page";
import createClient from "@utils/test_utils/createClient";
import TestLayout from "@utils/test_utils/TestLayout";

const renderPage = (mocks = []) => {
  const client = createClient({ mocks });
  return {
    ...render(
      <TestLayout client={client}>
        <Page />
      </TestLayout>
    ),
    client
  };
};
