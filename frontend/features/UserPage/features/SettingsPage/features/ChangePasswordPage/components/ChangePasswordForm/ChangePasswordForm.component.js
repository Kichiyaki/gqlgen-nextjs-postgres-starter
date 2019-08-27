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

const ChangePasswordForm = ({
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
      data-testid={constants.CHANGE_PASSWORD_FORM}
      className={classes.form}
      noValidate
    >
      <TextField
        variant="outlined"
        margin="normal"
        type="password"
        required
        fullWidth
        id={constants.CURRENT_PASSWORD}
        label={t("changePasswordForm.inputLabels.currentPassword")}
        name={constants.CURRENT_PASSWORD}
        autoComplete={constants.CURRENT_PASSWORD}
        autoFocus
        value={values.currentPassword}
        onBlur={handleBlur}
        onChange={handleChange}
        error={touched.currentPassword && !!errors.currentPassword}
        helperText={touched.currentPassword && errors.currentPassword}
      />
      <TextField
        variant="outlined"
        margin="normal"
        type="password"
        required
        fullWidth
        id={constants.NEW_PASSWORD}
        label={t("changePasswordForm.inputLabels.newPassword")}
        name={constants.NEW_PASSWORD}
        autoComplete={constants.NEW_PASSWORD}
        autoFocus
        value={values.newPassword}
        onBlur={handleBlur}
        onChange={handleChange}
        error={touched.newPassword && !!errors.newPassword}
        helperText={touched.newPassword && errors.newPassword}
      />
      <Button
        type="submit"
        fullWidth
        variant="contained"
        color="primary"
        className={classes.button}
        disabled={isSubmitting}
      >
        {t("changePasswordForm.submitButton")}
      </Button>
    </form>
  );
};

ChangePasswordForm.propTypes = {
  t: func.isRequired,
  ...formikPropTypes
};

export default ChangePasswordForm;
