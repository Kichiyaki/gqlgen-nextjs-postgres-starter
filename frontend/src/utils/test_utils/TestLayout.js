import React from 'react';
import i18n from 'i18next';
import { createMuiTheme } from '@material-ui/core/styles';
import { ApolloProvider } from '@apollo/react-hooks';
import { I18nextProvider } from 'react-i18next';
import { ThemeProvider } from '@material-ui/styles';

import MockRouter from './MockRouter';
import createApolloClient from './createApolloClient';

const theme = createMuiTheme();

const TestLayout = ({ client, routerProps, children }) => {
  return (
    <MockRouter {...routerProps}>
      <ApolloProvider client={client}>
        <I18nextProvider i18n={i18n}>
          <ThemeProvider theme={theme}>{children}</ThemeProvider>
        </I18nextProvider>
      </ApolloProvider>
    </MockRouter>
  );
};

TestLayout.defaultProps = {
  client: createApolloClient(),
  routerProps: {}
};

export default TestLayout;
