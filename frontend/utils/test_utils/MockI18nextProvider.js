import React from "react";
import i18n from "i18next";
import { I18nextProvider } from "react-i18next";

const MockI18nextProvider = ({ children }) => {
  return <I18nextProvider i18n={i18n}>{children}</I18nextProvider>;
};

export default MockI18nextProvider;
