import type { Payment } from "../types";

type Props = {
  payments: Payment[];
  onEdit: (p: Payment) => void;
  onDelete: (id: number) => void;
};

function formatDate(iso?: string): string {
  if (!iso) return "";
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return "";
  const y = d.getUTCFullYear();
  const m = String(d.getUTCMonth() + 1).padStart(2, "0");
  const day = String(d.getUTCDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}

export default function PaymentTable({ payments, onEdit, onDelete }: Props) {
  return (
    <section className="card">
      <div className="cardHeader">
        <h2 className="h2">Оплаты</h2>
        <div className="muted">Всего: {payments.length}</div>
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
              <th />
            </tr>
          </thead>
          <tbody>
            {payments.map((p) => (
              <tr key={p.id}>
                <td className="mono">{p.id}</td>
                <td className="mono">{formatDate(p.date)}</td>
                <td>{typeof (p as any).purpose === "string" ? String((p as any).purpose) : ""}</td>
                <td className="mono">{p.pupil_id}</td>
                <td>{typeof (p as any).surname === "string" ? String((p as any).surname) : ""}</td>
                <td className="mono">{p.summ}</td>
                <td className="rowActions">
                  <button className="btn btnSmall btnSecondary" onClick={() => onEdit(p)}>
                    Редактировать
                  </button>
                  <button className="btn btnSmall btnDanger" onClick={() => onDelete(p.id)}>
                    Удалить
                  </button>
                </td>
              </tr>
            ))}
            {payments.length === 0 ? (
              <tr>
                <td colSpan={7} className="empty">
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
