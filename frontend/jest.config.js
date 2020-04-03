module.exports = {
  setupFilesAfterEnv: ['jest-date-mock', '<rootDir>/src/setupTests'],
  coverageDirectory: '<rootDir>/reports/coverage',
  cache: false,
  testPathIgnorePatterns: [
    '<rootDir>/.next/',
    '<rootDir>/node_modules/',
    '<rootDir>/pages/'
  ],
  moduleNameMapper: {
    '\\.(css|less|scss)$': 'identity-obj-proxy'
  }
};
