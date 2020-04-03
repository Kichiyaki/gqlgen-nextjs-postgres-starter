import { useState } from 'react';

export const SEVERITY = {
  INFO: 'info',
  ERROR: 'error',
  WARNING: 'warning',
  SUCCESS: 'success'
};

export const DEFAULT_ANCHOR_ORIGIN = {
  vertical: 'top',
  horizontal: 'right'
};

export const DEFAULT_AUTO_HIDE_DURATION = 5000;

export default ({
  message: initialMessage,
  autoHideDuration: initialAutoHideDuration,
  severity: initialSeverity,
  anchorOrigin: initialAnchorOrigin
} = {}) => {
  const [message, setMessage] = useState(initialMessage || '');
  const [autoHideDuration, setAutoHideDuration] = useState(
    initialAutoHideDuration || DEFAULT_AUTO_HIDE_DURATION
  );
  const [severity, setSeverity] = useState(initialSeverity || SEVERITY.INFO);
  const [anchorOrigin, setAnchorOrigin] = useState(
    initialAnchorOrigin || DEFAULT_ANCHOR_ORIGIN
  );

  const reset = () => {
    setMessage(initialMessage || '');
    setAutoHideDuration(initialAutoHideDuration || DEFAULT_AUTO_HIDE_DURATION);
    setSeverity(initialSeverity || SEVERITY.INFO);
    setAnchorOrigin(initialAnchorOrigin || DEFAULT_ANCHOR_ORIGIN);
  };

  return {
    message,
    setMessage,
    setAnchorOrigin,
    setAutoHideDuration,
    setSeverity,
    reset,
    alertProps: {
      onClose: reset,
      severity
    },
    snackbarProps: {
      onClose: reset,
      autoHideDuration,
      anchorOrigin,
      open: !!message
    }
  };
};
