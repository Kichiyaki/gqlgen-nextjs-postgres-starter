import React from "react";
import Link from "next/link";
import { makeStyles } from "@material-ui/core/styles";

const useStyles = makeStyles(theme => ({
  link: {
    textDecoration: "none",
    color: "inherit",
    transition: "all .2s",
    "&:not(:last-child)": {
      marginRight: theme.spacing(2)
    },
    "&:hover": {
      transform: "translateY(-3px)",
      color: "#fff"
    }
  }
}));

const NavLink = ({ href, children }) => {
  const classes = useStyles();
  return (
    <Link href={href}>
      <a className={classes.link}>{children}</a>
    </Link>
  );
};

export default NavLink;
