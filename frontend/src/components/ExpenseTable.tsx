import React from "react";
import type { Expense } from "../types";

type Props = {
  expenses: Expense[];
  onEdit: (e: Expense) => void;
  onDelete: (id: number) => void;
};

function formatDate(iso: string): string {
  const d = new Date(iso);
  // show as YYYY-MM-DD (UTC) to match input expectations
  const y = d.getUTCFullYear();
  const m = String(d.getUTCMonth() + 1).padStart(2, "0");
  const day = String(d.getUTCDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}

export default function ExpenseTable({ expenses, onEdit, onDelete }: Props) {
  return (
    <section className="card">
      <div className="cardHeader">
        <h2 className="h2">Расходы</h2>
        <div className="muted">Всего: {expenses.length}</div>
      </div>

      <div className="tableWrap">
        <table className="table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Дата</th>
              <th>Подарок для</th>
              <th>Pupil ID</th>
              <th>Фамилия</th>
              <th>Сумма</th>
              <th />
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
                <td className="rowActions">
                  <button className="btn btnSmall btnSecondary" onClick={() => onEdit(e)}>
                    Редактировать
                  </button>
                  <button className="btn btnSmall btnDanger" onClick={() => onDelete(e.id)}>
                    Удалить
                  </button>
                </td>
              </tr>
            ))}
            {expenses.length === 0 ? (
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
