import { useState, useEffect } from 'react';
import { Activity, Calendar as CalendarIcon, Code2, Monitor } from 'lucide-react';
import { motion } from 'framer-motion';
import { Header } from './components/layout/Header';
import { Sidebar } from './components/layout/Sidebar';
import { KPICard } from './components/dashboard/KPICard';
import { ProductivityChart } from './components/dashboard/ProductivityChart';
import { HealthChart } from './components/dashboard/HealthChart';
import { ScheduleTimeline } from './components/dashboard/ScheduleTimeline';
import { LanguageDistribution } from './components/dashboard/LanguageDistribution';
import { TopProjects } from './components/dashboard/TopProjects';
import { TopApplications } from './components/dashboard/TopApplications';
import { Setup } from './components/Setup';
import { AuthSuccess } from './components/AuthSuccess';
import { fetchWakaTimeStats, fetchGoogleFitStats, fetchGoogleCalendarEvents, fetchActivityWatchStats, fetchTopLanguages, fetchTopProjects, DailyStat, DailyFitStat, CalendarEvent, AppStat, AggregatedLanguageStat, AggregatedProjectStat } from './lib/api';
import { getDateRange, formatDuration, fillMissingDates } from './lib/utils';
import { ThemeProvider } from './lib/theme';


function AppContent() {
  const [dateRange, setDateRange] = useState('today');
  const [loading, setLoading] = useState(true);
  const [view, setView] = useState<'dashboard' | 'setup' | 'auth-success'>('dashboard');
  
  const [wakaTimeData, setWakaTimeData] = useState<DailyStat[]>([]);
  const [productivityTrendData, setProductivityTrendData] = useState<DailyStat[]>([]);
  const [googleFitData, setGoogleFitData] = useState<DailyFitStat[]>([]);
  const [healthTrendData, setHealthTrendData] = useState<DailyFitStat[]>([]);
  const [calendarData, setCalendarData] = useState<CalendarEvent[]>([]);
  const [activityWatchData, setActivityWatchData] = useState<AppStat[]>([]);
  const [topLanguages, setTopLanguages] = useState<AggregatedLanguageStat[]>([]);
  const [topProjects, setTopProjects] = useState<AggregatedProjectStat[]>([]);

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    if (params.get('auth_success') === 'true') {
      setView('auth-success');
      window.history.replaceState({}, '', '/');
    }
  }, []);

  useEffect(() => {
    if (view === 'setup' || view === 'auth-success') return;

    const loadData = async () => {
      setLoading(true);
      try {
        const range = getDateRange(dateRange);
        
        // Determine range for trend charts (at least 7 days)
        let trendRange = range;
        if (dateRange === 'today' || dateRange === 'yesterday') {
             trendRange = getDateRange('7d');
        }

        const [wakaData, fitData, calData, awData, langData, projData, wakaTrend, fitTrend] = await Promise.all([
          fetchWakaTimeStats(range).catch(err => { console.error('WakaTime error:', err); return []; }),
          fetchGoogleFitStats(range).catch(err => { console.error('GoogleFit error:', err); return []; }),
          fetchGoogleCalendarEvents(range).catch(err => { console.error('Calendar error:', err); return []; }),
          fetchActivityWatchStats(range).catch(err => { console.error('ActivityWatch error:', err); return []; }),
          fetchTopLanguages(range, 5).catch(err => { console.error('Top Languages error:', err); return []; }),
          fetchTopProjects(range, 5).catch(err => { console.error('Top Projects error:', err); return []; }),
          (dateRange === 'today' || dateRange === 'yesterday') ? fetchWakaTimeStats(trendRange).catch(err => { console.error('WakaTrend error:', err); return []; }) : Promise.resolve(null),
          (dateRange === 'today' || dateRange === 'yesterday') ? fetchGoogleFitStats(trendRange).catch(err => { console.error('FitTrend error:', err); return []; }) : Promise.resolve(null),
        ]);

        // Helper to check if data is effectively empty (all zeros)
        const isEmpty = (data: any[]) => {
          if (!data || data.length === 0) return true;
          // Check for specific types
          if ('total_seconds' in data[0]) return data.every((d: any) => d.total_seconds === 0);
          if ('steps' in data[0]) return data.every((d: any) => d.steps === 0);
          if ('TotalDuration' in data[0]) return data.every((d: any) => d.TotalDuration === 0);
          return false;
        };

        setWakaTimeData(wakaData || []);
        setGoogleFitData(fitData || []);
        
        // Set trend data
        if (wakaTrend) {
            setProductivityTrendData(fillMissingDates(wakaTrend || [], trendRange.start_date, trendRange.end_date));
        } else {
            setProductivityTrendData(fillMissingDates(wakaData || [], trendRange.start_date, trendRange.end_date));
        }

        if (fitTrend) {
            setHealthTrendData(fillMissingDates(fitTrend || [], trendRange.start_date, trendRange.end_date));
        } else {
            setHealthTrendData(fillMissingDates(fitData || [], trendRange.start_date, trendRange.end_date));
        }

        setCalendarData(calData || []);
        setActivityWatchData(awData || []);
        setTopLanguages(langData || []);
        setTopProjects(projData || []);
      } catch (error) {
        console.error("Failed to fetch data:", error);
        setWakaTimeData([]);
        setProductivityTrendData([]);
        setGoogleFitData([]);
        setHealthTrendData([]);
        setTopLanguages([]);
        setTopProjects([]);
        setActivityWatchData([]);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [dateRange, view]);

  const totalCodingSeconds = wakaTimeData.reduce((acc, day) => acc + day.total_seconds, 0);
  const totalSteps = googleFitData.reduce((acc, day) => acc + day.steps, 0);
  const totalMeetingsDuration = calendarData.reduce((acc, event) => {
    const start = new Date(event.start_time).getTime();
    const end = new Date(event.end_time).getTime();
    return acc + (end - start) / 1000;
  }, 0);
  const totalPCActiveSeconds = activityWatchData.reduce((acc, app) => acc + app.TotalDuration, 0);

  return (
    <div className="flex min-h-screen bg-gradient-to-br from-[#F2F1EA] via-[#F8F7F4] to-[#FFFFFF] dark:from-background dark:via-background dark:to-background text-slate-900 dark:text-white font-sans selection:bg-primary/30 transition-colors duration-300 relative">
      <div className="absolute inset-0 bg-grid-pattern pointer-events-none fixed" />
      {view !== 'auth-success' && (
        <Sidebar 
            currentView={view === 'auth-success' ? 'dashboard' : view} 
            onViewChange={(v) => setView(v)} 
        />
      )}
      
      <main className="flex-1 p-4 md:p-8 overflow-y-auto h-screen relative z-10">
        <div className="max-w-7xl mx-auto space-y-8">
          {view !== 'auth-success' && (
            <Header 
              dateRangeLabel="Today" 
              currentRange={dateRange}
              onRangeChange={setDateRange}
              onSetupClick={() => setView(view === 'dashboard' ? 'setup' : 'dashboard')}
              isSetupMode={view === 'setup'}
            />
          )}
          
          {view === 'auth-success' ? (
            <AuthSuccess onContinue={() => setView('setup')} />
          ) : view === 'setup' ? (
            <Setup />
          ) : (
            <>
              <motion.div 
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5 }}
                className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-6"
              >
                <KPICard
                  title="Coding Time"
                  value={loading ? "..." : formatDuration(totalCodingSeconds)}
                  subtitle={dateRange === 'today' ? "Today" : "Selected Period"}
                  icon={Code2}
                  color="wakatime"
                />
                <KPICard
                  title="Steps"
                  value={loading ? "..." : totalSteps.toLocaleString()}
                  subtitle={dateRange === 'today' ? "Today" : "Selected Period"}
                  icon={Activity}
                  color="googlefit"
                />
                <KPICard
                  title="Meetings"
                  value={loading ? "..." : formatDuration(totalMeetingsDuration)}
                  subtitle={`${calendarData.length} events`}
                  icon={CalendarIcon}
                  color="calendar"
                />
                <KPICard
                  title="PC Active Time"
                  value={loading ? "..." : formatDuration(totalPCActiveSeconds)}
                  subtitle="ActivityWatch"
                  icon={Monitor}
                  color="activity"
                />
              </motion.div>

              <motion.div 
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: 0.2 }}
                className="grid grid-cols-1 lg:grid-cols-3 gap-6"
              >
                {/* Productivity Section */}
                <div className="lg:col-span-2 space-y-6">
                  <ProductivityChart data={productivityTrendData} loading={loading} />
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <LanguageDistribution data={topLanguages} loading={loading} />
                      <TopProjects data={topProjects} loading={loading} />
                  </div>
                  <TopApplications data={activityWatchData} loading={loading} />
                </div>

                {/* Health & Life Section */}
                <div className="space-y-6">
                  <HealthChart data={healthTrendData} loading={loading} />
                  <ScheduleTimeline data={calendarData} loading={loading} />
                </div>
              </motion.div>
            </>
          )}
        </div>
      </main>
    </div>
  );
}

function App() {
  return (
    <ThemeProvider defaultTheme="light" storageKey="vite-ui-theme">
      <AppContent />
    </ThemeProvider>
  )
}

export default App
