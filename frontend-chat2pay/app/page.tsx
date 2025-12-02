import { createApiClient } from "@/lib/api/client";

type Post = {
  id: number;
  userId: number;
  title: string;
  body: string;
};


const placeholderClient = createApiClient({
  baseUrl: "https://jsonplaceholder.typicode.com",
});

async function fetchPosts(): Promise<Post[]> {
  const { data } = await placeholderClient.get<Post[]>("/posts", {
    query: { _limit: 6 },
    next: { revalidate: 120 }, 
  });

  return data;
}

export default async function Home() {
  const posts = await fetchPosts();

  return (
    <div className="min-h-screen bg-zinc-50 font-sans text-black dark:bg-black dark:text-zinc-50">
      <main className="mx-auto flex max-w-5xl flex-col gap-8 px-6 py-16">
        <header className="flex flex-col gap-2">
          <p className="text-sm font-medium uppercase tracking-[0.16em] text-zinc-500 dark:text-zinc-400">
            Contoh Integrasi API
          </p>
          <h1 className="text-4xl font-semibold leading-tight">Posts dari JSONPlaceholder</h1>
          <p className="max-w-2xl text-lg text-zinc-600 dark:text-zinc-400">
            Halaman ini mengambil data dummy dari{" "}
            <span className="font-semibold text-black dark:text-zinc-100">
              https://jsonplaceholder.typicode.com/posts
            </span>{" "}
            dan merender enam postingan pertama.
          </p>
        </header>

        <section className="grid gap-4 sm:grid-cols-2">
          {posts.map((post) => (
            <article
              key={post.id}
              className="rounded-2xl border border-zinc-200 bg-white p-5 shadow-sm transition hover:-translate-y-0.5 hover:shadow-md dark:border-zinc-800 dark:bg-zinc-900"
            >
              <p className="mb-3 text-xs font-semibold uppercase tracking-[0.2em] text-amber-500">
                #{post.id} • User {post.userId}
              </p>
              <h2 className="mb-2 text-xl font-semibold leading-7">{post.title}</h2>
              <p className="text-sm leading-6 text-zinc-600 dark:text-zinc-400">{post.body}</p>
            </article>
          ))}
        </section>
      </main>
    </div>
  );
}
