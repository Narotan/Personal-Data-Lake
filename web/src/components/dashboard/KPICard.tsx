import { LucideIcon } from "lucide-react";
import { Card } from "../ui/Card";
import { cn } from "../../lib/utils";

interface KPICardProps {
  title: string;
  value: string | number;
  subtitle?: string;
  icon: LucideIcon;
  trend?: {
    value: number;
    isPositive: boolean;
  };
  color: "wakatime" | "googlefit" | "calendar" | "activity";
}

export function KPICard({ title, value, subtitle, icon: Icon, trend, color }: KPICardProps) {
  const colorStyles = {
    wakatime: "text-accent-wakatime bg-accent-wakatime/10",
    googlefit: "text-accent-googlefit bg-accent-googlefit/10",
    calendar: "text-accent-calendar bg-accent-calendar/10",
    activity: "text-accent-activity bg-accent-activity/10",
  };

  return (
    <Card className="relative overflow-hidden group transition-all duration-300 hover:shadow-lg hover:-translate-y-1 bg-white dark:bg-slate-800 border border-slate-100 dark:border-slate-700/50">
      <div className="flex justify-between items-start mb-4">
        <div className={cn("p-2 rounded-xl transition-transform duration-300 group-hover:scale-110", colorStyles[color])}>
          <Icon className="w-5 h-5" />
        </div>
        {trend && (
          <div className={cn(
            "text-xs font-medium px-2 py-1 rounded-full",
            trend.isPositive 
                ? "text-emerald-600 dark:text-emerald-400 bg-emerald-100 dark:bg-emerald-900/30" 
                : "text-rose-600 dark:text-rose-400 bg-rose-100 dark:bg-rose-900/30"
          )}>
            {trend.isPositive ? "+" : ""}{trend.value}%
          </div>
        )}
      </div>
      
      <div>
        <h3 className="text-slate-500 dark:text-slate-400 text-sm font-medium mb-1">{title}</h3>
        <div className="text-2xl font-bold text-slate-900 dark:text-white tracking-tight">{value}</div>
        {subtitle && <div className="text-xs text-slate-400 dark:text-slate-500 mt-1">{subtitle}</div>}
      </div>

      {/* Background Icon */}
      <Icon className={cn(
        "absolute -right-6 -bottom-6 w-32 h-32 opacity-[0.03] dark:opacity-[0.05] transform -rotate-12 transition-transform duration-500 group-hover:scale-110 pointer-events-none",
        color === 'wakatime' && "text-accent-wakatime",
        color === 'googlefit' && "text-accent-googlefit",
        color === 'calendar' && "text-accent-calendar",
        color === 'activity' && "text-accent-activity",
      )} />

      {/* Decorative gradient glow */}
      <div className={cn(
        "absolute -right-4 -bottom-4 w-24 h-24 rounded-full blur-3xl opacity-0 group-hover:opacity-10 transition-opacity pointer-events-none",
        color === 'wakatime' && "bg-accent-wakatime",
        color === 'googlefit' && "bg-accent-googlefit",
        color === 'calendar' && "bg-accent-calendar",
        color === 'activity' && "bg-accent-activity",
      )} />
    </Card>
  );
}
