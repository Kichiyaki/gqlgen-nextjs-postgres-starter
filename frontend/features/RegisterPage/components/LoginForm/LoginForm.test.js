import React from "react";
import { render, fireEvent, wait } from "@testing-library/react";
import { ToastContainer } from "react-toastify";
import { ApolloProvider } from "react-apollo";
import i18n from "i18next";

import constants from "@config/constants";
import { FETCH_CURRENT_USER_QUERY } from "@graphql/queries/user.queries";
import { LOGIN_MUTATION } from "../../mutations";
import LoginForm from "./LoginForm.container";
import createClient from "@utils/test_utils/createClient";
import MockRouter from "@utils/test_utils/MockRouter";
import { users } from "@utils/test_utils/seed";
import { testID } from "@common/Form/TextField/constants";
import registerPageConstants from "../../constants";

const t = i18n.getFixedT(null, constants.NAMESPACES.registerPage);

const renderLoginForm = (mocks = []) => {
  const client = createClient({ mocks });
  return {
    ...render(
      <MockRouter>
        <ApolloProvider client={client}>
          <LoginForm t={t} />
          <ToastContainer />
        </ApolloProvider>
      </MockRouter>
    ),
    client
  };
};

describe("LoginForm", () => {
  test("login and password are required", async () => {
    const { getAllByTestId, getByText } = renderLoginForm();

    getAllByTestId(testID).forEach(el => {
      fireEvent.change(el, { target: { value: "" } });
      fireEvent.blur(el);
    });

    await wait(() => {
      [
        t("registerForm.errors.validation.mustProvideLogin"),
        t("registerForm.errors.validation.mustProvidePassword")
      ].forEach(text => {
        expect(getByText(text)).toBeInTheDocument();
      });
    });
  });

  test(`length of username should be between ${
    constants.VALIDATION.minimumLengthOfLogin
  } and ${constants.VALIDATION.maximumLengthOfLogin} characters`, async () => {
    const { getAllByTestId, getByText } = renderLoginForm();
    let value = "as";

    getAllByTestId(testID).forEach(el => {
      if (el.id === registerPageConstants.LOGIN) {
        fireEvent.change(el, { target: { value } });
        fireEvent.blur(el);
      }
    });

    await wait(() =>
      expect(
        getByText(
          t("registerForm.errors.validation.minimumLengthOfLogin", {
            count: constants.VALIDATION.minimumLengthOfLogin
          })
        )
      ).toBeInTheDocument()
    );

    for (let k = 1; k <= constants.VALIDATION.maximumLengthOfLogin + 5; k++) {
      value += "a";
    }

    getAllByTestId(testID).forEach(el => {
      if (el.id === registerPageConstants.LOGIN) {
        fireEvent.change(el, {
          target: {
            value
          }
        });
        fireEvent.blur(el);
      }
    });

    await wait(() =>
      expect(
        getByText(
          t("registerForm.errors.validation.maximumLengthOfLogin", {
            count: constants.VALIDATION.maximumLengthOfLogin
          })
        )
      ).toBeInTheDocument()
    );
  });

  test(`length of password should be between ${
    constants.VALIDATION.minimumLengthOfPassword
  } and ${
    constants.VALIDATION.maximumLengthOfPassword
  } characters`, async () => {
    const { getAllByTestId, getByText } = renderLoginForm();
    let value = "asasd";

    getAllByTestId(testID).forEach(el => {
      if (el.id === registerPageConstants.PASSWORD) {
        fireEvent.change(el, { target: { value } });
        fireEvent.blur(el);
      }
    });

    await wait(() =>
      expect(
        getByText(
          t("registerForm.errors.validation.minimumLengthOfPassword", {
            count: constants.VALIDATION.minimumLengthOfPassword
          })
        )
      ).toBeInTheDocument()
    );

    for (
      let k = 1;
      k <= constants.VALIDATION.maximumLengthOfPassword + 5;
      k++
    ) {
      value += "a";
    }

    getAllByTestId(testID).forEach(el => {
      if (el.id === registerPageConstants.PASSWORD) {
        fireEvent.change(el, {
          target: {
            value
          }
        });
        fireEvent.blur(el);
      }
    });

    await wait(() =>
      expect(
        getByText(
          t("registerForm.errors.validation.maximumLengthOfPassword", {
            count: constants.VALIDATION.maximumLengthOfPassword
          })
        )
      ).toBeInTheDocument()
    );
  });

  test("password must contains min. 1 lowercase", async () => {
    const { getAllByTestId, getByText } = renderLoginForm();

    getAllByTestId(testID).forEach(el => {
      if (el.id === registerPageConstants.PASSWORD) {
        fireEvent.change(el, { target: { value: "ASDASDAASDASDASDA" } });
        fireEvent.blur(el);
      }
    });

    await wait(() =>
      expect(
        getByText(
          t("registerForm.errors.validation.passwordMustContainOneLowercase")
        )
      ).toBeInTheDocument()
    );
  });

  test("password must contains min. 1 uppercase", async () => {
    const { getAllByTestId, getByText } = renderLoginForm();

    getAllByTestId(testID).forEach(el => {
      if (el.id === registerPageConstants.PASSWORD) {
        fireEvent.change(el, { target: { value: "asdasdasdasdaadsa" } });
        fireEvent.blur(el);
      }
    });

    await wait(() =>
      expect(
        getByText(
          t("registerForm.errors.validation.passwordMustContainOneUppercase")
        )
      ).toBeInTheDocument()
    );
  });

  test("password must contains min. 1 digit", async () => {
    const { getAllByTestId, getByText } = renderLoginForm();

    getAllByTestId(testID).forEach(el => {
      if (el.id === registerPageConstants.PASSWORD) {
        fireEvent.change(el, { target: { value: "asdasdasdasASDdaadsa" } });
        fireEvent.blur(el);
      }
    });

    await wait(() =>
      expect(
        getByText(
          t("registerForm.errors.validation.passwordMustContainOneDigit")
        )
      ).toBeInTheDocument()
    );
  });

  test("successful submit", async () => {
    const password = "examplePassword123T";
    const user = {
      login: users[0].login,
      password
    };

    const mocks = [
      {
        request: {
          query: FETCH_CURRENT_USER_QUERY
        },
        result: {
          data: {
            fetchCurrentUser: users[0]
          }
        }
      },
      {
        request: {
          query: LOGIN_MUTATION,
          variables: {
            login: user.login,
            password
          }
        },
        result: {
          data: {
            login: users[0]
          }
        }
      }
    ];

    const {
      getAllByTestId,
      getByDisplayValue,
      getByTestId,
      getByText
    } = renderLoginForm(mocks);

    getAllByTestId(testID).forEach(el => {
      fireEvent.change(el, { target: { value: user[el.id] } });
      fireEvent.blur(el);
    });

    await wait(() => expect(getByDisplayValue(user.login)).toBeInTheDocument());

    fireEvent.submit(getByTestId(registerPageConstants.LOGIN_FORM));

    await wait(() =>
      expect(getByText(t("loginForm.success"))).toBeInTheDocument()
    );
  });
});
