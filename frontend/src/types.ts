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

// Payments на бэкенде сейчас минимальные (id, pupil_id, summ).
// Поля date/gift_for/surname оставляем опциональными, чтобы UI не падал,
// если бэкенд их не отдает, и чтобы можно было расширить модель позже.
export type Payment = {
  id: number;
  pupil_id: number;
  summ: number;
  date?: string;
  gift_for?: string;
  surname?: string;
} & Record<string, unknown>;

// Ответ бэкенда GET /expenses.
// Возможные варианты ключей зависят от json-тегов в Go.
export type ExpenseWithBalanceResponse = {
  expenses?: Expense[];
  Expenses?: Expense[];
  expense?: Expense[];
  Expense?: Expense[];
  balance?: number;
  Balance?: number;
} & Record<string, unknown>;
