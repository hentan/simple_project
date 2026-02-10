export type HomeLink = {
  title: string;
  description: string;
  href: string;
};

type Props = {
  links: HomeLink[];
};

export default function Home({ links }: Props) {
  return (
    <section className="card">
      <div className="cardHeader">
        <h2 className="h2">Главная</h2>
        <div className="muted">Выберите раздел</div>
      </div>

      <div className="menuGrid">
        {links.map((l) => (
          <a key={l.href} className="menuCard" href={l.href}>
            <div className="menuTitle">{l.title}</div>
            <div className="muted">{l.description}</div>
          </a>
        ))}
      </div>
    </section>
  );
}
