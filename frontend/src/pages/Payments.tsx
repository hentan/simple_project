import { useCallback, useEffect, useMemo, useState } from "react";
import type { Expense, Payment } from "../types";
import {
  addPayment,
  deletePayment,
  getPayments,
  getExpensesWithBalance,
  updatePayment,
  type PaymentInput
} from "../api";
import ExpenseTable from "../components/ExpenseTable";
import PaymentForm from "../components/PaymentForm";
import PaymentTable from "../components/PaymentTable";

export default function Payments() {
  const [payments, setPayments] = useState<Payment[]>([]);
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [balance, setBalance] = useState<number>(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [balanceError, setBalanceError] = useState<string | null>(null);
  const [editing, setEditing] = useState<Payment | null>(null);

  const total = useMemo(() => payments.reduce((acc, x) => acc + (x.summ ?? 0), 0), [payments]);

  const refresh = useCallback(async () => {
    setError(null);
    setBalanceError(null);
    setLoading(true);
    try {
      const data = await getPayments();
      setPayments(data);
    } catch (err: any) {
      setError(err?.message ?? "Ошибка загрузки платежей");
    } finally {
      setLoading(false);
    }

    try {
      const wb = await getExpensesWithBalance();
      setExpenses(wb.expenses);
      setBalance(wb.balance);
    } catch (err: any) {
      setBalanceError(err?.message ?? "Ошибка загрузки ExpenseWithBalance");
      setExpenses([]);
      setBalance(0);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const handleCreate = useCallback(async (payload: PaymentInput) => {
    await addPayment(payload);
    await refresh();
  }, [refresh]);

  const handleUpdate = useCallback(async (payload: PaymentInput & { id: number }) => {
    await updatePayment(payload);
    setEditing(null);
    await refresh();
  }, [refresh]);

  const handleDelete = useCallback(async (id: number) => {
    const ok = window.confirm(`Удалить платеж #${id}?`);
    if (!ok) return;
    try {
      await deletePayment(id);
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
          <h1 className="h1">Сданные деньги</h1>
          <div className="muted">
            Основная таблица: <code>/payments</code>. Ниже дополнительно выводится ответ <code>ExpenseWithBalance</code> из <code>/expenses</code>.
          </div>
        </div>
        <div className="kpiRow">
          <div className="kpi">
            <div className="kpiLabel">Сумма платежей</div>
            <div className="kpiValue mono">{total}</div>
          </div>
          <div className="kpi">
            <div className="kpiLabel">Баланс</div>
            <div className="kpiValue mono">{balance}</div>
          </div>
        </div>
      </header>

      <main className="grid">
        <PaymentForm
          editing={editing}
          onCreate={handleCreate}
          onUpdate={handleUpdate}
          onCancelEdit={() => setEditing(null)}
        />

        <section>
          {error ? <div className="error">Ошибка: {error}</div> : null}
          {loading ? <div className="muted">Загрузка...</div> : null}
          <PaymentTable payments={payments} onEdit={setEditing} onDelete={handleDelete} />

          {balanceError ? <div className="error" style={{ marginTop: 12 }}>Ошибка: {balanceError}</div> : null}
          <div style={{ marginTop: 12 }}>
            <ExpenseTable title="ExpenseWithBalance: потраченные деньги" expenses={expenses} />
          </div>
        </section>
      </main>
    </>
  );
}
