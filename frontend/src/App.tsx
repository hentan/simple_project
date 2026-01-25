import { useCallback, useEffect, useMemo, useState } from "react";
import type { Expense } from "./types";
import { addExpense, deleteExpense, getExpenses, updateExpense, type ExpenseInput } from "./api";
import ExpenseForm from "./components/ExpenseForm";
import ExpenseTable from "./components/ExpenseTable";

export default function App() {
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [editing, setEditing] = useState<Expense | null>(null);

  const total = useMemo(() => expenses.reduce((acc, x) => acc + (x.summ ?? 0), 0), [expenses]);

  const refresh = useCallback(async () => {
    setError(null);
    setLoading(true);
    try {
      const data = await getExpenses();
      setExpenses(data);
    } catch (err: any) {
      setError(err?.message ?? "Ошибка загрузки");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const handleCreate = useCallback(async (payload: ExpenseInput) => {
    await addExpense(payload);
    await refresh();
  }, [refresh]);

  const handleUpdate = useCallback(async (payload: ExpenseInput & { id: number }) => {
    await updateExpense(payload);
    setEditing(null);
    await refresh();
  }, [refresh]);

  const handleDelete = useCallback(async (id: number) => {
    const ok = window.confirm(`Удалить расход #${id}?`);
    if (!ok) return;
    try {
      await deleteExpense(id);
      if (editing?.id === id) setEditing(null);
      await refresh();
    } catch (err: any) {
      alert(err?.message ?? "Ошибка удаления");
    }
  }, [refresh, editing]);

  return (
    <div className="page">
      <header className="header">
        <div>
          <h1 className="h1">Учёт расходов</h1>
          <div className="muted">
            SPA (React + Vite) в Docker. Встроенный Nginx проксирует <code>/api/*</code> на бэкенд, чтобы не требовать CORS.
          </div>
        </div>
        <div className="kpi">
          <div className="kpiLabel">Сумма</div>
          <div className="kpiValue mono">{total}</div>
        </div>
      </header>

      <main className="grid">
        <ExpenseForm
          editing={editing}
          onCreate={handleCreate}
          onUpdate={handleUpdate}
          onCancelEdit={() => setEditing(null)}
        />

        <section>
          {error ? <div className="error">Ошибка: {error}</div> : null}
          {loading ? <div className="muted">Загрузка...</div> : null}
          <ExpenseTable expenses={expenses} onEdit={setEditing} onDelete={handleDelete} />
        </section>
      </main>

      <footer className="footer">
        <div className="muted">
          Чтобы изменить базовый путь API без прокси, установите <code>VITE_API_BASE</code> (например, <code>http://localhost:8080</code>) и включите CORS на бэкенде.
        </div>
      </footer>
    </div>
  );
}
