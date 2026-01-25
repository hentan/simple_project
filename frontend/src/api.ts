import type { Expense } from "./types";

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
    throw new Error(text || `HTTP ${res.status}`);
  }

  // 204 / empty
  const contentType = res.headers.get("content-type") || "";
  if (!contentType.includes("application/json")) {
    return undefined as unknown as T;
  }
  return (await res.json()) as T;
}

export async function getExpenses(): Promise<Expense[]> {
  return request<Expense[]>("/expenses", { method: "GET" });
}

export type ExpenseInput = {
  id?: number;
  date: string; // ISO string
  gift_for: string;
  pupil_id: number;
  summ: number;
};

export async function addExpense(payload: ExpenseInput): Promise<void> {
  await request<void>("/expense", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function updateExpense(payload: ExpenseInput & { id: number }): Promise<void> {
  await request<void>("/expense", {
    method: "PUT",
    body: JSON.stringify(payload)
  });
}

export async function deleteExpense(id: number): Promise<void> {
  // backend expects JSON body with {id: ...}
  await request<void>("/expense", {
    method: "DELETE",
    body: JSON.stringify({ id })
  });
}
