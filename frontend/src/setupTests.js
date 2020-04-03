import '@testing-library/jest-dom/extend-expect';
import 'react-act';
import { advanceTo, clear } from 'jest-date-mock';
import initializeI18NClient from '@utils/test_utils/initializeI18NClient';

beforeAll(done => {
  initializeI18NClient().then(() => {
    done();
  });
});

beforeEach(() => {
  advanceTo(0);
});

afterEach(() => {
  clear();
});

jest.setTimeout(10000);
