import React, { useMemo } from 'react';
import { string, object, node, func } from 'prop-types';

import formatDynamicRouteAsPath from '@utils/formatDynamicRouteAsPath';

import NextLink from 'next/link';
import { Link } from '@material-ui/core';

export default function MyLink({ href, params, children, onClick, ...rest }) {
  const as = useMemo(() => {
    return formatDynamicRouteAsPath(href, params);
  }, [href, params]);

  return (
    <NextLink href={href} as={as}>
      <Link color="secondary" underline="hover" href={as} {...rest}>
        {children}
      </Link>
    </NextLink>
  );
}

MyLink.defaultProps = {
  href: '',
  params: {}
};

MyLink.propTypes = {
  href: string.isRequired,
  params: object,
  children: node,
  onClick: func
};
