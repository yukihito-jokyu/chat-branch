import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import LanguageDetector from "i18next-browser-languagedetector";

import authJa from "../locales/auth/ja.json";
import authEn from "../locales/auth/en.json";

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources: {
      ja: {
        auth: authJa,
      },
      en: {
        auth: authEn,
      },
    },
    fallbackLng: "en",
    lng: "ja", // Default language
    ns: ["auth"],
    defaultNS: "auth",
    interpolation: {
      escapeValue: false, // React already safes from xss
    },
  });

export default i18n;
