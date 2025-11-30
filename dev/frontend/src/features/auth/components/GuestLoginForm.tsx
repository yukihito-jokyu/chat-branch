import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import { loginGuest } from "../api/guest-auth";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";

export const GuestLoginForm = () => {
  const { t } = useTranslation("auth");
  const [userId, setUserId] = useState("");
  const navigate = useNavigate();

  const mutation = useMutation({
    mutationFn: loginGuest,
    onSuccess: (data) => {
      localStorage.setItem("token", data.token);
      toast.success(t("guestLogin.successMessage"));
      navigate({ to: "/chat" });
    },
    onError: () => {
      toast.error(t("guestLogin.errorMessage"));
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!userId.trim()) {
      toast.error(t("guestLogin.validation.userIdRequired"));
      return;
    }
    mutation.mutate({ user_id: userId });
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t("guestLogin.title")}</CardTitle>
        <CardDescription>{t("guestLogin.description")}</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="userId">{t("guestLogin.userIdLabel")}</Label>
            <Input
              id="userId"
              placeholder={t("guestLogin.userIdPlaceholder")}
              value={userId}
              onChange={(e) => setUserId(e.target.value)}
              disabled={mutation.isPending}
            />
          </div>
          <Button
            type="submit"
            className="w-full"
            disabled={mutation.isPending}
          >
            {mutation.isPending
              ? t("guestLogin.submittingButton")
              : t("guestLogin.submitButton")}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
};
