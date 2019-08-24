import React from "react";
import { render, fireEvent, wait } from "@testing-library/react";
import { ToastContainer } from "react-toastify";
import { omit } from "lodash";
import i18n from "i18next";
import { ApolloProvider } from "react-apollo";

import { FETCH_CURRENT_USER_QUERY } from "@graphql/queries/user.queries";
import { SIGNUP_MUTATION } from "../../mutations";
import RegisterForm from "./RegisterForm.container";
import createClient from "@utils/test_utils/createClient";
import MockRouter from "@utils/test_utils/MockRouter";
import { users } from "@utils/test_utils/seed";
import { testID } from "@common/Form/TextField/constants";
import constants from "@config/constants";
import registerPageConstants from "../../constants";

const t = i18n.getFixedT(null, constants.NAMESPACES.registerPage);

const renderRegisterForm = (mocks = []) => {
  const client = createClient({ mocks });
  return {
    ...render(
      <MockRouter>
        <ApolloProvider client={client}>
          <RegisterForm t={t} />
          <ToastContainer />
        </ApolloProvider>
      </MockRouter>
    ),
    client
  };
};

describe("RegisterForm", () => {
  test("login, email and password are required", async () => {
    const { getAllByTestId, getByText, getAllByText } = renderRegisterForm();

    getAllByTestId(testID).forEach(el => {
      fireEvent.change(el, { target: { value: "" } });
      fireEvent.blur(el);
    });

    await wait(() =>
      expect(
        getAllByText(t("registerForm.errors.validation.mustProvidePassword"))
      ).toHaveLength(2)
    );

    [
      t("registerForm.errors.validation.mustProvideLogin"),
      t("registerForm.errors.validation.mustProvideEmail")
    ].forEach(text => {
      expect(getByText(text)).toBeInTheDocument();
    });
  });

  test(`length of username should be between ${constants.VALIDATION.minimumLengthOfLogin} and ${constants.VALIDATION.maximumLengthOfLogin} characters`, async () => {
    const { getAllByTestId, getByText } = renderRegisterForm();
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

  test("should show error when email is invalid", async () => {
    const { getAllByTestId, getByText } = renderRegisterForm();

    getAllByTestId(testID).forEach(el => {
      if (el.id === registerPageConstants.EMAIL) {
        fireEvent.change(el, { target: { value: "asdasd" } });
        fireEvent.blur(el);
      }
    });

    await wait(() =>
      expect(
        getByText(t("registerForm.errors.validation.invalidEmail"))
      ).toBeInTheDocument()
    );
  });

  test(`length of password should be between ${constants.VALIDATION.minimumLengthOfPassword} and ${constants.VALIDATION.maximumLengthOfPassword} characters`, async () => {
    const { getAllByTestId, getByText } = renderRegisterForm();
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
      k <= constants.VALIDATION.maximumLengthOfPassword + 40;
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

  test("password must contains at least 1 lowercase", async () => {
    const { getAllByTestId, getByText } = renderRegisterForm();

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

  test("password must contains at least 1 uppercase", async () => {
    const { getAllByTestId, getByText } = renderRegisterForm();

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

  test("password must contains at least 1 digit", async () => {
    const { getAllByTestId, getByText } = renderRegisterForm();

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

  test("passwordConfirmation must be the same as password", async () => {
    const { getAllByTestId, getByText } = renderRegisterForm();

    getAllByTestId(testID).forEach(el => {
      if (el.id === registerPageConstants.PASSWORD) {
        fireEvent.change(el, { target: { value: "asdasdasdASDa123" } });
        fireEvent.blur(el);
      } else if (el.id === registerPageConstants.PASSWORD_CONFIRMATION) {
        fireEvent.change(el, { target: { value: "asdasd123TASD" } });
        fireEvent.blur(el);
      }
    });

    await wait(() =>
      expect(
        getByText(t("registerForm.errors.validation.passwordsAreNotTheSame"))
      ).toBeInTheDocument()
    );
  });

  test("successful submit", async () => {
    const password = "examplePassword123T";
    const user = {
      login: users[0].login,
      password,
      passwordConfirmation: password,
      email: users[0].email
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
          query: SIGNUP_MUTATION,
          variables: {
            user: omit(user, [registerPageConstants.PASSWORD_CONFIRMATION])
          }
        },
        result: {
          data: {
            signup: users[0]
          }
        }
      }
    ];

    const {
      getAllByTestId,
      getByDisplayValue,
      getByTestId,
      getByText
    } = renderRegisterForm(mocks);

    getAllByTestId(testID).forEach(el => {
      fireEvent.change(el, { target: { value: user[el.id] } });
      fireEvent.blur(el);
    });

    await wait(() => expect(getByDisplayValue(user.login)).toBeInTheDocument());

    fireEvent.submit(getByTestId(registerPageConstants.REGISTER_FORM));

    await wait(() =>
      expect(getByText(t("registerForm.success"))).toBeInTheDocument()
    );
  });
});
