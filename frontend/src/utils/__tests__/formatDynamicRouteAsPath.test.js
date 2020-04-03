import formatDynamicRouteAsPath from '../formatDynamicRouteAsPath';

describe('utils > formatDynamicRouteAsPath', () => {
  test('should correctly fill dynamic route asPath params with values', () => {
    const dynamicRoute = '/test/[param1]/test2/[param2]/[param3]';
    const params = {
      param1: 'sdsss',
      param2: 'sdsdsds'
    };
    const asPath = formatDynamicRouteAsPath(dynamicRoute, params);
    expect(asPath).toBe(
      `/test/${params.param1}/test2/${params.param2}/[param3]`
    );
  });
});
