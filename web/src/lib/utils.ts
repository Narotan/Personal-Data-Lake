import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { subDays, format, startOfMonth, endOfMonth, addDays, parseISO } from "date-fns";

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
    case 'year':
      start = subDays(today, 364); // 365 days including today
      end = today;
      break;
    case 'all_time':
      start = new Date(2020, 0, 1); // Assuming project start around 2020
      end = today;
      break;
    default:
      if (range.startsWith('custom:')) {
        const parts = range.split(':');
        if (parts.length === 3) {
            return {
                start_date: parts[1],
                end_date: parts[2]
            };
        }
      }
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

export function fillMissingDates(data: any[], startDate: string, endDate: string): any[] {
  const start = parseISO(startDate);
  const end = parseISO(endDate);
  const filled = [];
  
  let current = start;
  while (current <= end) {
    const dateStr = format(current, 'yyyy-MM-dd');
    // Check if data exists for this date (ignoring time part if present)
    const existing = data.find(item => item.date.startsWith(dateStr));
    
    if (existing) {
      filled.push(existing);
    } else {
      filled.push({
        date: dateStr,
        steps: 0,
        distance: 0,
        total_seconds: 0,
        grand_total: { total_seconds: 0 }, // For wakatime structure if needed
      });
    }
    current = addDays(current, 1);
  }
  return filled;
}

export function aggregateByMonth(data: any[]): any[] {
  const monthlyData: { [key: string]: any } = {};
  
  data.forEach(item => {
    const date = parseISO(item.date);
    const monthKey = format(date, 'yyyy-MM');
    
    if (!monthlyData[monthKey]) {
      monthlyData[monthKey] = {
        date: format(date, 'yyyy-MM-01'),
        total_seconds: 0,
        steps: 0,
        distance: 0,
        count: 0,
      };
    }
    
    monthlyData[monthKey].total_seconds += item.total_seconds || 0;
    monthlyData[monthKey].steps += item.steps || 0;
    monthlyData[monthKey].distance += item.distance || 0;
    monthlyData[monthKey].count += 1;
  });
  
  return Object.values(monthlyData).sort((a, b) => a.date.localeCompare(b.date));
}
