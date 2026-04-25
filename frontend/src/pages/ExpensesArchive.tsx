import { useCallback, useEffect, useMemo, useState } from "react";
import { getExpensesArchive } from "../api";
import ExpenseArchiveTable from "../components/ExpenseArchiveTable";
import type { ExpenseArchive } from "../types";

export default function ExpensesArchive() {
  const [items, setItems] = useState<ExpenseArchive[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const total = useMemo(() => items.reduce((acc, x) => acc + (x.summ ?? 0), 0), [items]);

  const refresh = useCallback(async () => {
    setError(null);
    setLoading(true);
    try {
      const data = await getExpensesArchive();
      setItems(data);
    } catch (err: any) {
      setError(err?.message ?? "Ошибка загрузки архива расходов");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  return (
    <>
      <header className="header">
        <div>
          <h1 className="h1">Архив расходов</h1>
          <div className="muted">
            Старые версии расходов из <code>/expenses/archive</code>.
          </div>
        </div>
        <div className="kpiRow">
          <div className="kpi">
            <div className="kpiLabel">Записей в архиве</div>
            <div className="kpiValue mono">{items.length}</div>
          </div>
          <div className="kpi">
            <div className="kpiLabel">Сумма архивных записей</div>
            <div className="kpiValue mono">{total}</div>
          </div>
        </div>
      </header>

      <main>
        {error ? <div className="error">Ошибка: {error}</div> : null}
        {loading ? <div className="muted">Загрузка...</div> : null}
        <ExpenseArchiveTable expenses={items} />
      </main>
    </>
  );
}
