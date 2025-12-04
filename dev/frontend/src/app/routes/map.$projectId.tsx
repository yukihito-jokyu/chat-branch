import { createFileRoute } from "@tanstack/react-router";
import { ProjectMapFlow } from "@/features/map/components/ProjectMapFlow";

export const Route = createFileRoute("/map/$projectId")({
  component: MapPage,
});

function MapPage() {
  const { projectId } = Route.useParams();

  return (
    <div className="h-screen w-screen">
      <ProjectMapFlow projectId={projectId} />
    </div>
  );
}
