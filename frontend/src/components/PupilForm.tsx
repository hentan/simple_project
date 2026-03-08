import { useEffect, useMemo, useState } from "react";
import type { Pupil } from "../types";
import type { PupilInput } from "../api";

type Props = {
  editing?: Pupil | null;
  onCreate: (payload: PupilInput) => Promise<void>;
  onUpdate: (payload: PupilInput & { id: number }) => Promise<void>;
  onCancelEdit: () => void;
};

export default function PupilForm({ editing, onCreate, onUpdate, onCancelEdit }: Props) {
  const isEdit = !!editing;

  const initialSurname = useMemo(() => editing?.surname ?? "", [editing]);
  const initialName = useMemo(() => {
    const raw = (editing as any)?.name;
    return typeof raw === "string" ? raw : "";
  }, [editing]);

    const initialParentName = useMemo(() => {
        const raw = (editing as any)?.parent_name;
        return typeof raw === "string" ? raw : "";
    }, [editing]);

    const initialParentPhone = useMemo(() => {
        const raw = (editing as any)?.parent_phone;
        return typeof raw === "string" ? raw : "";
    }, [editing]);

  const [surname, setSurname] = useState<string>("");
  const [name, setName] = useState<string>("");
  const [parentName, setParentName] = useState<string>("");
  const [parentPhone, setParentPhone] = useState<string>("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (editing) {
            setSurname(initialSurname);
            setName(initialName);
            setParentName(initialParentName);
            setParentPhone(initialParentPhone);
            setError(null);
        } else {
            setSurname("");
            setName("");
            setParentName("");
            setParentPhone("");
            setError(null);
        }
    }, [editing, initialSurname, initialName, initialParentName, initialParentPhone]);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);

    if (!surname.trim()) return setError("Укажите фамилию.");

      const payload: PupilInput = {
          surname: surname.trim(),
          name: name.trim() || undefined,
          parent_name: parentName.trim() || undefined,
          parent_phone: parentPhone.trim() || undefined
      };

    try {
      setSubmitting(true);
      if (isEdit && editing) {
        await onUpdate({ ...payload, id: editing.id });
      } else {
        await onCreate(payload);
        setSurname("");
        setName("");
        setParentName("");
        setParentPhone("");
      }
    } catch (err: any) {
      setError(err?.message ?? "Ошибка запроса");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section className="card">
      <div className="cardHeader">
        <h2 className="h2">{isEdit ? `Редактирование ученика #${editing?.id}` : "Добавить ученика"}</h2>
        {isEdit ? (
          <button type="button" className="btn btnSecondary" onClick={onCancelEdit} disabled={submitting}>
            Отменить
          </button>
        ) : null}
      </div>

      <form onSubmit={handleSubmit} className="formGrid">
        <label className="field">
          <span className="label">Фамилия</span>
          <input value={surname} onChange={(e) => setSurname(e.target.value)} placeholder="Например: Иванов" />
        </label>

        <label className="field">
          <span className="label">Имя (необязательно)</span>
          <input value={name} onChange={(e) => setName(e.target.value)} placeholder="Например: Пётр" />
        </label>

        <label className="field">
          <span className="label">Имя родителя (необязательно)</span>
          <input
              value={parentName}
              onChange={(e) => setParentName(e.target.value)}
              placeholder="Например: Мария Ивановна"
          />
        </label>

        <label className="field">
          <span className="label">Телефон родителя (необязательно)</span>
          <input
              value={parentPhone}
              onChange={(e) => setParentPhone(e.target.value)}
              placeholder="Например: +7 999 123-45-67"
          />
        </label>

        <div className="actions">
          <button className="btn" type="submit" disabled={submitting}>
            {submitting ? "Сохранение..." : isEdit ? "Сохранить" : "Добавить"}
          </button>
          {error ? <div className="error">{error}</div> : null}
          <div className="hint">
            API: <code>/pupils</code> (GET/POST/PUT/DELETE).
          </div>
        </div>
      </form>
    </section>
  );
}
