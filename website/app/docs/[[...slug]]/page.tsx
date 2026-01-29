import { notFound } from "next/navigation";
import { getDocBySlug, getAllDocSlugs, getFlatDocList } from "../../lib/docs";
import DocContent from "../components/DocContent";
import TableOfContents from "../components/TableOfContents";
import Breadcrumbs from "../components/Breadcrumbs";
import PrevNextLinks from "../components/PrevNextLinks";

interface PageProps {
  params: Promise<{ slug?: string[] }>;
}

export async function generateStaticParams() {
  const slugs = getAllDocSlugs();
  return slugs.map((slug) => ({
    slug: slug.length > 0 ? slug : undefined,
  }));
}

export async function generateMetadata({ params }: PageProps) {
  const { slug } = await params;
  const doc = await getDocBySlug(slug || []);
  if (!doc) return { title: "Not Found" };
  return { title: `${doc.title} — sreq docs` };
}

export default async function DocPage({ params }: PageProps) {
  const { slug } = await params;
  const resolvedSlug = slug || [];
  const doc = await getDocBySlug(resolvedSlug);

  if (!doc) notFound();

  const flatList = getFlatDocList();
  const currentPath =
    resolvedSlug.length > 0
      ? `/docs/${resolvedSlug.join("/")}`
      : "/docs";
  const currentIndex = flatList.findIndex((l) => l.href === currentPath);
  const prev = currentIndex > 0 ? flatList[currentIndex - 1] : null;
  const next =
    currentIndex >= 0 && currentIndex < flatList.length - 1
      ? flatList[currentIndex + 1]
      : null;

  return (
    <>
      <Breadcrumbs slug={resolvedSlug} />
      <DocContent html={doc.html} />
      <TableOfContents items={doc.toc} />
      <PrevNextLinks prev={prev} next={next} />
    </>
  );
}
