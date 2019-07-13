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

const ActivateAccountFormCmp = ({
  errors,
  values,
  touched,
  handleSubmit,
  handleChange,
  handleBlur,
  translations,
  isSubmitting
}) => {
  const classes = useStyles();
  const {
    USER_PAGE: {
      SETTINGS_PAGE: {
        ACCOUNT_ACTIVATION_PAGE: {
          activateAccountForm: { inputLabels, submitButton }
        }
      }
    }
  } = translations;
  return (
    <form
      onSubmit={handleSubmit}
      data-testid={constants.ACTIVATE_ACCOUNT_FORM}
      className={classes.form}
      noValidate
    >
      <TextField
        variant="outlined"
        margin="normal"
        required
        fullWidth
        id={constants.TOKEN}
        label={inputLabels.token}
        name={constants.TOKEN}
        autoComplete={constants.TOKEN}
        autoFocus
        value={values.token}
        onBlur={handleBlur}
        onChange={handleChange}
        error={touched.token && !!errors.token}
        helperText={touched.token && errors.token}
      />
      <Button
        type="submit"
        fullWidth
        variant="contained"
        color="primary"
        className={classes.button}
        disabled={isSubmitting}
      >
        {submitButton}
      </Button>
    </form>
  );
};

ActivateAccountFormCmp.propTypes = {
  translations: object.isRequired,
  ...formikPropTypes
};

export default ActivateAccountFormCmp;
