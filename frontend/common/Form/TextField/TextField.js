import React from "react";
import MaterialUITextField from "@material-ui/core/TextField";
import { testID } from "./constants";

const TextField = ({ inputProps = {}, ...rest } = {}) => {
  return (
    <MaterialUITextField
      {...rest}
      inputProps={{ "data-testid": testID, ...inputProps }}
    />
  );
};

export default TextField;
