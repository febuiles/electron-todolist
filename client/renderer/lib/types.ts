interface Todo {
  id: number;
  title: string;
  user_id: number;
  lastUpdated: string;
  column: ColumnType;
  creator: string;
}

type ColumnType = 'todo' | 'ongoing' | 'done';

interface Column {
  id: ColumnType;
  title: string;
}
