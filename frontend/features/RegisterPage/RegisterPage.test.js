import React from "react";
import { render } from "@testing-library/react";
import { ToastContainer } from "react-toastify";

import RegisterPage from "./RegisterPage";
import createClient from "@utils/test_utils/createClient";
import MockRouter from "@utils/test_utils/MockRouter";
import TranslationProvider from "@lib/i18n/Provider";
import ApolloProvider from "@common/ApolloProvider/ApolloProvider";
import { testID } from "@common/Form/TextField/constants";
import constants from "@config/constants";
import registerPageConstants from "./constants";

const renderPage = (mocks = [], registerPage = true) => {
  const client = createClient({ mocks });
  return {
    ...render(
      <MockRouter
        route={
          registerPage ? constants.ROUTES.register : constants.ROUTES.login
        }
      >
        <ApolloProvider client={client}>
          <TranslationProvider locale="pl">
            <RegisterPage />
            <ToastContainer />
          </TranslationProvider>
        </ApolloProvider>
      </MockRouter>
    ),
    client
  };
};

describe("RegisterPage", () => {
  test("should render RegiserPage correctly", () => {
    const { asFragment, getByTestId, getAllByTestId } = renderPage();
    expect(asFragment()).toMatchSnapshot();
    expect(
      getByTestId(registerPageConstants.REGISTER_FORM)
    ).toBeInTheDocument();
    expect(getAllByTestId(testID)).toHaveLength(4);
  });

  test("should render LoginPage correctly", () => {
    const { asFragment, getByTestId, getAllByTestId } = renderPage([], false);
    expect(asFragment()).toMatchSnapshot();
    expect(getByTestId(registerPageConstants.LOGIN_FORM)).toBeInTheDocument();
    expect(getAllByTestId(testID)).toHaveLength(2);
  });
});
