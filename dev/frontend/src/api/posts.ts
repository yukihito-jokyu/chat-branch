export interface Post {
  userId: number;
  id: number;
  title: string;
  body: string;
}

export const fetchPosts = async (): Promise<Post[]> => {
  const response = await fetch("https://jsonplaceholder.typicode.com/posts");
  if (!response.ok) {
    throw new Error("Network response was not ok");
  }
  return response.json();
};
