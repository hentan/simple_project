import type { Pupil } from "../types";

type Props = {
  pupils: Pupil[];
  onEdit: (p: Pupil) => void;
  onDelete: (id: number) => void;
};

export default function PupilTable({ pupils, onEdit, onDelete }: Props) {
  return (
    <section className="card">
      <div className="cardHeader">
        <h2 className="h2">Ученики</h2>
        <div className="muted">Всего: {pupils.length}</div>
      </div>

      <div className="tableWrap">
        <table className="table">
            <thead>
            <tr>
                <th>ID</th>
                <th>Фамилия</th>
                <th>Имя</th>
                <th>Родитель</th>
                <th>Телефон родителя</th>
                <th />
            </tr>
            </thead>
          <tbody>
            {pupils.map((p) => (
                <tr key={p.id}>
                    <td className="mono">{p.id}</td>
                    <td>{p.surname}</td>
                    <td>{typeof (p as any).name === "string" ? String((p as any).name) : ""}</td>
                    <td>{typeof (p as any).parent_name === "string" ? String((p as any).parent_name) : ""}</td>
                    <td>{typeof (p as any).parent_phone === "string" ? String((p as any).parent_phone) : ""}</td>
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
            {pupils.length === 0 ? (
              <tr>
                  <td colSpan={6} className="empty">
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
