import React from "react";
import { render, fireEvent, wait } from "@testing-library/react";
import { ToastContainer } from "react-toastify";

import { GENERATE_NEW_ACTIVATION_TOKEN_FOR_CURRENT_USER } from "../../mutations";
import GenerateNewActivationTokenForm from "./GenerateNewActivationTokenForm";
import createClient from "@utils/test_utils/createClient";
import MockRouter from "@utils/test_utils/MockRouter";
import { users } from "@utils/test_utils/seed";
import ApolloProvider from "@common/ApolloProvider/ApolloProvider";
import accountActivationPageConstants from "../../constants";
import translations from "@lib/i18n/translations/pl";

const {
  USER_PAGE: {
    SETTINGS_PAGE: {
      ACCOUNT_ACTIVATION_PAGE: { generateNewActivationTokenForm }
    }
  }
} = translations;

const renderGenerateNewActivationTokenForm = (mocks = []) => {
  const client = createClient({ mocks, user: users[0] });
  return {
    ...render(
      <MockRouter>
        <ApolloProvider client={client}>
          <GenerateNewActivationTokenForm translations={translations} />
          <ToastContainer />
        </ApolloProvider>
      </MockRouter>
    ),
    client
  };
};

describe("GenerateNewActivationTokenForm", () => {
  test("successful submit", async () => {
    const mocks = [
      {
        request: {
          query: GENERATE_NEW_ACTIVATION_TOKEN_FOR_CURRENT_USER
        },
        result: {
          data: {
            generateNewActivationTokenForCurrentUser: "Sukces"
          }
        }
      }
    ];

    const { getByTestId, getByText } = renderGenerateNewActivationTokenForm(
      mocks
    );

    fireEvent.click(
      getByTestId(
        accountActivationPageConstants.GENERATE_NEW_ACTIVATION_TOKEN_BUTTON_TESTID
      )
    );

    await wait(() => {
      expect(
        getByText(generateNewActivationTokenForm.success)
      ).toBeInTheDocument();
    });
  });
});
