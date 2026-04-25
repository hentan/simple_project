import type { ExpenseArchive } from "../types";

type Props = {
  expenses: ExpenseArchive[];
};

function formatDate(iso: string): string {
  if (!iso) return "";
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return "";
  const y = d.getUTCFullYear();
  const m = String(d.getUTCMonth() + 1).padStart(2, "0");
  const day = String(d.getUTCDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}

function formatDateTime(iso: string): string {
  if (!iso) return "";
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return "";
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  const hh = String(d.getHours()).padStart(2, "0");
  const mm = String(d.getMinutes()).padStart(2, "0");
  return `${y}-${m}-${day} ${hh}:${mm}`;
}

export default function ExpenseArchiveTable({ expenses }: Props) {
  return (
    <section className="card">
      <div className="cardHeader">
        <h2 className="h2">Архив расходов</h2>
        <div className="muted">Всего: {expenses.length}</div>
      </div>

      <div className="tableWrap">
        <table className="table">
          <thead>
            <tr>
              <th>Архивировано</th>
              <th>Операция</th>
              <th>ID</th>
              <th>Дата расхода</th>
              <th>Назначение</th>
              <th>Pupil ID</th>
              <th>Фамилия</th>
              <th>Сумма</th>
            </tr>
          </thead>
          <tbody>
            {expenses.map((e, index) => (
              <tr key={`${e.archived_at}-${e.op}-${e.id}-${index}`}>
                <td className="mono">{formatDateTime(e.archived_at)}</td>
                <td className="mono">{e.op}</td>
                <td className="mono">{e.id}</td>
                <td className="mono">{formatDate(e.date)}</td>
                <td>{e.gift_for}</td>
                <td className="mono">{e.pupil_id}</td>
                <td>{e.surname}</td>
                <td className="mono">{e.summ}</td>
              </tr>
            ))}
            {expenses.length === 0 ? (
              <tr>
                <td colSpan={8} className="empty">
                  Нет архивных данных.
                </td>
              </tr>
            ) : null}
          </tbody>
        </table>
      </div>
    </section>
  );
}
