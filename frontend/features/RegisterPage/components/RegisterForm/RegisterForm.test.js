import React from "react";
import { render, fireEvent, wait } from "@testing-library/react";
import { ToastContainer } from "react-toastify";
import { omit } from "lodash";

import { FETCH_CURRENT_USER_QUERY } from "@graphql/queries/user.queries";
import { SIGNUP_MUTATION } from "../../mutations";
import RegisterForm from "./RegisterForm.container";
import createClient from "@utils/test_utils/createClient";
import MockRouter from "@utils/test_utils/MockRouter";
import { users } from "@utils/test_utils/seed";
import ApolloProvider from "@common/ApolloProvider/ApolloProvider";
import { testID } from "@common/Form/TextField/constants";
import constants from "@config/constants";
import registerPageConstants from "../../constants";
import translations from "@lib/i18n/translations/pl";

const renderRegisterForm = (mocks = []) => {
  const client = createClient({ mocks });
  return {
    ...render(
      <MockRouter>
        <ApolloProvider client={client}>
          <RegisterForm translations={translations} />
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
        getAllByText(
          translations.REGISTER_PAGE.registerForm.errors.validation
            .mustProvidePassword
        )
      ).toHaveLength(2)
    );

    [
      translations.REGISTER_PAGE.registerForm.errors.validation
        .mustProvideLogin,
      translations.REGISTER_PAGE.registerForm.errors.validation.mustProvideEmail
    ].forEach(text => {
      expect(getByText(text)).toBeInTheDocument();
    });
  });

  test(`length of username should be between ${
    constants.VALIDATION.minimumLengthOfLogin
  } and ${constants.VALIDATION.maximumLengthOfLogin} characters`, async () => {
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
          translations.REGISTER_PAGE.registerForm.errors.validation
            .minimumLengthOfLogin
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
          translations.REGISTER_PAGE.registerForm.errors.validation
            .maximumLengthOfLogin
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
        getByText(
          translations.REGISTER_PAGE.registerForm.errors.validation.invalidEmail
        )
      ).toBeInTheDocument()
    );
  });

  test(`length of password should be between ${
    constants.VALIDATION.minimumLengthOfPassword
  } and ${
    constants.VALIDATION.maximumLengthOfPassword
  } characters`, async () => {
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
          translations.REGISTER_PAGE.registerForm.errors.validation
            .minimumLengthOfPassword
        )
      ).toBeInTheDocument()
    );

    for (let k = 1; k <= 140; k++) {
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
          translations.REGISTER_PAGE.registerForm.errors.validation
            .maximumLengthOfPassword
        )
      ).toBeInTheDocument()
    );
  });

  test("password must contains min. 1 lowercase", async () => {
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
          translations.REGISTER_PAGE.registerForm.errors.validation
            .passwordMustContainsOneLowercase
        )
      ).toBeInTheDocument()
    );
  });

  test("password must contains min. 1 uppercase", async () => {
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
          translations.REGISTER_PAGE.registerForm.errors.validation
            .passwordMustContainsOneUppercase
        )
      ).toBeInTheDocument()
    );
  });

  test("password must contains min. 1 digit", async () => {
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
          translations.REGISTER_PAGE.registerForm.errors.validation
            .passwordMustContainsOneDigit
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
        getByText(
          translations.REGISTER_PAGE.registerForm.errors.validation
            .passwordsAreNotTheSame
        )
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
      expect(
        getByText(translations.REGISTER_PAGE.registerForm.success)
      ).toBeInTheDocument()
    );
  });
});
