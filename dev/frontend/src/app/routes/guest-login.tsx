import { createFileRoute } from "@tanstack/react-router";
import { GuestLoginPageContent } from "@/features/auth/components/GuestLoginPageContent";

export const Route = createFileRoute("/guest-login")({
  component: GuestLoginPageContent,
});
