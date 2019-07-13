import React from "react";
import { render, wait } from "@testing-library/react";
import { ToastContainer } from "react-toastify";

import UserAccountActivationPage from "./UserAccountActivationPage";
import { users } from "@utils/test_utils/seed";
import createClient from "@utils/test_utils/createClient";
import MockRouter from "@utils/test_utils/MockRouter";
import TranslationProvider from "@lib/i18n/Provider";
import plTranslations from "@lib/i18n/translations/pl";
import ApolloProvider from "@common/ApolloProvider/ApolloProvider";
import { ACTIVATE_USER_ACCOUNT_QUERY } from "./queries";

const query = {
  id: 123,
  token: "asdd-asdd-asdd-asdd"
};
let push;

const renderPage = (mocks = []) => {
  const client = createClient({ mocks });
  push = jest.fn();
  return {
    ...render(
      <MockRouter push={push} query={query}>
        <ApolloProvider client={client}>
          <TranslationProvider locale="pl">
            <UserAccountActivationPage />
            <ToastContainer />
          </TranslationProvider>
        </ApolloProvider>
      </MockRouter>
    ),
    client
  };
};

describe("UserAccountActivationPage", () => {
  test("should correcly call onCompleted", async () => {
    const mocks = [
      {
        request: {
          query: ACTIVATE_USER_ACCOUNT_QUERY,
          variables: query
        },
        result: {
          data: {
            activateUserAccount: users[0]
          }
        }
      }
    ];

    const { getByText } = renderPage(mocks);
    await wait(() => {
      expect(
        getByText(
          plTranslations.USER_ACCOUNT_ACTIVATION_PAGE.success(users[0].login)
        )
      ).toBeInTheDocument();
      expect(push).toHaveBeenCalled();
    });
  });

  test("should correctly call onError", async () => {
    const errMsg = "example error msg";
    const mocks = [
      {
        request: {
          query: ACTIVATE_USER_ACCOUNT_QUERY,
          variables: query
        },
        result: {
          data: null,
          errors: [
            {
              message: errMsg
            }
          ]
        }
      }
    ];
    const { getByText } = renderPage(mocks);
    await wait(() => {
      expect(getByText(errMsg)).toBeInTheDocument();
      expect(push).toHaveBeenCalled();
    });
  });
});
