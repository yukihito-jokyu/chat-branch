import { usePosts } from "../hooks/usePosts";

export const PostList = () => {
  const { data: posts, isLoading, error } = usePosts();

  if (isLoading) return <div className="p-4">Loading...</div>;
  if (error)
    return <div className="p-4 text-red-500">Error: {error.message}</div>;

  return (
    <div className="p-4">
      <h2 className="text-2xl font-bold mb-4">Posts</h2>
      <div className="grid gap-4">
        {posts?.map((post) => (
          <div key={post.id} className="border p-4 rounded shadow">
            <h3 className="font-bold text-lg">{post.title}</h3>
            <p className="text-gray-600">{post.body}</p>
          </div>
        ))}
      </div>
    </div>
  );
};
