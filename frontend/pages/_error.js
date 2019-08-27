import React from "react";
import NextError from "next/error";
import constants from "@config/constants";

const ErrorPage = ({ statusCode }) => {
  return <NextError statusCode={statusCode} />;
};

ErrorPage.getInitialProps = ({ res, err }) => {
  const statusCode = res ? res.statusCode : err ? err.statusCode : null;
  return { statusCode, namespacesRequired: [constants.NAMESPACES.common] };
};

export default ErrorPage;
