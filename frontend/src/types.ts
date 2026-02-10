export type Expense = {
  id: number;
  date: string; // RFC3339 / ISO string from Go time.Time
  gift_for: string;
  pupil_id: number;
  summ: number;
  surname: string;
};

export type Pupil = {
  id: number;
  surname: string;
  name?: string;
} & Record<string, unknown>;

export type Payment = {
  id: number;
  date: string;
  gift_for?: string;
  pupil_id: number;
  summ: number;
  surname?: string;
} & Record<string, unknown>;

// Ответ бэкенда GET /expenses.
// В Go это struct models.ExpenseWithBalance { Expense []Expense; Balance int }
// (возможны json-теги, поэтому на клиенте нормализуем ключи).
export type ExpenseWithBalanceResponse = {
  expense: Expense[];
  balance: number;
} & Record<string, unknown>;
