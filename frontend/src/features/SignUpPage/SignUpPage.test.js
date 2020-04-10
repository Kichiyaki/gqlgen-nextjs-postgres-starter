import { render } from '@testing-library/react';
import i18n from 'i18next';
import { SIGN_IN_PAGE, SIGN_UP_PAGE } from '@config/namespaces';
import * as routes from '@config/routes';
import TestLayout from '@utils/test_utils/TestLayout';
import createApolloClient from '@utils/test_utils/createApolloClient';
import SignUpPage from './SignUpPage';

const t1 = i18n.getFixedT(null, SIGN_UP_PAGE);
const t2 = i18n.getFixedT(null, SIGN_IN_PAGE);

const renderPage = (pathname) => {
  return render(
    <TestLayout client={createApolloClient()}>
      <SignUpPage pathname={pathname} />
    </TestLayout>
  );
};

describe('features > SignUpPage', () => {
  test('should correctly render signup page', async () => {
    const { queryByText, asFragment } = renderPage(routes.SIGN_UP_PAGE);
    expect(queryByText(t1('signUpForm.title'))).toBeInTheDocument();
    expect(queryByText(t1('signUpForm.inputLabel.login'))).toBeInTheDocument();
    expect(queryByText(t1('signUpForm.inputLabel.email'))).toBeInTheDocument();
    expect(
      queryByText(t1('signUpForm.inputLabel.password'))
    ).toBeInTheDocument();
    expect(
      queryByText(t1('signUpForm.inputLabel.confirmPassword'))
    ).toBeInTheDocument();
    expect(
      queryByText(t1('signUpForm.alreadyHaveAnAccount'))
    ).toBeInTheDocument();
    expect(queryByText(t1('signUpForm.submitButton'))).toBeInTheDocument();
    expect(asFragment()).toMatchSnapshot();
  });

  test('should correctly render signin page', async () => {
    const { queryByText, asFragment } = renderPage(routes.SIGN_IN_PAGE);
    expect(queryByText(t2('signInForm.title'))).toBeInTheDocument();
    expect(queryByText(t2('signInForm.inputLabel.login'))).toBeInTheDocument();
    expect(
      queryByText(t2('signInForm.inputLabel.password'))
    ).toBeInTheDocument();
    expect(queryByText(t2('signInForm.dontHaveAnAccount'))).toBeInTheDocument();
    expect(queryByText(t2('signInForm.forgotPassword'))).toBeInTheDocument();
    expect(queryByText(t2('signInForm.submitButton'))).toBeInTheDocument();
    expect(asFragment()).toMatchSnapshot();
  });
});
