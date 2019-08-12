import React from "react";
import { render } from "@testing-library/react";
import { ToastContainer } from "react-toastify";
import i18n from "i18next";

import AccountActivationPage from "./AccountActivationPage";
import createClient from "@utils/test_utils/createClient";
import MockRouter from "@utils/test_utils/MockRouter";
import MockI18nextProvider from "@utils/test_utils/MockI18nextProvider";
import ApolloProvider from "@common/ApolloProvider/ApolloProvider";
import { users } from "@utils/test_utils/seed";
import pageConstants from "./constants";

const t = i18n.getFixedT(null, pageConstants.NAMESPACE);

const renderPage = (mocks = []) => {
  const client = createClient({ mocks, user: users[0] });
  return {
    ...render(
      <MockRouter>
        <ApolloProvider client={client}>
          <MockI18nextProvider>
            <AccountActivationPage />
            <ToastContainer />
          </MockI18nextProvider>
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
      getByTestId(pageConstants.GENERATE_NEW_ACTIVATION_TOKEN_BUTTON_TESTID)
    ).toBeInTheDocument();
    expect(
      getByTestId(pageConstants.ACTIVATE_ACCOUNT_FORM)
    ).toBeInTheDocument();
    expect(getByText(t("title"))).toBeInTheDocument();
  });
});
