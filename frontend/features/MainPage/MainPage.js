import React from "react";
import AppLayout from "@common/AppLayout/AppLayout";
import withCurrentUser from "../../hocs/withCurrentUser";

const MainPage = () => {
  return (
    <AppLayout>
      <p>Kicha to kozak</p>
    </AppLayout>
  );
};

export default withCurrentUser(MainPage);
