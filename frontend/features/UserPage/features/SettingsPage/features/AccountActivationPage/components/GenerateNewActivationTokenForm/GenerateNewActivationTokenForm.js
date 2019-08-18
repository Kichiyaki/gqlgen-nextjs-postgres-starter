import React from "react";
import { func } from "prop-types";
import { Mutation } from "react-apollo";
import Button from "@material-ui/core/Button";
import { makeStyles } from "@material-ui/core/styles";
import { GENERATE_NEW_ACTIVATION_TOKEN_FOR_CURRENT_USER } from "../../mutations";
import { showErrorMessage, showSuccessMessage } from "@services/toastify";
import constants from "../../constants";

const useStyles = makeStyles(theme => ({
  button: {
    marginBottom: theme.spacing(2),
    marginTop: theme.spacing(2)
  }
}));

const GenerateNewActivationTokenForm = ({ t }) => {
  const classes = useStyles();

  const handleClick = generateNewActivationToken => async () => {
    try {
      await generateNewActivationToken();
      showSuccessMessage(t("generateNewActivationTokenForm.success"));
    } catch (error) {
      if (error.graphQLErrors && error.graphQLErrors[0]) {
        showErrorMessage(error.graphQLErrors[0].message);
      } else {
        showErrorMessage(t("generateNewActivationTokenForm.errors.default"));
      }
    }
  };

  return (
    <Mutation mutation={GENERATE_NEW_ACTIVATION_TOKEN_FOR_CURRENT_USER}>
      {(generateNewActivationToken, { loading }) => {
        return (
          <Button
            type="submit"
            fullWidth
            variant="contained"
            color="primary"
            onClick={handleClick(generateNewActivationToken)}
            disabled={loading}
            className={classes.button}
            data-testid={constants.GENERATE_NEW_ACTIVATION_TOKEN_BUTTON_TESTID}
          >
            {t("generateNewActivationTokenForm.submitButton")}
          </Button>
        );
      }}
    </Mutation>
  );
};

GenerateNewActivationTokenForm.propTypes = {
  t: func.isRequired
};

export default GenerateNewActivationTokenForm;
