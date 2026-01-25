export type Expense = {
  id: number;
  date: string; // RFC3339 / ISO string from Go time.Time
  gift_for: string;
  pupil_id: number;
  summ: number;
  surname: string;
};
