import React, { useContext } from "react";
import { makeStyles } from "@material-ui/core/styles";
import AppBar from "@material-ui/core/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from "@material-ui/core/Typography";
import Button from "@material-ui/core/Button";
import I18NContext from "@lib/i18n/context";

const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1
  },
  title: {
    flexGrow: 1
  }
}));

const AppHeaderCmp = () => {
  const classes = useStyles();
  const { APPLICATION } = useContext(I18NContext);

  return (
    <div className={classes.root}>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" className={classes.title}>
            {APPLICATION.name}
          </Typography>
          <div>
            <Button color="inherit">Informacje ogólne</Button>
            <Button color="inherit">Regulamin</Button>
            <Button color="inherit">O autorze</Button>
          </div>
        </Toolbar>
      </AppBar>
    </div>
  );
};

export default AppHeaderCmp;
