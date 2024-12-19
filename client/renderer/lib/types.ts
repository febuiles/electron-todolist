export interface Todo {
  id: number;
  title: string;
  user_id: number;
  lastUpdated: string;
  column: ColumnType;
  creator: string;
}

export type ColumnType = 'todo' | 'ongoing' | 'done';

export interface Columnable {
  id: ColumnType;
  title: string;
}

export interface User {
  id: number;
  username: string;
  lastUsedTodolistId: number;
}
