import React from "react";
import { makeStyles } from "@material-ui/core/styles";
import Typography from "@material-ui/core/Typography";
import Paper from "@material-ui/core/Paper";
import Grid from "@material-ui/core/Grid";
import Divider from "@material-ui/core/Divider";
import AppLayout from "@common/AppLayout/AppLayout";
import Navigation from "../Navigation/Navigation";

const useStyles = makeStyles(theme => ({
  hide: {
    [theme.breakpoints.down("sm")]: {
      display: "none"
    }
  },
  contentContainer: {
    padding: theme.spacing(2),
    height: "100%"
  }
}));

const PageLayout = ({ children, title }) => {
  const classes = useStyles();
  return (
    <AppLayout>
      <Grid item md={2} className={classes.hide} />
      <Grid item xs={4} md={2}>
        <Navigation />
      </Grid>
      <Grid item xs={8} md={6}>
        <Paper className={classes.contentContainer}>
          <Typography align="center" variant="h5" component="h2">
            {title}
          </Typography>
          <Divider />
          <div>{children}</div>
        </Paper>
      </Grid>
    </AppLayout>
  );
};

export default PageLayout;
