import { apiClient } from "@/lib/api-client";
import type { GetProjectTreeResponse } from "../types";

export const getProjectTree = async (
  projectUuid: string
): Promise<GetProjectTreeResponse> => {
  return apiClient.get(`/api/projects/${projectUuid}/tree`);
};
