import path from "path";
import fs from "fs";
import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import walkSync from "../walkSync";

const getDefaultResources = (language = "pl") => {
  const result = {};
  walkSync(path.join(process.cwd(), `/static/locales/${language}`)).map(
    file => {
      const json = fs.readFileSync(file, { encoding: "utf-8" });
      const parsed = JSON.parse(json);
      const namespace = file
        .replace(path.join(process.cwd(), `/static/locales/${language}/`), "")
        .replace(/\\/g, "/")
        .replace(".json", "");
      result[namespace] = parsed;
    }
  );
  return result;
};

export default async ({
  resources = {
    pl: getDefaultResources("pl")
  },
  ns = ["common"],
  defaultNS = "common",
  lng = "pl"
} = {}) => {
  if (!i18n.isInitialized) {
    await i18n.use(initReactI18next).init({
      fallbackLng: "pl",
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
        nsMode: "fallback"
      }
    });
  }
  return i18n;
};
