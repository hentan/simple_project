import type { Expense } from "../types";

type Props = {
  expenses: Expense[];
  title?: string;
  onEdit?: (e: Expense) => void;
  onDelete?: (id: number) => void;
};

function formatDate(iso: string): string {
  const d = new Date(iso);
  // show as YYYY-MM-DD (UTC) to match input expectations
  const y = d.getUTCFullYear();
  const m = String(d.getUTCMonth() + 1).padStart(2, "0");
  const day = String(d.getUTCDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}

export default function ExpenseTable({ expenses, title = "Потраченные деньги", onEdit, onDelete }: Props) {
  const canEdit = typeof onEdit === "function" && typeof onDelete === "function";
  return (
    <section className="card">
      <div className="cardHeader">
        <h2 className="h2">{title}</h2>
        <div className="muted">Всего: {expenses.length}</div>
      </div>

      <div className="tableWrap">
        <table className="table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Дата</th>
              <th>Назначение</th>
              <th>Pupil ID</th>
              <th>Фамилия</th>
              <th>Сумма</th>
              {canEdit ? <th /> : null}
            </tr>
          </thead>
          <tbody>
            {expenses.map((e) => (
              <tr key={e.id}>
                <td className="mono">{e.id}</td>
                <td className="mono">{formatDate(e.date)}</td>
                <td>{e.gift_for}</td>
                <td className="mono">{e.pupil_id}</td>
                <td>{e.surname}</td>
                <td className="mono">{e.summ}</td>
                {canEdit ? (
                  <td className="rowActions">
                    <button className="btn btnSmall btnSecondary" onClick={() => onEdit!(e)}>
                      Редактировать
                    </button>
                    <button className="btn btnSmall btnDanger" onClick={() => onDelete!(e.id)}>
                      Удалить
                    </button>
                  </td>
                ) : null}
              </tr>
            ))}
            {expenses.length === 0 ? (
              <tr>
                <td colSpan={canEdit ? 7 : 6} className="empty">
                  Нет данных.
                </td>
              </tr>
            ) : null}
          </tbody>
        </table>
      </div>
    </section>
  );
}
