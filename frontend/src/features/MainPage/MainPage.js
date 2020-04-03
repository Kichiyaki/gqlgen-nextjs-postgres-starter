import React from 'react';
import { COMMON } from '@config/namespaces';
import AppLayout from '@common/AppLayout/AppLayout';

function MainPage() {
  return <AppLayout>MainPage</AppLayout>;
}

MainPage.getInitialProps = () => {
  return {
    namespacesRequired: [COMMON]
  };
};

export default MainPage;
