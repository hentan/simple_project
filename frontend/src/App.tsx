import { useEffect, useMemo, useState } from "react";
import Home, { type HomeLink } from "./pages/Home";
import ExpensesArchive from "./pages/ExpensesArchive";
import Money from "./pages/Money";
import Payments from "./pages/Payments";
import Pupils from "./pages/Pupils";

type RouteKey = "" | "money" | "pupils" | "payments" | "expenses-archive";

function normalizeHashToRoute(hash: string): RouteKey {
  const raw = hash.replace(/^#\/?/, "").trim();
  if (raw === "money") return "money";
  if (raw === "pupils") return "pupils";
  if (raw === "payments") return "payments";
  if (raw === "expenses-archive") return "expenses-archive";
  return "";
}

function linkClass(isActive: boolean): string {
  return `navLink${isActive ? " navLinkActive" : ""}`;
}

export default function App() {
  const [route, setRoute] = useState<RouteKey>(() => normalizeHashToRoute(window.location.hash));

  useEffect(() => {
    const onHashChange = () => setRoute(normalizeHashToRoute(window.location.hash));
    window.addEventListener("hashchange", onHashChange);
    return () => window.removeEventListener("hashchange", onHashChange);
  }, []);

  const links: HomeLink[] = useMemo(() => (
    [
      {
        title: "Сданные деньги",
        description: "Учёт поступлений (то, что сдали).",
        href: "#/money"
      },
      {
        title: "Ученики",
        description: "Справочник учеников.",
        href: "#/pupils"
      },
      {
        title: "Потраченные деньги",
        description: "Сданные деньги + расчёт потраченные",
        href: "#/payments"
      },
      {
        title: "Архив расходов",
        description: "Старые версии расходов после изменения или удаления.",
        href: "#/expenses-archive"
      }
    ]
  ), []);

  const content = useMemo(() => {
    switch (route) {
      case "money":
        return <Money />;
      case "pupils":
        return <Pupils />;
      case "payments":
        return <Payments />;
      case "expenses-archive":
        return <ExpensesArchive />;
      default:
        return <Home links={links} />;
    }
  }, [route, links]);

  return (
    <div className="page">
      <div className="topNav">
        <a href="#/" className="brand">
          Учёт
        </a>
        <div className="nav">
          <a className={linkClass(route === "")} href="#/">
            Главная
          </a>
          <a className={linkClass(route === "money")} href="#/money">
            Оплаты
          </a>
          <a className={linkClass(route === "pupils")} href="#/pupils">
            Ученики
          </a>
          <a className={linkClass(route === "payments")} href="#/payments">
            Сданные деньги
          </a>
          <a className={linkClass(route === "expenses-archive")} href="#/expenses-archive">
            Архив расходов
          </a>
        </div>
      </div>

      {content}

      <footer className="footer">
        <div className="muted">
          SPA (React + Vite) в Docker. Встроенный Nginx проксирует <code>/api/*</code> на бэкенд, чтобы не требовать CORS.
          <br />
          Чтобы изменить базовый путь API без прокси, установите <code>VITE_API_BASE</code> (например, <code>http://localhost:8080</code>) и включите CORS на бэкенде.
        </div>
      </footer>
    </div>
  );
}
