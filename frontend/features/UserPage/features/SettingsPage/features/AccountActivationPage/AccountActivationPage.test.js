import React from "react";
import { render } from "@testing-library/react";
import { ToastContainer } from "react-toastify";

import AccountActivationPage from "./AccountActivationPage";
import createClient from "@utils/test_utils/createClient";
import MockRouter from "@utils/test_utils/MockRouter";
import TranslationProvider from "@lib/i18n/Provider";
import translations from "@lib/i18n/translations/pl";
import ApolloProvider from "@common/ApolloProvider/ApolloProvider";
import { users } from "@utils/test_utils/seed";
import accountActivationPageConstants from "./constants";

const {
  USER_PAGE: {
    SETTINGS_PAGE: { ACCOUNT_ACTIVATION_PAGE }
  }
} = translations;

const renderPage = (mocks = []) => {
  const client = createClient({ mocks, user: users[0] });
  return {
    ...render(
      <MockRouter>
        <ApolloProvider client={client}>
          <TranslationProvider locale="pl">
            <AccountActivationPage />
            <ToastContainer />
          </TranslationProvider>
        </ApolloProvider>
      </MockRouter>
    ),
    client
  };
};

describe("AccountActivationPage", () => {
  test("should render AccountActivationPage correctly", () => {
    const { asFragment, getByTestId, getByText } = renderPage();
    expect(asFragment()).toMatchSnapshot();
    expect(
      getByTestId(
        accountActivationPageConstants.GENERATE_NEW_ACTIVATION_TOKEN_BUTTON_TESTID
      )
    ).toBeInTheDocument();
    expect(
      getByTestId(accountActivationPageConstants.ACTIVATE_ACCOUNT_FORM)
    ).toBeInTheDocument();
    expect(getByText(ACCOUNT_ACTIVATION_PAGE.title)).toBeInTheDocument();
  });
});
