import { useQuery } from "@tanstack/react-query";
import { fetchPosts, type Post } from "../api/posts";

export const usePosts = () => {
  return useQuery<Post[], Error>({
    queryKey: ["posts"],
    queryFn: fetchPosts,
  });
};
