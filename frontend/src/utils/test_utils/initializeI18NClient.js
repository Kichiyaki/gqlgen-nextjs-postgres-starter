import path from 'path';
import fs from 'fs';
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import walkSync from '../walkSync';

const getResources = () => {
  const result = {};
  const p = path.join(process.cwd(), `/public/locales/en`);
  walkSync(p).map(file => {
    if (file.includes('json')) {
      const json = fs.readFileSync(file, { encoding: 'utf-8' });
      const parsed = JSON.parse(json);
      const namespace = file
        .replace(p + '\\', '')
        .replace(/\\/g, '/')
        .replace('.json', '');
      result[namespace] = parsed;
    }
  });
  return result;
};

export default async ({
  resources = {
    en: getResources()
  },
  ns = ['common'],
  defaultNS = 'common',
  lng = 'en'
} = {}) => {
  if (!i18n.isInitialized) {
    await i18n.use(initReactI18next).init({
      fallbackLng: 'en',
      lng,
      resources,
      // have a common namespace used around the full app
      ns,
      defaultNS,
      interpolation: {
        escapeValue: false // not needed for react!!
      },
      react: {
        wait: false,
        nsMode: 'fallback'
      }
    });
  }
  return i18n;
};
