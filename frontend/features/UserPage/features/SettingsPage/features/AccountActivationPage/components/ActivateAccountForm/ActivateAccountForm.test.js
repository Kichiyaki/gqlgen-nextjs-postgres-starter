import React from "react";
import { render, fireEvent, wait } from "@testing-library/react";
import { ToastContainer } from "react-toastify";

import { ACTIVATE_USER_ACCOUNT_MUTATION } from "../../mutations";
import ActivateAccountForm from "./ActivateAccountForm.container";
import createClient from "@utils/test_utils/createClient";
import MockRouter from "@utils/test_utils/MockRouter";
import { users } from "@utils/test_utils/seed";
import ApolloProvider from "@common/ApolloProvider/ApolloProvider";
import { testID } from "@common/Form/TextField/constants";
import accountActivationPageConstants from "../../constants";
import translations from "@lib/i18n/translations/pl";

const {
  USER_PAGE: {
    SETTINGS_PAGE: {
      ACCOUNT_ACTIVATION_PAGE: { activateAccountForm }
    }
  }
} = translations;

const renderActivateAccountForm = (mocks = []) => {
  const client = createClient({ mocks, user: users[0] });
  return {
    ...render(
      <MockRouter>
        <ApolloProvider client={client}>
          <ActivateAccountForm translations={translations} />
          <ToastContainer />
        </ApolloProvider>
      </MockRouter>
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
      [activateAccountForm.errors.validation.mustProvideToken].forEach(text => {
        expect(getByText(text)).toBeInTheDocument();
      });
    });
  });

  test("token must be uuid", async () => {
    const { getAllByTestId, getByText } = renderActivateAccountForm();

    getAllByTestId(testID).forEach(el => {
      if (el.id === accountActivationPageConstants.TOKEN) {
        fireEvent.change(el, { target: { value: "asdasdadsada" } });
        fireEvent.blur(el);
      }
    });

    await wait(() => {
      [activateAccountForm.errors.validation.tokenIsInvalid].forEach(text => {
        expect(getByText(text)).toBeInTheDocument();
      });
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

    fireEvent.submit(
      getByTestId(accountActivationPageConstants.ACTIVATE_ACCOUNT_FORM)
    );

    await wait(() =>
      expect(getByText(activateAccountForm.success)).toBeInTheDocument()
    );
  });
});
