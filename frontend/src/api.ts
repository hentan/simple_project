import type { Expense, ExpenseWithBalanceResponse, Payment, Pupil } from "./types";

const API_BASE = (import.meta as any).env?.VITE_API_BASE ?? "/api";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {})
    },
    ...init
  });

  if (!res.ok) {
    const text = await res.text().catch(() => "");
    const err = new Error(text || `HTTP ${res.status}`);
    (err as any).status = res.status;
    throw err;
  }

  // Явно обработаем "нет контента"
  if (res.status === 204) {
    return undefined as unknown as T;
  }

  const contentType = res.headers.get("content-type") || "";
  if (!contentType.includes("application/json")) {
    return undefined as unknown as T;
  }

  // ВАЖНО: тело может быть пустым даже при application/json
  const raw = await res.text();
  if (!raw.trim()) {
    return undefined as unknown as T;
  }

  return JSON.parse(raw) as T;
}

function normalizeExpenseWithBalance(raw: any): { expenses: Expense[]; balance: number } {
  // При отсутствии json-тегов Go сериализует как {"Expense": [...], "Balance": 123}
  // С тегами может быть {"expense": [...], "balance": 123}
  const expenses = (raw?.expenses ?? raw?.Expenses ?? raw?.expense ?? raw?.Expense ?? []) as Expense[];
  const balance = Number(raw?.balance ?? raw?.Balance ?? 0);
  return { expenses: Array.isArray(expenses) ? expenses : [], balance: Number.isFinite(balance) ? balance : 0 };
}

// GET /expenses -> ExpenseWithBalanceResponse (expenses + balance)
export async function getExpensesWithBalance(): Promise<{ expenses: Expense[]; balance: number }> {
  const raw = await request<ExpenseWithBalanceResponse>("/expenses", { method: "GET" });
  return normalizeExpenseWithBalance(raw);
}

// Иногда нужен только список расходов.
export async function getExpenses(): Promise<Expense[]> {
  const { expenses } = await getExpensesWithBalance();
  return expenses;
}

export type ExpenseInput = {
  id?: number;
  date: string; // ISO string
  gift_for: string;
  pupil_id: number;
  summ: number;
};

export async function addExpense(payload: ExpenseInput): Promise<void> {
  await request<void>("/expenses", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function updateExpense(payload: ExpenseInput & { id: number }): Promise<void> {
  await request<void>("/expenses", {
    method: "PUT",
    body: JSON.stringify(payload)
  });
}

export async function deleteExpense(id: number): Promise<void> {
  // backend expects JSON body with {id: ...}
  await request<void>("/expenses", {
    method: "DELETE",
    body: JSON.stringify({ id })
  });
}

// --- Pupils ---

function normalizePupil(raw: any): Pupil {
    const id = Number(raw?.id ?? raw?.ID ?? raw?.Id ?? 0);
    const surname = String(raw?.surname ?? raw?.Surname ?? "");

    const nameRaw = raw?.name ?? raw?.Name;
    const name = typeof nameRaw === "string" ? nameRaw : undefined;

    const parentNameRaw = raw?.parent_name ?? raw?.parentName ?? raw?.ParentName;
    const parent_name = typeof parentNameRaw === "string" ? parentNameRaw : undefined;

    const parentPhoneRaw = raw?.parent_phone ?? raw?.parentPhone ?? raw?.ParentPhone;
    const parent_phone = typeof parentPhoneRaw === "string" ? parentPhoneRaw : undefined;

    return {
        id,
        surname,
        name,
        parent_name,
        parent_phone
    } as Pupil;
}

export async function getPupils(): Promise<Pupil[]> {
  const raw = await request<any>("/pupils", { method: "GET" });
  // Go может сериализовать nil-slice как null
  const arr = Array.isArray(raw) ? raw : [];
  return arr.map(normalizePupil).filter((p) => Number.isFinite(p.id) && p.id > 0);
}

export type PupilInput = {
    id?: number;
    surname: string;
    name?: string;
    parent_name?: string;
    parent_phone?: string;
};

export async function addPupil(payload: PupilInput): Promise<void> {
  await request<void>("/pupils", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function updatePupil(payload: PupilInput & { id: number }): Promise<void> {
  await request<void>("/pupils", {
    method: "PUT",
    body: JSON.stringify(payload)
  });
}

export async function deletePupil(id: number): Promise<void> {
  await request<void>("/pupils", {
    method: "DELETE",
    body: JSON.stringify({ id })
  });
}

// --- Payments ---

export async function getPayments(): Promise<Payment[]> {
  const raw = await request<any>("/payments", { method: "GET" });
  // Go может сериализовать nil-slice как null
  return Array.isArray(raw) ? (raw as Payment[]) : [];
}

export type PaymentInput = {
  id?: number;
  date: string;
  gift_for?: string;
  pupil_id: number;
  summ: number;
};

export async function addPayment(payload: PaymentInput): Promise<void> {
  await request<void>("/payments", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function updatePayment(payload: PaymentInput & { id: number }): Promise<void> {
  await request<void>("/payments", {
    method: "PUT",
    body: JSON.stringify(payload)
  });
}

export async function deletePayment(id: number): Promise<void> {
  await request<void>("/payments", {
    method: "DELETE",
    body: JSON.stringify({ id })
  });
}
