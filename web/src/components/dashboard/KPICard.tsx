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
    <Card className="relative overflow-hidden group hover:border-white/20 transition-colors">
      <div className="flex justify-between items-start mb-4">
        <div className={cn("p-2 rounded-lg", colorStyles[color])}>
          <Icon className="w-5 h-5" />
        </div>
        {trend && (
          <div className={cn(
            "text-xs font-medium px-2 py-1 rounded-full",
            trend.isPositive ? "text-emerald-400 bg-emerald-400/10" : "text-rose-400 bg-rose-400/10"
          )}>
            {trend.isPositive ? "+" : ""}{trend.value}%
          </div>
        )}
      </div>
      
      <div>
        <h3 className="text-slate-400 text-sm font-medium mb-1">{title}</h3>
        <div className="text-2xl font-bold text-white tracking-tight">{value}</div>
        {subtitle && <div className="text-xs text-slate-500 mt-1">{subtitle}</div>}
      </div>

      {/* Decorative gradient glow */}
      <div className={cn(
        "absolute -right-4 -bottom-4 w-24 h-24 rounded-full blur-3xl opacity-0 group-hover:opacity-20 transition-opacity pointer-events-none",
        color === 'wakatime' && "bg-accent-wakatime",
        color === 'googlefit' && "bg-accent-googlefit",
        color === 'calendar' && "bg-accent-calendar",
        color === 'activity' && "bg-accent-activity",
      )} />
    </Card>
  );
}
