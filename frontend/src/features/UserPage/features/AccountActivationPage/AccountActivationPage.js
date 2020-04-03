import React from 'react';
import i18n from 'i18next';
import isUUID from 'validator/lib/isUUID';
import { useTranslation } from '@libs/i18n';
import GraphQLError from '@graphql/GraphQLError';
import isGraphQLError from '@graphql/isGraphQLError';
import { COMMON, USER_PAGE } from '@config/namespaces';
import { ACTIVATE_USER_ACCOUNT_QUERY } from './constants';

import { makeStyles } from '@material-ui/core/styles';
import { Typography, Container } from '@material-ui/core';
import ErrorPage from '@features/ErrorPage/ErrorPage';
import AppLayout from '@common/AppLayout/AppLayout';

const useStyles = makeStyles(theme => ({
  image: {
    maxHeight: '25vh'
  },
  container: {
    '& > *:not(:last-child)': {
      marginBottom: theme.spacing(2)
    }
  },
  appLayout: {
    display: 'flex',
    justifyContent: 'center',
    flexDirection: 'column',
    textAlign: 'center'
  }
}));

export default function AccountActivationPage({ status, message }) {
  const classes = useStyles();
  const { t } = useTranslation(USER_PAGE.ACCOUNT_ACTIVATION_PAGE);

  if (status != 200) {
    return <ErrorPage title={message} statusCode={status} />;
  }

  return (
    <AppLayout className={classes.appLayout}>
      <Container maxWidth="sm">
        <Typography variant="h2" component="h1">
          {t('title')}
        </Typography>
        <Typography variant="h3" component="h2">
          {t('success')}
        </Typography>
      </Container>
    </AppLayout>
  );
}

AccountActivationPage.getInitialProps = async ({
  query,
  apolloClient,
  req
}) => {
  const props = {
    namespacesRequired: [COMMON, USER_PAGE.ACCOUNT_ACTIVATION_PAGE],
    status: 200
  };

  try {
    if (
      !query.id ||
      !query.token ||
      isNaN(parseInt(query.id)) ||
      !isUUID(query.token, '4')
    ) {
      throw new GraphQLError(
        req
          ? req.t(`${USER_PAGE.ACCOUNT_ACTIVATION_PAGE}:defaultError`)
          : i18n.t(`${USER_PAGE.ACCOUNT_ACTIVATION_PAGE}:defaultError`)
      );
    }
    await apolloClient.query({
      query: ACTIVATE_USER_ACCOUNT_QUERY,
      variables: query
    });
  } catch (error) {
    if (isGraphQLError(error)) {
      props.message = error.graphQLErrors[0].message;
    }
    props.status = 500;
  }

  return props;
};
