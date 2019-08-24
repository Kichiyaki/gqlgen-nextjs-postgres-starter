import React from "react";
import { render, fireEvent, wait } from "@testing-library/react";
import i18n from "i18next";

import { CHANGE_PASSWORD_MUTATION } from "../../mutations";
import ChangePasswordForm from "./ChangePasswordForm.container";
import createClient from "@utils/test_utils/createClient";
import TestLayout from "@utils/test_utils/TestLayout";
import { users } from "@utils/test_utils/seed";
import { testID } from "@common/Form/TextField/constants";
import constants from "@config/constants";
import pageConstants from "../../constants";
import { FETCH_CURRENT_USER_QUERY } from "@graphql/queries/user.queries";

const t = i18n.getFixedT(null, [
  pageConstants.NAMESPACE,
  constants.NAMESPACES.registerPage
]);

const renderChangePasswordForm = (mocks = []) => {
  const client = createClient({ mocks, user: users[0] });
  return {
    ...render(
      <TestLayout client={client}>
        <ChangePasswordForm t={t} />
      </TestLayout>
    ),
    client
  };
};

describe("ChangePasswordForm", () => {
  test("currentPassword and newPassword are required", async () => {
    const { getAllByTestId, getByText } = renderChangePasswordForm();

    getAllByTestId(testID).forEach(el => {
      fireEvent.change(el, { target: { value: "" } });
      fireEvent.blur(el);
    });

    await wait(() => {
      expect(
        getByText(
          t(
            constants.NAMESPACES.registerPage +
              ":registerForm.errors.validation.mustProvidePassword"
          )
        )
      ).toBeInTheDocument();
      expect(
        getByText(
          t("changePasswordForm.errors.validation.passwordsCannotBeTheSame")
        )
      ).toBeInTheDocument();
    });
  });

  test(`length of password should be between ${constants.VALIDATION.minimumLengthOfPassword} and ${constants.VALIDATION.maximumLengthOfPassword} characters`, async () => {
    const { getAllByTestId, getAllByText } = renderChangePasswordForm();
    let value = "asasd";

    getAllByTestId(testID).forEach(el => {
      fireEvent.change(el, {
        target: {
          value: el.id === pageConstants.NEW_PASSWORD ? value + "a" : value
        }
      });
      fireEvent.blur(el);
    });

    await wait(() =>
      expect(
        getAllByText(
          t(
            constants.NAMESPACES.registerPage +
              ":registerForm.errors.validation.minimumLengthOfPassword",
            {
              count: constants.VALIDATION.minimumLengthOfPassword
            }
          )
        )
      ).toHaveLength(2)
    );

    for (
      let k = 1;
      k <= constants.VALIDATION.maximumLengthOfPassword + 30;
      k++
    ) {
      value += "a";
    }

    getAllByTestId(testID).forEach(el => {
      fireEvent.change(el, {
        target: {
          value: el.id === pageConstants.NEW_PASSWORD ? value + "a" : value
        }
      });
      fireEvent.blur(el);
    });

    await wait(() =>
      expect(
        getAllByText(
          t(
            constants.NAMESPACES.registerPage +
              ":registerForm.errors.validation.maximumLengthOfPassword",
            {
              count: constants.VALIDATION.maximumLengthOfPassword
            }
          )
        )
      ).toHaveLength(2)
    );
  });

  test("password must contains at least 1 lowercase", async () => {
    const { getAllByTestId, getAllByText } = renderChangePasswordForm();

    const value = "ASDASDAASDASDASDA";

    getAllByTestId(testID).forEach(el => {
      fireEvent.change(el, {
        target: {
          value: el.id === pageConstants.NEW_PASSWORD ? value + "A" : value
        }
      });
      fireEvent.blur(el);
    });

    await wait(() =>
      expect(
        getAllByText(
          t(
            constants.NAMESPACES.registerPage +
              ":registerForm.errors.validation.passwordMustContainOneLowercase"
          )
        )
      ).toHaveLength(2)
    );
  });

  test("password must contains at least 1 uppercase", async () => {
    const { getAllByTestId, getAllByText } = renderChangePasswordForm();

    const value = "asdasdasdasdaadsa";

    getAllByTestId(testID).forEach(el => {
      fireEvent.change(el, {
        target: {
          value: el.id === pageConstants.NEW_PASSWORD ? value + "a" : value
        }
      });
      fireEvent.blur(el);
    });

    await wait(() =>
      expect(
        getAllByText(
          t(
            constants.NAMESPACES.registerPage +
              ":registerForm.errors.validation.passwordMustContainOneUppercase"
          )
        )
      ).toHaveLength(2)
    );
  });

  test("password must contains at least 1 digit", async () => {
    const { getAllByTestId, getAllByText } = renderChangePasswordForm();

    const value = "asdasdasdasASDdaadsa";

    getAllByTestId(testID).forEach(el => {
      fireEvent.change(el, {
        target: {
          value: el.id === pageConstants.NEW_PASSWORD ? value + "a" : value
        }
      });
      fireEvent.blur(el);
    });

    await wait(() =>
      expect(
        getAllByText(
          t(
            constants.NAMESPACES.registerPage +
              ":registerForm.errors.validation.passwordMustContainOneDigit"
          )
        )
      ).toHaveLength(2)
    );
  });

  test("currentPassword and newPassword cannot be the same", async () => {
    const { getAllByTestId, getByText } = renderChangePasswordForm();

    const value = "asdasdasdasASDda123asdadsa";

    getAllByTestId(testID).forEach(el => {
      fireEvent.change(el, { target: { value } });
      fireEvent.blur(el);
    });

    await wait(() =>
      expect(
        getByText(
          t("changePasswordForm.errors.validation.passwordsCannotBeTheSame")
        )
      ).toBeInTheDocument()
    );
  });

  test("successful submit", async () => {
    const values = {
      currentPassword: "currentPasswordas123",
      newPassword: "newPassword123"
    };

    const mocks = [
      {
        request: {
          query: CHANGE_PASSWORD_MUTATION,
          variables: values
        },
        result: {
          data: {
            changePassword: "Success"
          }
        }
      },
      {
        request: {
          query: FETCH_CURRENT_USER_QUERY
        },
        result: {
          data: {
            fetchCurrentUser: null
          }
        }
      }
    ];

    const {
      getAllByTestId,
      getByDisplayValue,
      getByTestId,
      getByText
    } = renderChangePasswordForm(mocks);

    getAllByTestId(testID).forEach(el => {
      fireEvent.change(el, { target: { value: values[el.id] } });
      fireEvent.blur(el);
    });

    await wait(() => {
      for (let i in values) {
        expect(getByDisplayValue(values[i])).toBeInTheDocument();
      }
    });

    fireEvent.submit(getByTestId(pageConstants.CHANGE_PASSWORD_FORM));

    await wait(() =>
      expect(getByText(t("changePasswordForm.success"))).toBeInTheDocument()
    );
  });
});
