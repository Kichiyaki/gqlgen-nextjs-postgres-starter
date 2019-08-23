import React from "react";
import { render, fireEvent, wait } from "@testing-library/react";
import i18n from "i18next";

import { GENERATE_NEW_ACTIVATION_TOKEN_FOR_CURRENT_USER } from "../../mutations";
import GenerateNewActivationTokenForm from "./GenerateNewActivationTokenForm";
import createClient from "@utils/test_utils/createClient";
import TestLayout from "@utils/test_utils/TestLayout";
import { users } from "@utils/test_utils/seed";
import pageConstants from "../../constants";

const t = i18n.getFixedT(null, pageConstants.NAMESPACE);

const renderGenerateNewActivationTokenForm = (mocks = []) => {
  const client = createClient({ mocks, user: users[0] });
  return {
    ...render(
      <TestLayout client={client}>
        <GenerateNewActivationTokenForm t={t} />
      </TestLayout>
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
      getByTestId(pageConstants.GENERATE_NEW_ACTIVATION_TOKEN_BUTTON_TESTID)
    );

    await wait(() => {
      expect(
        getByText(t("generateNewActivationTokenForm.success"))
      ).toBeInTheDocument();
    });
  });
});
