import React, { useState } from "react";
import { useMutation } from "react-apollo";
import isEmail from "validator/lib/isEmail";
import { bool, func } from "prop-types";
import Button from "@material-ui/core/Button";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogTitle from "@material-ui/core/DialogTitle";

import { showErrorMessage, showSuccessMessage } from "@services/toastify";
import { GENERATE_NEW_RESET_PASSWORD_TOKEN_MUTATION } from "../../mutations";
import TextField from "@common/Form/TextField/TextField";

const ResetPasswordModal = ({ open, handleClose, t }) => {
  const [email, setEmail] = useState("");
  const [generateNewResetPasswordToken] = useMutation(
    GENERATE_NEW_RESET_PASSWORD_TOKEN_MUTATION
  );

  const handleChange = e => {
    setEmail(e.target.value);
  };

  const handleSubmit = async () => {
    if (email.trim().length === 0)
      return showErrorMessage(
        t("resetPasswordModal.errors.validation.mustProvideEmail")
      );
    if (!isEmail(email)) {
      return showErrorMessage(
        t("resetPasswordModal.errors.validation.invalidEmail", { email })
      );
    }

    try {
      await generateNewResetPasswordToken({ variables: { email } });
      showSuccessMessage(t("resetPasswordModal.success"));
      handleClose();
    } catch (error) {
      if (error.graphQLErrors && error.graphQLErrors[0]) {
        showErrorMessage(error.graphQLErrors[0].message);
      } else {
        showErrorMessage(t("resetPasswordModal.errors.default"));
      }
    }
  };

  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogTitle>{t("resetPasswordModal.title")}</DialogTitle>
      <DialogContent>
        <TextField
          autoFocus
          margin="dense"
          id="name"
          label={t("resetPasswordModal.inputLabels.email")}
          type="email"
          fullWidth
          value={email}
          onChange={handleChange}
          required
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} color="primary">
          {t("resetPasswordModal.cancel")}
        </Button>
        <Button onClick={handleSubmit} color="primary">
          {t("resetPasswordModal.submitButton")}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

ResetPasswordModal.propTypes = {
  open: bool.isRequired,
  handleClose: func.isRequired,
  t: func.isRequired
};

export default ResetPasswordModal;
