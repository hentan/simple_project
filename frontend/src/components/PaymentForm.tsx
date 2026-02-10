import { useEffect, useMemo, useState } from "react";
import type { Payment } from "../types";
import type { PaymentInput } from "../api";

type Props = {
  editing?: Payment | null;
  onCreate: (payload: PaymentInput) => Promise<void>;
  onUpdate: (payload: PaymentInput & { id: number }) => Promise<void>;
  onCancelEdit: () => void;
};

function isoFromDateInput(value: string): string {
  return new Date(`${value}T00:00:00Z`).toISOString();
}

function dateInputFromIso(iso: string): string {
  const d = new Date(iso);
  const y = d.getUTCFullYear();
  const m = String(d.getUTCMonth() + 1).padStart(2, "0");
  const day = String(d.getUTCDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}

export default function PaymentForm({ editing, onCreate, onUpdate, onCancelEdit }: Props) {
  const isEdit = !!editing;

  const initialDate = useMemo(() => {
    if (!editing?.date) return "";
    return dateInputFromIso(editing.date);
  }, [editing]);

  const [date, setDate] = useState<string>("");
  const [giftFor, setGiftFor] = useState<string>("");
  const [pupilId, setPupilId] = useState<string>("");
  const [summ, setSumm] = useState<string>("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (editing) {
      setDate(initialDate);
      setGiftFor((editing as any)?.gift_for ? String((editing as any).gift_for) : "");
      setPupilId(String(editing.pupil_id ?? ""));
      setSumm(String(editing.summ ?? ""));
      setError(null);
    } else {
      setDate("");
      setGiftFor("");
      setPupilId("");
      setSumm("");
      setError(null);
    }
  }, [editing, initialDate]);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);

    if (!date) return setError("Укажите дату.");
    const pid = Number(pupilId);
    if (!Number.isFinite(pid) || pid <= 0) return setError("pupil_id должен быть положительным числом.");
    const s = Number(summ);
    if (!Number.isFinite(s) || s < 0) return setError("summ должен быть числом (>= 0).");

    const payload: PaymentInput = {
      date: isoFromDateInput(date),
      pupil_id: pid,
      summ: s,
      gift_for: giftFor.trim() || undefined
    };

    try {
      setSubmitting(true);
      if (isEdit && editing) {
        await onUpdate({ ...payload, id: editing.id });
      } else {
        await onCreate(payload);
        setDate("");
        setGiftFor("");
        setPupilId("");
        setSumm("");
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
        <h2 className="h2">{isEdit ? `Редактирование оплаты #${editing?.id}` : "Добавить оплату"}</h2>
        {isEdit ? (
          <button type="button" className="btn btnSecondary" onClick={onCancelEdit} disabled={submitting}>
            Отменить
          </button>
        ) : null}
      </div>

      <form onSubmit={handleSubmit} className="formGrid">
        <label className="field">
          <span className="label">Дата</span>
          <input type="date" value={date} onChange={(e) => setDate(e.target.value)} required />
        </label>

        <label className="field">
          <span className="label">Назначение (опционально)</span>
          <input value={giftFor} onChange={(e) => setGiftFor(e.target.value)} placeholder="Например: занятие / подарок" />
        </label>

        <label className="field">
          <span className="label">Pupil ID</span>
          <input inputMode="numeric" value={pupilId} onChange={(e) => setPupilId(e.target.value)} placeholder="Например: 12" />
        </label>

        <label className="field">
          <span className="label">Summ</span>
          <input inputMode="numeric" value={summ} onChange={(e) => setSumm(e.target.value)} placeholder="Например: 1500" />
        </label>

        <div className="actions">
          <button className="btn" type="submit" disabled={submitting}>
            {submitting ? "Сохранение..." : isEdit ? "Сохранить" : "Добавить"}
          </button>
          {error ? <div className="error">{error}</div> : null}
          <div className="hint">
            API: <code>/payments</code> (GET/POST/PUT/DELETE).
          </div>
        </div>
      </form>
    </section>
  );
}
