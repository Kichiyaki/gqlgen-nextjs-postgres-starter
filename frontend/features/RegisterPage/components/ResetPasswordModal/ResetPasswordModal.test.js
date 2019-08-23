import React from "react";
import { render, fireEvent, wait } from "@testing-library/react";
import { ToastContainer } from "react-toastify";
import i18n from "i18next";
import { ApolloProvider } from "react-apollo";

import { GENERATE_NEW_RESET_PASSWORD_TOKEN_MUTATION } from "../../mutations";
import ResetPasswordModal from "./ResetPasswordModal";
import createClient from "@utils/test_utils/createClient";
import MockRouter from "@utils/test_utils/MockRouter";
import constants from "@config/constants";
import { testID } from "@common/Form/TextField/constants";
import pageConstants from "../../constants";

let handleClose;
const t = i18n.getFixedT(null, pageConstants.NAMESPACE);

const renderResetPasswordModal = (mocks = []) => {
  const client = createClient({ mocks });
  handleClose = jest.fn();
  return {
    ...render(
      <MockRouter>
        <ApolloProvider client={client}>
          <ResetPasswordModal handleClose={handleClose} open t={t} />
          <ToastContainer />
        </ApolloProvider>
      </MockRouter>
    ),
    client
  };
};

describe("ResetPasswordModal", () => {
  test("should render ResetPasswordModal correctly", () => {
    const { getByText } = renderResetPasswordModal();
    expect(getByText(t("resetPasswordModal.title"))).toBeInTheDocument();
    expect(getByText(t("resetPasswordModal.cancel"))).toBeInTheDocument();
    expect(getByText(t("resetPasswordModal.submitButton"))).toBeInTheDocument();
    expect(
      getByText(t("resetPasswordModal.inputLabels.email"))
    ).toBeInTheDocument();
  });

  test("email is required", async () => {
    const { getByText } = renderResetPasswordModal();
    fireEvent.click(
      getByText(t("resetPasswordModal.submitButton")).parentElement
    );

    await wait(() =>
      getByText(t("resetPasswordModal.errors.validation.mustProvideEmail"))
    );
  });

  test("email address must be valid", async () => {
    const {
      getByText,
      getByTestId,
      getByDisplayValue
    } = renderResetPasswordModal();
    const email = "asdf";
    fireEvent.change(getByTestId(testID), { target: { value: email } });

    await wait(() => getByDisplayValue(email));

    fireEvent.click(
      getByText(t("resetPasswordModal.submitButton")).parentElement
    );

    await wait(() =>
      getByText(
        t("resetPasswordModal.errors.validation.invalidEmail", { email })
      )
    );
  });

  test("successful submit", async () => {
    const email = "test@test.com";
    const mocks = [
      {
        request: {
          query: GENERATE_NEW_RESET_PASSWORD_TOKEN_MUTATION,
          variables: { email }
        },
        result: {
          data: {
            generateNewResetPasswordToken: "Sukces"
          }
        }
      }
    ];

    const {
      getByText,
      getByTestId,
      getByDisplayValue
    } = renderResetPasswordModal(mocks);
    fireEvent.change(getByTestId(testID), { target: { value: email } });

    await wait(() => getByDisplayValue(email));

    fireEvent.click(
      getByText(t("resetPasswordModal.submitButton")).parentElement
    );

    await wait(() => getByText(t("resetPasswordModal.success")));
  });
});
