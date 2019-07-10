import React from "react";
import { object } from "prop-types";
import Button from "@material-ui/core/Button";
import { makeStyles } from "@material-ui/core/styles";
import formikPropTypes from "@utils/formikPropTypes";
import constants from "../../constants";

import TextField from "@common/Form/TextField/TextField";

const useStyles = makeStyles(theme => ({
  form: {
    width: "100%", // Fix IE 11 issue.
    marginTop: theme.spacing(1)
  },
  button: {
    marginBottom: theme.spacing(2)
  }
}));

const LoginFormCmp = ({
  errors,
  values,
  touched,
  handleSubmit,
  handleChange,
  handleBlur,
  translations
}) => {
  const classes = useStyles();
  return (
    <form
      onSubmit={handleSubmit}
      data-testid={constants.LOGIN_FORM}
      className={classes.form}
      noValidate
    >
      <TextField
        variant="outlined"
        margin="normal"
        required
        fullWidth
        id={constants.LOGIN}
        label={translations.REGISTER_PAGE.registerForm.inputLabels.login}
        name={constants.LOGIN}
        autoComplete={constants.LOGIN}
        autoFocus
        value={values.login}
        onBlur={handleBlur}
        onChange={handleChange}
        error={touched.login && !!errors.login}
        helperText={touched.login && errors.login}
      />
      <TextField
        variant="outlined"
        margin="normal"
        required
        fullWidth
        name={constants.PASSWORD}
        label={translations.REGISTER_PAGE.registerForm.inputLabels.password}
        type={constants.PASSWORD}
        id={constants.PASSWORD}
        autoComplete={constants.CURRENT_PASSWORD}
        value={values.password}
        onBlur={handleBlur}
        onChange={handleChange}
        error={touched.password && !!errors.password}
        helperText={touched.password && errors.password}
      />
      <Button
        type="submit"
        fullWidth
        variant="contained"
        color="primary"
        className={classes.button}
      >
        {translations.REGISTER_PAGE.loginForm.submitButton}
      </Button>
    </form>
  );
};

LoginFormCmp.propTypes = {
  translations: object.isRequired,
  ...formikPropTypes
};

export default LoginFormCmp;
