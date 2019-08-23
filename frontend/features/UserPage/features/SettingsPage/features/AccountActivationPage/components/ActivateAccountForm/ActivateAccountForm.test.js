import React from "react";
import { render, fireEvent, wait } from "@testing-library/react";
import i18n from "i18next";

import { ACTIVATE_USER_ACCOUNT_MUTATION } from "../../mutations";
import ActivateAccountForm from "./ActivateAccountForm.container";
import createClient from "@utils/test_utils/createClient";
import TestLayout from "@utils/test_utils/TestLayout";
import { users } from "@utils/test_utils/seed";
import { testID } from "@common/Form/TextField/constants";
import pageConstants from "../../constants";

const t = i18n.getFixedT(null, pageConstants.NAMESPACE);

const renderActivateAccountForm = (mocks = []) => {
  const client = createClient({ mocks, user: users[0] });
  return {
    ...render(
      <TestLayout client={client}>
        <ActivateAccountForm t={t} />
      </TestLayout>
    ),
    client
  };
};

describe("ActivateAccountForm", () => {
  test("token is required", async () => {
    const { getAllByTestId, getByText } = renderActivateAccountForm();

    getAllByTestId(testID).forEach(el => {
      fireEvent.change(el, { target: { value: "" } });
      fireEvent.blur(el);
    });

    await wait(() => {
      [t("activateAccountForm.errors.validation.mustProvideToken")].forEach(
        text => {
          expect(getByText(text)).toBeInTheDocument();
        }
      );
    });
  });

  test("token must be uuid", async () => {
    const { getAllByTestId, getByText } = renderActivateAccountForm();

    getAllByTestId(testID).forEach(el => {
      if (el.id === pageConstants.TOKEN) {
        fireEvent.change(el, { target: { value: "asdasdadsada" } });
        fireEvent.blur(el);
      }
    });

    await wait(() => {
      [t("activateAccountForm.errors.validation.tokenIsInvalid")].forEach(
        text => {
          expect(getByText(text)).toBeInTheDocument();
        }
      );
    });
  });

  test("successful submit", async () => {
    const values = {
      token: "80131cba-6cb7-4799-98fd-8c6e117218fe"
    };

    const mocks = [
      {
        request: {
          query: ACTIVATE_USER_ACCOUNT_MUTATION,
          variables: {
            token: values.token,
            id: users[0].id
          }
        },
        result: {
          data: {
            activateUserAccount: users[0]
          }
        }
      }
    ];

    const {
      getAllByTestId,
      getByDisplayValue,
      getByTestId,
      getByText
    } = renderActivateAccountForm(mocks);

    getAllByTestId(testID).forEach(el => {
      fireEvent.change(el, { target: { value: values[el.id] } });
      fireEvent.blur(el);
    });

    await wait(() =>
      expect(getByDisplayValue(values.token)).toBeInTheDocument()
    );

    fireEvent.submit(getByTestId(pageConstants.ACTIVATE_ACCOUNT_FORM));

    await wait(() =>
      expect(getByText(t("activateAccountForm.success"))).toBeInTheDocument()
    );
  });
});
