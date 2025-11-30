import { GuestSignupForm } from "./GuestSignupForm";
import { GuestLoginForm } from "./GuestLoginForm";
import { useTranslation } from "react-i18next";

export const GuestLoginPageContent = () => {
  const { t } = useTranslation("auth");
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-background p-4">
      <div className="w-full max-w-4xl space-y-8">
        <div className="text-center space-y-2">
          <h1 className="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl">
            {t("guestLoginPage.title")}
          </h1>
          <p className="mx-auto max-w-[700px] text-muted-foreground md:text-xl">
            {t("guestLoginPage.description")}
          </p>
        </div>

        <div className="grid gap-8 md:grid-cols-2">
          <div className="space-y-4">
            <h2 className="text-2xl font-semibold text-center">
              {t("guestLoginPage.newUsersTitle")}
            </h2>
            <GuestSignupForm />
          </div>
          <div className="space-y-4">
            <h2 className="text-2xl font-semibold text-center">
              {t("guestLoginPage.existingUsersTitle")}
            </h2>
            <GuestLoginForm />
          </div>
        </div>
      </div>
    </div>
  );
};
