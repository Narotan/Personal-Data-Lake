import { Check, X, ExternalLink, Activity, Calendar, Code2 } from 'lucide-react';
import { cn } from '../lib/utils';
import { useEffect, useState } from 'react';
import { fetchAuthStatus, AuthStatus } from '../lib/api';

export function Setup() {
  const [status, setStatus] = useState<AuthStatus>({
    wakatime: false,
    googlefit: false,
    googlecalendar: false,
  });

  useEffect(() => {
    fetchAuthStatus().then(setStatus).catch(console.error);
  }, []);

  return (
    <div className="max-w-4xl mx-auto p-6">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-slate-900 dark:text-white mb-2">
          Setup & Integrations
        </h1>
        <p className="text-slate-500 dark:text-slate-400">
          Connect your accounts to start collecting data for your personal data lake.
        </p>
      </div>

      <div className="grid gap-6">
        {/* WakaTime */}
        <div className="bg-white dark:bg-slate-800 rounded-3xl p-6 shadow-sm border border-slate-100 dark:border-slate-700/50">
          <div className="flex items-start justify-between mb-4">
            <div className="flex items-center gap-3">
              <div className="p-3 bg-purple-100 dark:bg-purple-900/20 rounded-xl">
                <Code2 className="w-6 h-6 text-purple-600 dark:text-purple-400" />
              </div>
              <div>
                <h2 className="text-xl font-semibold text-slate-900 dark:text-white">WakaTime</h2>
                <p className="text-sm text-slate-500 dark:text-slate-400">Coding activity tracking</p>
              </div>
            </div>
            <div className={cn(
              "px-3 py-1 rounded-full text-sm font-medium flex items-center gap-2 border",
              status.wakatime 
                ? "bg-green-100 border-green-200 text-green-700 dark:bg-green-900/30 dark:border-green-800 dark:text-green-400" 
                : "bg-red-100 border-red-200 text-red-700 dark:bg-red-900/30 dark:border-red-800 dark:text-red-400"
            )}>
              {status.wakatime ? <Check className="w-3 h-3" /> : <X className="w-3 h-3" />}
              {status.wakatime ? "Connected" : "Not Connected"}
            </div>
          </div>
          
          <p className="text-slate-600 dark:text-slate-300 mb-6">
            Tracks time spent in your IDE and code editors. Requires a WakaTime account.
          </p>

          <a 
            href="https://wakatime.com/oauth/authorize?client_id=rtf2GDb7hmgMXGN7MsDupbSn&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback&response_type=code&state=wakatime"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 px-4 py-2 bg-slate-900 hover:bg-slate-800 dark:bg-white dark:hover:bg-slate-200 text-white dark:text-slate-900 rounded-xl font-medium transition-all shadow-md"
          >
            {status.wakatime ? "Reconnect WakaTime" : "Connect WakaTime"} <ExternalLink className="w-4 h-4" />
          </a>
        </div>

        {/* Google Fit */}
        <div className="bg-white dark:bg-slate-800 rounded-3xl p-6 shadow-sm border border-slate-100 dark:border-slate-700/50">
          <div className="flex items-start justify-between mb-4">
            <div className="flex items-center gap-3">
              <div className="p-3 bg-green-100 dark:bg-green-900/20 rounded-xl">
                <Activity className="w-6 h-6 text-green-600 dark:text-green-400" />
              </div>
              <div>
                <h2 className="text-xl font-semibold text-slate-900 dark:text-white">Google Fit</h2>
                <p className="text-sm text-slate-500 dark:text-slate-400">Physical activity & health</p>
              </div>
            </div>
            <div className={cn(
              "px-3 py-1 rounded-full text-sm font-medium flex items-center gap-2 border",
              status.googlefit 
                ? "bg-green-100 border-green-200 text-green-700 dark:bg-green-900/30 dark:border-green-800 dark:text-green-400" 
                : "bg-red-100 border-red-200 text-red-700 dark:bg-red-900/30 dark:border-red-800 dark:text-red-400"
            )}>
              {status.googlefit ? <Check className="w-3 h-3" /> : <X className="w-3 h-3" />}
              {status.googlefit ? "Connected" : "Not Connected"}
            </div>
          </div>
          
          <p className="text-slate-600 dark:text-slate-300 mb-6">
            Collects steps, distance, and sleep data from Google Fit.
          </p>

          <a 
            href="https://accounts.google.com/o/oauth2/v2/auth?access_type=offline&client_id=511018240872-sc6t6pctjkivoolcggr49v4nl3267dqa.apps.googleusercontent.com&include_granted_scopes=true&prompt=consent&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Foauth2callback&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Ffitness.activity.read+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Ffitness.body.read+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Ffitness.sleep.read+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Ffitness.location.read&state=googlefit"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 px-4 py-2 bg-slate-900 hover:bg-slate-800 dark:bg-white dark:hover:bg-slate-200 text-white dark:text-slate-900 rounded-xl font-medium transition-all shadow-md"
          >
            {status.googlefit ? "Reconnect Google Fit" : "Connect Google Fit"} <ExternalLink className="w-4 h-4" />
          </a>
        </div>

        {/* Google Calendar */}
        <div className="bg-white dark:bg-slate-800 rounded-3xl p-6 shadow-sm border border-slate-100 dark:border-slate-700/50">
          <div className="flex items-start justify-between mb-4">
            <div className="flex items-center gap-3">
              <div className="p-3 bg-blue-100 dark:bg-blue-900/20 rounded-xl">
                <Calendar className="w-6 h-6 text-blue-600 dark:text-blue-400" />
              </div>
              <div>
                <h2 className="text-xl font-semibold text-slate-900 dark:text-white">Google Calendar</h2>
                <p className="text-sm text-slate-500 dark:text-slate-400">Meetings & events</p>
              </div>
            </div>
            <div className={cn(
              "px-3 py-1 rounded-full text-sm font-medium flex items-center gap-2 border",
              status.googlecalendar 
                ? "bg-green-100 border-green-200 text-green-700 dark:bg-green-900/30 dark:border-green-800 dark:text-green-400" 
                : "bg-red-100 border-red-200 text-red-700 dark:bg-red-900/30 dark:border-red-800 dark:text-red-400"
            )}>
              {status.googlecalendar ? <Check className="w-3 h-3" /> : <X className="w-3 h-3" />}
              {status.googlecalendar ? "Connected" : "Not Connected"}
            </div>
          </div>
          
          <p className="text-slate-600 dark:text-slate-300 mb-6">
            Imports your calendar events to analyze time usage.
          </p>

          <a 
            href="https://accounts.google.com/o/oauth2/v2/auth?access_type=offline&client_id=511018240872-sc6t6pctjkivoolcggr49v4nl3267dqa.apps.googleusercontent.com&include_granted_scopes=true&prompt=consent&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Foauth2callback%2Fcalendar&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fcalendar.readonly+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fcalendar.events.readonly&state=googlecalendar"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 px-4 py-2 bg-slate-900 hover:bg-slate-800 dark:bg-white dark:hover:bg-slate-200 text-white dark:text-slate-900 rounded-xl font-medium transition-all shadow-md"
          >
            {status.googlecalendar ? "Reconnect Calendar" : "Connect Calendar"} <ExternalLink className="w-4 h-4" />
          </a>
        </div>
      </div>
    </div>
  );
}
