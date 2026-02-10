import { useCallback, useEffect, useState } from "react";
import type { Pupil } from "../types";
import { addPupil, deletePupil, getPupils, updatePupil, type PupilInput } from "../api";
import PupilForm from "../components/PupilForm";
import PupilTable from "../components/PupilTable";

export default function Pupils() {
  const [pupils, setPupils] = useState<Pupil[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [editing, setEditing] = useState<Pupil | null>(null);

  const refresh = useCallback(async () => {
    setError(null);
    setLoading(true);
    try {
      const data = await getPupils();
      setPupils(data);
    } catch (err: any) {
      setError(err?.message ?? "Ошибка загрузки");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const handleCreate = useCallback(async (payload: PupilInput) => {
    await addPupil(payload);
    await refresh();
  }, [refresh]);

  const handleUpdate = useCallback(async (payload: PupilInput & { id: number }) => {
    await updatePupil(payload);
    setEditing(null);
    await refresh();
  }, [refresh]);

  const handleDelete = useCallback(async (id: number) => {
    const ok = window.confirm(`Удалить ученика #${id}?`);
    if (!ok) return;
    try {
      await deletePupil(id);
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
          <h1 className="h1">Ученики</h1>
          <div className="muted">
            Справочник учеников. API: <code>/pupils</code> (GET/POST/PUT/DELETE).
          </div>
        </div>
      </header>

      <main className="grid">
        <PupilForm
          editing={editing}
          onCreate={handleCreate}
          onUpdate={handleUpdate}
          onCancelEdit={() => setEditing(null)}
        />

        <section>
          {error ? <div className="error">Ошибка: {error}</div> : null}
          {loading ? <div className="muted">Загрузка...</div> : null}
          <PupilTable pupils={pupils} onEdit={setEditing} onDelete={handleDelete} />
        </section>
      </main>
    </>
  );
}
