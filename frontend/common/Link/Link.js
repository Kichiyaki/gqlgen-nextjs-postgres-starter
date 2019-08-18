import React from "react";
import { node, object } from "prop-types";
import NextLink from "next/link";
import MaterialUILink from "@material-ui/core/Link";

export default function Link({ children, linkProps, ...props }) {
  return (
    <NextLink {...props}>
      <MaterialUILink {...linkProps}>{children}</MaterialUILink>
    </NextLink>
  );
}

Link.defaultProps = {
  linkProps: {}
};

Link.propTypes = {
  children: node.isRequired,
  linkProps: object.isRequired
};
