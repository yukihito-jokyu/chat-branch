import { useQuery } from "@tanstack/react-query";
import { useNavigate, useLocation } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { PlusCircle, MessageSquare } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { getProjects, getProject } from "../api/chat";

export function Sidebar() {
  const navigate = useNavigate();
  const location = useLocation();
  const { t } = useTranslation("chat");

  const { data: projects } = useQuery({
    queryKey: ["projects"],
    queryFn: getProjects,
  });

  const isChatIndex =
    location.pathname === "/chat" || location.pathname === "/chat/";

  const handleProjectClick = async (projectUuid: string) => {
    try {
      const data = await getProject(projectUuid);
      navigate({ to: "/chat/$chatId", params: { chatId: data.chat_uuid } });
    } catch (error) {
      console.error("Failed to fetch project chat:", error);
    }
  };

  return (
    <div className="flex flex-col h-full py-4">
      <div className="px-4 mb-4">
        <Button
          className="w-full justify-start gap-2"
          variant="outline"
          onClick={() => navigate({ to: "/chat" })}
          disabled={isChatIndex}
        >
          <PlusCircle className="h-4 w-4" />
          {t("new_chat")}
        </Button>
      </div>

      <ScrollArea className="flex-1 px-2">
        <div className="space-y-1">
          {projects?.map((project) => (
            <div
              key={project.uuid}
              onClick={() => handleProjectClick(project.uuid)}
              className="flex items-center gap-2 px-3 py-2 text-sm font-medium rounded-md hover:bg-accent hover:text-accent-foreground transition-colors cursor-pointer"
            >
              <MessageSquare className="h-4 w-4" />
              <span className="truncate">
                {project.title || t("default_new_chat_title")}
              </span>
            </div>
          ))}
        </div>
      </ScrollArea>
    </div>
  );
}
