import { DailyStat, DailyFitStat, CalendarEvent, AppStat, AggregatedLanguageStat, AggregatedProjectStat } from './api';

export const mockWakaTimeData: DailyStat[] = Array.from({ length: 7 }, (_, i) => {
  const date = new Date();
  date.setDate(date.getDate() - i);
  const dateStr = date.toISOString().split('T')[0];
  return {
    date: dateStr,
    total_seconds: Math.random() * 28800 + 3600, // 1-9 hours
    text: "Coding",
    projects: [],
    languages: []
  };
}).reverse();

export const mockGoogleFitData: DailyFitStat[] = Array.from({ length: 7 }, (_, i) => {
  const date = new Date();
  date.setDate(date.getDate() - i);
  const dateStr = date.toISOString().split('T')[0];
  return {
    date: dateStr,
    steps: Math.floor(Math.random() * 10000) + 2000,
    distance: Math.random() * 5000 + 1000
  };
}).reverse();

export const mockTopLanguages: AggregatedLanguageStat[] = [
  { name: 'TypeScript', total_seconds: 150000, percent: 45 },
  { name: 'Go', total_seconds: 100000, percent: 30 },
  { name: 'Python', total_seconds: 50000, percent: 15 },
  { name: 'SQL', total_seconds: 20000, percent: 6 },
  { name: 'HTML/CSS', total_seconds: 13000, percent: 4 },
];

export const mockTopProjects: AggregatedProjectStat[] = [
  { name: 'personal-data-lake', total_seconds: 120000 },
  { name: 'website-v2', total_seconds: 80000 },
  { name: 'learning-go', total_seconds: 60000 },
  { name: 'advent-of-code', total_seconds: 40000 },
];

export const mockActivityWatchData: AppStat[] = [
    { App: "Code", TotalDuration: 18000, EventCount: 100 },
    { App: "Chrome", TotalDuration: 12000, EventCount: 500 },
    { App: "Slack", TotalDuration: 3600, EventCount: 50 },
    { App: "Terminal", TotalDuration: 7200, EventCount: 200 },
];
