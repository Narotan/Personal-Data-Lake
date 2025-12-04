import { LayoutDashboard, Settings, BarChart2, PieChart, Activity, Database } from "lucide-react";
import { cn } from "../../lib/utils";

interface SidebarProps {
  currentView: 'dashboard' | 'setup' | 'auth-success';
  onViewChange: (view: 'dashboard' | 'setup') => void;
}

export function Sidebar({ currentView, onViewChange }: SidebarProps) {
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
    <aside className="hidden md:flex flex-col w-64 bg-slate-900/50 border-r border-white/5 min-h-screen p-4 sticky top-0 h-screen">
      <div className="flex items-center gap-3 px-2 mb-8 mt-2">
        <div className="p-2 bg-primary/20 rounded-lg">
          <LayoutDashboard className="w-6 h-6 text-primary" />
        </div>
        <div>
          <h1 className="text-lg font-bold tracking-tight text-white">Data Lake</h1>
          <p className="text-xs text-slate-400">Personal Analytics</p>
        </div>
      </div>

      <div className="space-y-1 flex-1">
        <p className="px-2 text-xs font-semibold text-slate-500 mb-2 uppercase tracking-wider">Menu</p>
        {menuItems.map((item) => (
          <button
            key={item.label}
            onClick={item.onClick}
            className={cn(
              "w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-200",
              (currentView === item.value) || (item.value === 'setup' && currentView === 'setup')
                ? "bg-primary text-white shadow-lg shadow-primary/20"
                : "text-slate-400 hover:text-white hover:bg-white/5"
            )}
          >
            <item.icon className="w-5 h-5" />
            {item.label}
          </button>
        ))}
      </div>

      <div className="mt-auto pt-4 border-t border-white/5">
        <button 
            onClick={() => onViewChange('setup')}
            className={cn(
                "w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-200",
                currentView === 'setup' 
                    ? "bg-white/10 text-white" 
                    : "text-slate-400 hover:text-white hover:bg-white/5"
            )}
        >
          <Settings className="w-5 h-5" />
          Settings
        </button>
      </div>
    </aside>
  );
}
