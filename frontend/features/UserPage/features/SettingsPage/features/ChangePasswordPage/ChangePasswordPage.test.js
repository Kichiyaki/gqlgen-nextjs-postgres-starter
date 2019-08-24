import React from "react";
import { render } from "@testing-library/react";
import i18n from "i18next";

import ChangePasswordPage from "./ChangePasswordPage";
import createClient from "@utils/test_utils/createClient";
import TestLayout from "@utils/test_utils/TestLayout";
import { users } from "@utils/test_utils/seed";
import pageConstants from "./constants";

const t = i18n.getFixedT(null, pageConstants.NAMESPACE);

const renderPage = (mocks = []) => {
  const client = createClient({ mocks, user: users[0] });
  return {
    ...render(
      <TestLayout client={client}>
        <ChangePasswordPage />
      </TestLayout>
    ),
    client
  };
};

describe("ChangePasswordPage", () => {
  test("should render ChangePasswordPage correctly", () => {
    const { asFragment, getByTestId, getByText } = renderPage();
    expect(asFragment()).toMatchSnapshot();
    expect(getByTestId(pageConstants.CHANGE_PASSWORD_FORM)).toBeInTheDocument();
    expect(getByText(t("title"))).toBeInTheDocument();
  });
});
