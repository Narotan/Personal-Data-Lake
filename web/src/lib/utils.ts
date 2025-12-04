import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { subDays, format, startOfMonth, endOfMonth } from "date-fns";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function getDateRange(range: string): { start_date: string; end_date: string } {
  const today = new Date();
  let start = today;
  let end = today;

  switch (range) {
    case 'today':
      start = today;
      end = today;
      break;
    case 'yesterday':
      start = subDays(today, 1);
      end = subDays(today, 1);
      break;
    case '7d':
      start = subDays(today, 6);
      end = today;
      break;
    case '30d':
      start = subDays(today, 29);
      end = today;
      break;
    case 'month':
      start = startOfMonth(today);
      end = endOfMonth(today);
      break;
    default:
      start = today;
      end = today;
  }

  return {
    start_date: format(start, 'yyyy-MM-dd'),
    end_date: format(end, 'yyyy-MM-dd'),
  };
}

export function formatDuration(seconds: number): string {
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  if (h > 0) return `${h}h ${m}m`;
  return `${m}m`;
}
