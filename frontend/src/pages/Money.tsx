import { useCallback, useEffect, useMemo, useState } from "react";
import type { Expense } from "../types";
import { addExpense, deleteExpense, getExpensesWithBalance, updateExpense, type ExpenseInput } from "../api";
import ExpenseForm from "../components/ExpenseForm";
import ExpenseTable from "../components/ExpenseTable";

type Props = {
  title?: string;
};

export default function Money({ title = "Сданные деньги" }: Props) {
  const [items, setItems] = useState<Expense[]>([]);
  const [balance, setBalance] = useState<number>(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [editing, setEditing] = useState<Expense | null>(null);

  const total = useMemo(() => items.reduce((acc, x) => acc + (x.summ ?? 0), 0), [items]);

  const refresh = useCallback(async () => {
    setError(null);
    setLoading(true);
    try {
      const data = await getExpensesWithBalance();
      setItems(data.expenses);
      setBalance(data.balance);
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
    const ok = window.confirm(`Удалить запись #${id}?`);
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
    <>
      <header className="header">
        <div>
          <h1 className="h1">{title}</h1>
          <div className="muted">
            Данные загружаются из <code>/expenses</code> (в ответе также приходит <code>Balance</code>),
            и изменяются через <code>/expenses</code> (POST/PUT/DELETE).
          </div>
        </div>
        <div className="kpiRow">
          <div className="kpi">
            <div className="kpiLabel">Итого (сдано)</div>
            <div className="kpiValue mono">{total}</div>
          </div>
          <div className="kpi">
            <div className="kpiLabel">Баланс</div>
            <div className="kpiValue mono">{balance}</div>
          </div>
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
          <ExpenseTable expenses={items} onEdit={setEditing} onDelete={handleDelete} />
        </section>
      </main>
    </>
  );
}
