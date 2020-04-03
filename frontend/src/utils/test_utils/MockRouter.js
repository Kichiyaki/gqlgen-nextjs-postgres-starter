import React from 'react';
import { RouterContext } from 'next/dist/next-server/lib/router-context';
import { string as propString, object, func } from 'prop-types';

const MockRouter = ({
  children,
  route,
  pathname,
  query,
  asPath,
  back,
  beforePopState,
  on,
  off,
  emit,
  prefetch,
  push,
  reload,
  replace
}) => {
  return (
    <RouterContext.Provider
      value={{
        route,
        pathname,
        query,
        asPath,
        back,
        beforePopState,
        prefetch,
        push,
        reload,
        replace,
        events: {
          on,
          off,
          emit
        }
      }}
    >
      {children}
    </RouterContext.Provider>
  );
};

const defaultFunction = () => {};

MockRouter.defaultProps = {
  route: '/',
  pathname: '/',
  query: {},
  asPath: '/',
  back: defaultFunction,
  beforePopState: defaultFunction,
  on: defaultFunction,
  off: defaultFunction,
  emit: defaultFunction,
  prefetch: defaultFunction,
  push: defaultFunction,
  reload: defaultFunction,
  replace: defaultFunction
};

MockRouter.propTypes = {
  route: propString.isRequired,
  pathname: propString.isRequired,
  asPath: propString.isRequired,
  query: object.isRequired,
  back: func.isRequired,
  beforePopState: func.isRequired,
  on: func.isRequired,
  off: func.isRequired,
  emit: func.isRequired,
  prefetch: func.isRequired,
  push: func.isRequired,
  reload: func.isRequired,
  replace: func.isRequired
};

export default MockRouter;
