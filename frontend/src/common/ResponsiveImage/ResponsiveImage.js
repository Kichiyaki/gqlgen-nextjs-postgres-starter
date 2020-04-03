import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import classnames from 'classnames';

const useStyles = makeStyles(() => {
  return {
    image: {
      display: 'block',
      width: '100%',
      height: 'auto'
    }
  };
});

const ResponsiveImage = props => {
  const classes = useStyles();
  return (
    <img {...props} className={classnames(classes.image, props.className)} />
  );
};

export default ResponsiveImage;
