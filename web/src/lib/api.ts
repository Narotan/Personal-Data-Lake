import axios from 'axios';

const api = axios.create({
  baseURL: '/api/v1',
  headers: {
    'X-API-Key': import.meta.env.VITE_API_KEY,
  },
});

export interface DateRange {
  start_date: string; // YYYY-MM-DD
  end_date: string;   // YYYY-MM-DD
}

export interface ProjectStat {
  name: string;
  total_seconds: number;
}

export interface LanguageStat {
  name: string;
  total_seconds: number;
  percent: number;
}

export interface DailyStat {
  date: string;
  total_seconds: number;
  text: string;
  projects: ProjectStat[];
  languages: LanguageStat[];
}

export interface DailyFitStat {
  date: string;
  steps: number;
  distance: number;
}

export interface CalendarEvent {
  id: string;
  summary: string;
  description: string;
  start_time: string;
  end_time: string;
}

export const fetchWakaTimeStats = async (range: DateRange): Promise<DailyStat[]> => {
  const { data } = await api.get('/wakatime/stats', { params: range });
  return data;
};

export const fetchGoogleFitStats = async (range: DateRange): Promise<DailyFitStat[]> => {
  const { data } = await api.get('/googlefit/stats', { params: range });
  return data;
};

export const fetchGoogleCalendarEvents = async (range: DateRange): Promise<CalendarEvent[]> => {
  const { data } = await api.get('/googlecalendar/events', { params: range });
  return data;
};

export interface AuthStatus {
  wakatime: boolean;
  googlefit: boolean;
  googlecalendar: boolean;
}

export const fetchAuthStatus = async (): Promise<AuthStatus> => {
  const { data } = await api.get('/auth/status');
  return data;
};

export interface AppStat {
  App: string;
  TotalDuration: number;
  EventCount: number;
}

export const fetchActivityWatchStats = async (range: DateRange): Promise<AppStat[]> => {
  const { data } = await api.get('/activitywatch/stats', { 
    params: {
      start: range.start_date + 'T00:00:00Z',
      end: range.end_date + 'T23:59:59Z'
    } 
  });
  return data;
};

export interface AggregatedLanguageStat {
  name: string;
  total_seconds: number;
  percent: number;
}

export interface AggregatedProjectStat {
  name: string;
  total_seconds: number;
}

export const fetchTopLanguages = async (range: DateRange, limit: number = 5): Promise<AggregatedLanguageStat[]> => {
  const { data } = await api.get('/wakatime/top-languages', { 
    params: { ...range, limit } 
  });
  return data;
};

export const fetchTopProjects = async (range: DateRange, limit: number = 5): Promise<AggregatedProjectStat[]> => {
  const { data } = await api.get('/wakatime/top-projects', { 
    params: { ...range, limit } 
  });
  return data;
};
