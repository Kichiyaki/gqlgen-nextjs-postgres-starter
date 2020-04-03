import { isObject, isEmpty, isString } from 'lodash';
import { isDynamicRoute } from 'next/dist/next-server/lib/router/utils/is-dynamic';

const TEST_REGEX = /\/\[[^/]+?\](?=\/|$)/g;

export default (route, params) => {
  if (
    !isString(route) ||
    !isDynamicRoute(route) ||
    !isObject(params) ||
    isEmpty(params)
  ) {
    return route;
  }

  const paramsToReplace = route
    .match(TEST_REGEX)
    .map(param => {
      return {
        param,
        value: params[param.replace('/[', '').replace(']', '')] || param
      };
    })
    .filter(({ param, value }) => param != value);

  let asPath = route;
  paramsToReplace.forEach(({ param, value }) => {
    asPath = asPath.replace(param, `/${value}`);
  });
  return asPath;
};
