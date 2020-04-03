import React from 'react';
import App from 'next/app';
import { ApolloProvider } from '@apollo/react-hooks';
import { createMuiTheme, responsiveFontSizes } from '@material-ui/core/styles';
import withApollo from '@libs/withApollo';
import { appWithTranslation } from '@libs/i18n';
import { ME } from '@graphql/queries/auth.queries';
import { ThemeProvider } from '@material-ui/styles';

const theme = responsiveFontSizes(
  createMuiTheme({
    palette: {
      type: 'dark'
    }
  })
);

class MyApp extends App {
  static async getInitialProps({ Component, ctx }) {
    await ctx.apolloClient.query({ query: ME });
    let pageProps = {};
    if (Component.getInitialProps) {
      pageProps = await Component.getInitialProps(ctx);
    }

    return { pageProps };
  }

  componentDidMount() {
    const jssStyles = document.querySelector('#jss-server-side');
    if (jssStyles) {
      jssStyles.parentNode.removeChild(jssStyles);
    }
  }

  render() {
    const { Component, pageProps, apollo } = this.props;
    return (
      <ApolloProvider client={apollo}>
        <ThemeProvider theme={theme}>
          <Component {...pageProps} />
        </ThemeProvider>
      </ApolloProvider>
    );
  }
}

export default withApollo(appWithTranslation(MyApp));
