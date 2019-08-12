import React from "react";
import { func } from "prop-types";
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
  t,
  isSubmitting
}) => {
  const classes = useStyles();
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
        label={t("activateAccountForm.inputLabels.token")}
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
        {t("activateAccountForm.submitButton")}
      </Button>
    </form>
  );
};

ActivateAccountFormCmp.propTypes = {
  t: func.isRequired,
  ...formikPropTypes
};

export default ActivateAccountFormCmp;
