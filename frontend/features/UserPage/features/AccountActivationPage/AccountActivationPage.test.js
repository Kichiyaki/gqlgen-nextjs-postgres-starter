import React from "react";
import { render, wait } from "@testing-library/react";
import i18n from "i18next";

import AccountActivationPage from "./AccountActivationPage";
import { users } from "@utils/test_utils/seed";
import createClient from "@utils/test_utils/createClient";

import TestLayout from "@utils/test_utils/TestLayout";
import { ACTIVATE_USER_ACCOUNT_QUERY } from "./queries";
import pageConstants from "./constants";

const query = {
  id: 123,
  token: "asdd-asdd-asdd-asdd"
};
let push;
const t = i18n.getFixedT(null, pageConstants.NAMESPACE);

const renderPage = (mocks = []) => {
  const client = createClient({ mocks });
  push = jest.fn();
  return {
    ...render(
      <TestLayout client={client} routerProps={{ push, query }}>
        <AccountActivationPage />
      </TestLayout>
    ),
    client
  };
};

describe("UserPage/AccountActivationPage", () => {
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
        getByText(t("success", { login: users[0].login }))
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
