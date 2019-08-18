import "@testing-library/react/cleanup-after-each";
import "jest-dom/extend-expect";
import "react-act";
import initializeI18nClient from "@utils/test_utils/initializeI18nClient";

beforeAll(done => {
  initializeI18nClient().then(() => {
    done();
  });
});

jest.setTimeout(10000);
