import React, { useEffect, useMemo, useState } from "react";
import type { Expense, Pupil } from "../types";
import { getPupils, type ExpenseInput } from "../api";

type Props = {
  editing?: Expense | null;
  onCreate: (payload: ExpenseInput) => Promise<void>;
  onUpdate: (payload: ExpenseInput & { id: number }) => Promise<void>;
  onCancelEdit: () => void;
};

function isoFromDateInput(value: string): string {
  // Keep a stable "date-only" representation without timezone surprises:
  // send midnight UTC for the chosen day.
  return new Date(`${value}T00:00:00Z`).toISOString();
}

function dateInputFromIso(iso: string): string {
  const d = new Date(iso);
  const y = d.getUTCFullYear();
  const m = String(d.getUTCMonth() + 1).padStart(2, "0");
  const day = String(d.getUTCDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}

function normalizePupil(raw: any): Pupil {
  // Поддержка старых ответов Go без json-тегов: ID/Name/Surname
  const id = Number(raw?.id ?? raw?.ID ?? raw?.Id ?? 0);
  const surname = String(raw?.surname ?? raw?.Surname ?? "");
  const nameRaw = raw?.name ?? raw?.Name;
  const name = typeof nameRaw === "string" ? nameRaw : undefined;
  return { id, surname, name } as Pupil;
}

function pupilLabel(p: Pupil): string {
  const surname = (p.surname ?? "").trim();
  const name = (typeof (p as any).name === "string" ? String((p as any).name) : "").trim();
  return `${surname}${name ? " " + name : ""}`.trim();
}

export default function ExpenseForm({ editing, onCreate, onUpdate, onCancelEdit }: Props) {
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

  // Справочник учеников для выпадающего списка
  const [pupils, setPupils] = useState<Pupil[]>([]);
  const [pupilsLoading, setPupilsLoading] = useState(false);
  const [pupilsError, setPupilsError] = useState<string | null>(null);

  useEffect(() => {
    let alive = true;
    (async () => {
      setPupilsError(null);
      setPupilsLoading(true);
      try {
        const list = await getPupils();
        const normalized = list
          .map(normalizePupil)
          .filter((p) => Number.isFinite(p.id) && p.id > 0 && String(p.surname ?? "").trim().length > 0)
          .sort((a, b) =>
            String(a.surname ?? "").localeCompare(String(b.surname ?? ""), "ru") ||
            String((a as any).name ?? "").localeCompare(String((b as any).name ?? ""), "ru")
          );
        if (alive) setPupils(normalized);
      } catch (err: any) {
        if (alive) {
          setPupils([]);
          setPupilsError(err?.message ?? "Не удалось загрузить список учеников");
        }
      } finally {
        if (alive) setPupilsLoading(false);
      }
    })();
    return () => {
      alive = false;
    };
  }, []);

  useEffect(() => {
    if (editing) {
      setDate(initialDate);
      setGiftFor(editing.gift_for ?? "");
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
    if (!giftFor.trim()) return setError("Укажите назначение (gift_for).");
    const pid = Number(pupilId);
    if (!Number.isFinite(pid) || pid <= 0) return setError("Выберите ученика (pupil_id)." );
    const s = Number(summ);
    if (!Number.isFinite(s) || s < 0) return setError("summ должен быть числом (>= 0)." );

    const payload: ExpenseInput = {
      date: isoFromDateInput(date),
      gift_for: giftFor.trim(),
      pupil_id: pid,
      summ: s
    };

    try {
      setSubmitting(true);
      if (isEdit && editing) {
        await onUpdate({ ...payload, id: editing.id });
      } else {
        await onCreate(payload);
      }
      if (!isEdit) {
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

  const selectedMissing = pupilId && !pupils.some((p) => String(p.id) === String(pupilId));

  return (
    <section className="card">
      <div className="cardHeader">
        <h2 className="h2">{isEdit ? `Редактирование записи #${editing?.id}` : "Добавить сдачу"}</h2>
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
          <span className="label">Назначение</span>
          <input value={giftFor} onChange={(e) => setGiftFor(e.target.value)} placeholder="Например: Новый год" />
        </label>

        <label className="field">
          <span className="label">Ученик</span>
          <select
            value={pupilId}
            onChange={(e) => setPupilId(e.target.value)}
            disabled={pupilsLoading || submitting}
          >
            <option value="">— выберите ученика —</option>
            {selectedMissing ? <option value={pupilId}>ID {pupilId} (нет в списке)</option> : null}
            {pupils.map((p) => (
              <option key={p.id} value={p.id}>
                {pupilLabel(p)}
              </option>
            ))}
          </select>
          {pupilsError ? <div className="hint error">{pupilsError}</div> : null}
          {!pupilsError && pupils.length === 0 && !pupilsLoading ? (
            <div className="hint">Сначала добавьте учеников в справочник (раздел «Ученики»).</div>
          ) : null}
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
            API: <code>/expenses</code> (GET/POST/PUT/DELETE), проксируется через <code>/api</code>.
          </div>
        </div>
      </form>
    </section>
  );
}
