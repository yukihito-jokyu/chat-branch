import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import LanguageDetector from "i18next-browser-languagedetector";

import authJa from "../locales/auth/ja.json";
import authEn from "../locales/auth/en.json";
import chatJa from "../locales/chat/ja.json";
import chatEn from "../locales/chat/en.json";

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources: {
      ja: {
        auth: authJa,
        chat: chatJa,
      },
      en: {
        auth: authEn,
        chat: chatEn,
      },
    },
    fallbackLng: "en",
    lng: "ja",
    ns: ["auth", "chat"],
    defaultNS: "auth",
    interpolation: {
      escapeValue: false,
    },
  });

export default i18n;
