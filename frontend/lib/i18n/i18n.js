const NextI18Next = require("next-i18next").default;

let NextI18NextInstance;
if (process.browser || process.env.NODE_ENV === "test") {
  NextI18NextInstance = new NextI18Next({
    defaultLanguage: "pl",
    otherLanguages: ["en"]
  });
} else {
  const path = require("path");
  const walkSync = require("../../utils/walkSync");

  NextI18NextInstance = new NextI18Next({
    defaultLanguage: "pl",
    otherLanguages: ["en"],
    ns: walkSync(path.join(process.cwd(), `/static/locales/pl`)).map(val => {
      return val
        .replace(path.join(process.cwd(), `/static/locales/pl/`), "")
        .replace(/\\/g, "/")
        .replace(".json", "");
    })
  });
}

module.exports = NextI18NextInstance;
