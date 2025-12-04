import { LayoutDashboard, Settings, BarChart2, PieChart, Activity, Database, Moon, Sun } from "lucide-react";
import { cn } from "../../lib/utils";
import { useTheme } from "../../lib/theme";

interface SidebarProps {
  currentView: 'dashboard' | 'setup' | 'auth-success';
  onViewChange: (view: 'dashboard' | 'setup') => void;
}

export function Sidebar({ currentView, onViewChange }: SidebarProps) {
  const { theme, setTheme } = useTheme();

  const menuItems = [
    { 
      icon: LayoutDashboard, 
      label: "Dashboard", 
      value: 'dashboard' as const,
      onClick: () => onViewChange('dashboard')
    },
    { 
      icon: BarChart2, 
      label: "Analytics", 
      value: 'analytics', // Placeholder
      onClick: () => {} 
    },
    { 
      icon: Activity, 
      label: "Activity", 
      value: 'activity', // Placeholder
      onClick: () => {} 
    },
    { 
      icon: Database, 
      label: "Data Sources", 
      value: 'setup' as const,
      onClick: () => onViewChange('setup')
    },
  ];

  return (
    <aside className="hidden md:flex flex-col w-64 bg-white/30 dark:bg-slate-900/30 backdrop-blur-xl border-r border-slate-200/50 dark:border-slate-800/50 min-h-screen p-4 sticky top-0 h-screen z-20">
      <div className="flex items-center gap-3 px-2 mb-8 mt-2">
        <div className="p-2 bg-white dark:bg-slate-800 rounded-xl shadow-sm transition-colors ring-1 ring-slate-900/5 dark:ring-white/10">
          <LayoutDashboard className="w-6 h-6 text-slate-900 dark:text-white" />
        </div>
        <div>
          <h1 className="text-lg font-bold tracking-tight text-slate-900 dark:text-white">Data Lake</h1>
          <p className="text-xs text-slate-500 dark:text-slate-400">Personal Analytics</p>
        </div>
      </div>

      <div className="space-y-1 flex-1">
        <p className="px-2 text-xs font-semibold text-slate-400 dark:text-slate-500 mb-2 uppercase tracking-wider">Menu</p>
        {menuItems.map((item) => (
          <button
            key={item.label}
            onClick={item.onClick}
            className={cn(
              "w-full flex items-center gap-3 px-3 py-2.5 rounded-2xl text-sm font-medium transition-all duration-200 group relative",
              (currentView === item.value) || (item.value === 'setup' && currentView === 'setup')
                ? "bg-slate-900 dark:bg-white text-white dark:text-slate-900 shadow-lg shadow-slate-900/20 dark:shadow-white/20"
                : "text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white hover:bg-white/50 dark:hover:bg-white/5"
            )}
          >
            <item.icon className="w-5 h-5 transition-transform group-hover:scale-110 duration-200" />
            {item.label}
            {(currentView === item.value || (item.value === 'setup' && currentView === 'setup')) && (
                <div className="absolute right-2 w-1.5 h-1.5 rounded-full bg-white dark:bg-slate-900 animate-pulse" />
            )}
          </button>
        ))}
      </div>

      <div className="mt-auto pt-4 space-y-2">
        <button 
            onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
            className="w-full flex items-center gap-3 px-3 py-2.5 rounded-2xl text-sm font-medium transition-all duration-200 text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white hover:bg-white/50 dark:hover:bg-white/5 group"
        >
          {theme === 'dark' ? <Sun className="w-5 h-5 transition-transform group-hover:rotate-90 duration-500" /> : <Moon className="w-5 h-5 transition-transform group-hover:-rotate-12 duration-300" />}
          {theme === 'dark' ? 'Light Mode' : 'Dark Mode'}
        </button>

        <button 
            onClick={() => onViewChange('setup')}
            className={cn(
                "w-full flex items-center gap-3 px-3 py-2.5 rounded-2xl text-sm font-medium transition-all duration-200 group",
                currentView === 'setup' 
                    ? "bg-slate-900 dark:bg-white text-white dark:text-slate-900 shadow-lg shadow-slate-900/20 dark:shadow-white/20" 
                    : "text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white hover:bg-white/50 dark:hover:bg-white/5"
            )}
        >
          <Settings className="w-5 h-5 transition-transform group-hover:rotate-90 duration-500" />
          Settings
        </button>
      </div>
    </aside>
  );
}
