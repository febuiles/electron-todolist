interface Todo {
  id: number;
  title: string;
  creator: string;
  lastUpdated: string;
  column: ColumnType;
}

type ColumnType = 'todo' | 'ongoing' | 'done';

interface Column {
  id: ColumnType;
  title: string;
}
