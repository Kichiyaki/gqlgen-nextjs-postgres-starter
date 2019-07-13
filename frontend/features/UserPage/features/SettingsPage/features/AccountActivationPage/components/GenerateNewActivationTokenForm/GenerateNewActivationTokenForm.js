import React from "react";
import { object } from "prop-types";
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

const GenerateNewActivationTokenForm = ({ translations }) => {
  const {
    USER_PAGE: {
      SETTINGS_PAGE: {
        ACCOUNT_ACTIVATION_PAGE: {
          generateNewActivationTokenForm: { submitButton, success, errors }
        }
      }
    }
  } = translations;
  const classes = useStyles();

  const handleClick = generateNewActivationToken => async () => {
    try {
      await generateNewActivationToken();
      showSuccessMessage(success);
    } catch (error) {
      if (error.graphQLErrors && error.graphQLErrors[0]) {
        showErrorMessage(error.graphQLErrors[0].message);
      } else {
        showErrorMessage(errors.default);
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
            {submitButton}
          </Button>
        );
      }}
    </Mutation>
  );
};

GenerateNewActivationTokenForm.propTypes = {
  translations: object.isRequired
};

export default GenerateNewActivationTokenForm;
