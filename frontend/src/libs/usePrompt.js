import { useEffect, useRef } from 'react';

export default (message = '') => {
  const callback = useRef(function(event) {
    event.returnValue = message;
  });
  useEffect(() => {
    callback.current = function(event) {
      event.returnValue = message;
    };
  }, [message]);
  useEffect(() => {
    window.addEventListener('beforeunload', callback.current);
    return () => {
      window.removeEventListener('beforeunload', callback.current);
    };
  });
};
