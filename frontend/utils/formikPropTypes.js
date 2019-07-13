import { objectOf, func, string as propString, bool } from "prop-types";

export default {
  values: objectOf(propString).isRequired,
  touched: objectOf(bool).isRequired,
  errors: objectOf(propString).isRequired,
  handleChange: func.isRequired,
  handleBlur: func.isRequired,
  handleSubmit: func.isRequired,
  setFieldError: func.isRequired,
  setFieldTouched: func.isRequired,
  setFieldValue: func.isRequired,
  setSubmitting: func.isRequired,
  setTouched: func.isRequired,
  setValues: func.isRequired,
  isSubmitting: bool.isRequired
};
