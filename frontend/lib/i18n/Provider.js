import React from "react";
import Context from "./context";
import plTranslation from "./translations/pl";
import { node, string as propString } from "prop-types";

const Provider = ({ children, locale }) => {
  const matchTranslation = () => {
    switch (locale.toLowerCase()) {
      case "pl":
        return plTranslation;
      default:
        throw new Error(`Unknow locale: ${locale}`);
    }
  };

  return (
    <Context.Provider value={matchTranslation()}>{children}</Context.Provider>
  );
};

Provider.defaultProps = {
  locale: "pl"
};

Provider.propTypes = {
  children: node.isRequired,
  locale: propString.isRequired
};

export default Provider;
