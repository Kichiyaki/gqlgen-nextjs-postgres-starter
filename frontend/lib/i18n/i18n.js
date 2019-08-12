const NextI18Next = require("next-i18next").default;
const walkSync = require("../../utils/walkSync");

let NextI18NextInstance;
if (process.browser || process.env.NODE_ENV === "test") {
  NextI18NextInstance = new NextI18Next({
    defaultLanguage: "pl",
    otherLanguages: ["en"]
  });
} else {
  const path = require("path");

  NextI18NextInstance = new NextI18Next({
    defaultLanguage: "pl",
    otherLanguages: ["en"],
    ns: walkSync(path.join(process.cwd(), `/static/locales/en`)).map(val => {
      return val
        .replace(path.join(process.cwd(), `/static/locales/en/`), "")
        .replace(/\\/g, "/")
        .replace(".json", "");
    })
  });
}

module.exports = NextI18NextInstance;
