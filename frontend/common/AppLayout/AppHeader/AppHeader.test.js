import React from "react";
import { render, fireEvent, wait } from "@testing-library/react";
import { ToastContainer } from "react-toastify";

import createClient from "@utils/test_utils/createClient";
import MockRouter from "@utils/test_utils/MockRouter";
import MockI18nextProvider from "@utils/test_utils/MockI18nextProvider";
import { users } from "@utils/test_utils/seed";
import ApolloProvider from "@common/ApolloProvider/ApolloProvider";
import common from "@static/locales/pl/common.json";
import { FETCH_CURRENT_USER_QUERY } from "@graphql/queries/user.queries";
import { LOGOUT_USER_MUTATION } from "./mutations";
import constants from "./constants";
import AppHeader from "./AppHeader";

const renderHeader = (mocks = [], user = undefined) => {
  const client = createClient({ mocks, user });
  return {
    ...render(
      <MockRouter>
        <MockI18nextProvider>
          <ApolloProvider client={client}>
            <AppHeader />
            <ToastContainer />
          </ApolloProvider>
        </MockI18nextProvider>
      </MockRouter>
    ),
    client
  };
};

describe("AppHeader", () => {
  test("should render header when user is logged out correctly", () => {
    const { asFragment, queryByTestId, getByText } = renderHeader();
    expect(asFragment()).toMatchSnapshot();
    expect(getByText(common.APPLICATION.name)).toBeInTheDocument();
    expect(getByText(common.HEADER.buttons.logout)).toBeInTheDocument();
    expect(queryByTestId(constants.LOGOUT_BUTTON)).not.toBeInTheDocument();
  });

  test("should render header when user is logged in correctly", () => {
    const { asFragment, getByTestId, getByText } = renderHeader([], users[0]);
    expect(asFragment()).toMatchSnapshot();
    expect(getByText(common.APPLICATION.name)).toBeInTheDocument();
    expect(getByText(common.HEADER.buttons.logout)).toBeInTheDocument();
    expect(getByTestId(constants.LOGOUT_BUTTON)).toBeInTheDocument();
    expect(getByTestId(constants.LOGOUT_BUTTON)).toHaveTextContent(
      common.HEADER.buttons.logout
    );
  });

  test("should render header when user is logged in and has a deactivated account correctly", () => {
    const { asFragment, getByTestId, getByText } = renderHeader(
      [],
      users.find(user => !user.activated)
    );
    expect(asFragment()).toMatchSnapshot();
    expect(getByText(common.APPLICATION.name)).toBeInTheDocument();
    expect(getByText(common.HEADER.links.activateAccount)).toBeInTheDocument();
    expect(getByText(common.HEADER.buttons.logout)).toBeInTheDocument();
    expect(getByTestId(constants.LOGOUT_BUTTON)).toBeInTheDocument();
    expect(getByTestId(constants.LOGOUT_BUTTON)).toHaveTextContent(
      common.HEADER.buttons.logout
    );
  });

  test("should correctly logout user", async () => {
    const mocks = [
      {
        request: {
          query: FETCH_CURRENT_USER_QUERY
        },
        result: {
          data: {
            fetchCurrentUser: null
          }
        }
      },
      {
        request: {
          query: LOGOUT_USER_MUTATION
        },
        result: {
          data: {
            logout: "Success"
          }
        }
      }
    ];

    const { getByTestId, getByText } = renderHeader(mocks, users[0]);
    fireEvent.click(getByTestId(constants.LOGOUT_BUTTON));

    await wait(() => {
      expect(getByText(common.HEADER.logout.success)).toBeInTheDocument();
    });
  });
});
