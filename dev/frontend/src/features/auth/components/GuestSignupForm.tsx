import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { signupGuest, type SignupResponse } from "../api/guest-auth";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Copy } from "lucide-react";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";

export const GuestSignupForm = () => {
  const { t } = useTranslation("auth");
  const [isOpen, setIsOpen] = useState(false);
  const [userData, setUserData] = useState<SignupResponse["user"] | null>(null);

  const mutation = useMutation({
    mutationFn: signupGuest,
    onSuccess: (data) => {
      setUserData(data.user);
      setIsOpen(true);
    },
    onError: () => {
      toast.error(t("guestSignup.errorMessage"));
    },
  });

  const handleCopyId = () => {
    if (userData?.uuid) {
      navigator.clipboard.writeText(userData.uuid);
      toast.success(t("guestSignup.dialog.copySuccess"));
    }
  };

  return (
    <>
      <Card>
        <CardHeader>
          <CardTitle>{t("guestSignup.title")}</CardTitle>
          <CardDescription>{t("guestSignup.description")}</CardDescription>
        </CardHeader>
        <CardContent>
          <Button
            className="w-full"
            onClick={() => mutation.mutate()}
            disabled={mutation.isPending}
          >
            {mutation.isPending
              ? t("guestSignup.submittingButton")
              : t("guestSignup.submitButton")}
          </Button>
        </CardContent>
      </Card>

      <Dialog open={isOpen} onOpenChange={setIsOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("guestSignup.dialog.title")}</DialogTitle>
            <DialogDescription>
              {t("guestSignup.dialog.description")}
            </DialogDescription>
          </DialogHeader>
          {userData && (
            <div className="space-y-4">
              <div className="p-4 bg-muted rounded-lg space-y-4">
                <div className="space-y-2">
                  <Label>{t("guestSignup.dialog.userNameLabel")}</Label>
                  <div className="font-medium px-1">{userData.name}</div>
                </div>
                <div className="space-y-2">
                  <Label>{t("guestSignup.dialog.userIdLabel")}</Label>
                  <div className="flex items-center gap-2">
                    <Input
                      readOnly
                      value={userData.uuid}
                      className="font-mono bg-background"
                    />
                    <Button
                      size="icon"
                      variant="outline"
                      onClick={handleCopyId}
                    >
                      <Copy className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </div>
              <Button className="w-full" onClick={() => setIsOpen(false)}>
                {t("guestSignup.dialog.closeButton")}
              </Button>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </>
  );
};
