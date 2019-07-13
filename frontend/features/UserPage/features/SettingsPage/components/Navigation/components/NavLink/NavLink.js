import React from "react";
import Link from "next/link";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";
import { string, bool, node } from "prop-types";

const NavLink = ({ children, text, href, selected }) => {
  return (
    <Link href={href}>
      <ListItem button selected={selected} disabled={selected}>
        <ListItemIcon>{children}</ListItemIcon>
        <ListItemText primary={text} />
      </ListItem>
    </Link>
  );
};

NavLink.propTypes = {
  text: string.isRequired,
  href: string.isRequired,
  selected: bool.isRequired,
  children: node
};

export default NavLink;
