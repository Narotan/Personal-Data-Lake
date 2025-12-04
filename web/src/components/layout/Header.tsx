import { Calendar, ChevronDown, LayoutDashboard, Settings } from "lucide-react";
import { Button } from "../ui/Button";
import { cn } from "../../lib/utils";

interface HeaderProps {
  dateRangeLabel: string;
  onRangeChange: (range: string) => void;
  currentRange: string;
  onSetupClick: () => void;
  isSetupMode: boolean;
}

export function Header({ dateRangeLabel, onRangeChange, currentRange, onSetupClick, isSetupMode }: HeaderProps) {
  const ranges = [
    { label: 'Today', value: 'today' },
    { label: 'Yesterday', value: 'yesterday' },
    { label: '7 Days', value: '7d' },
    { label: '30 Days', value: '30d' },
    { label: 'Month', value: 'month' },
  ];

  return (
    <header className="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 mb-8">
      <div className="flex items-center gap-3 cursor-pointer" onClick={() => isSetupMode && onSetupClick()}>
        <div className="p-2 bg-primary/20 rounded-lg">
          <LayoutDashboard className="w-6 h-6 text-primary" />
        </div>
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-white">Personal Data Lake</h1>
          <p className="text-sm text-slate-400">Analytics Dashboard</p>
        </div>
      </div>

      <div className="flex items-center gap-4">
        {!isSetupMode && (
          <div className="flex items-center gap-2 bg-surface/50 p-1 rounded-xl border border-white/5 backdrop-blur-sm">
            {ranges.map((range) => (
              <button
                key={range.value}
                onClick={() => onRangeChange(range.value)}
                className={cn(
                  "px-3 py-1.5 text-sm font-medium rounded-lg transition-all",
                  currentRange === range.value
                    ? "bg-primary text-white shadow-lg shadow-primary/25"
                    : "text-slate-400 hover:text-white hover:bg-white/5"
                )}
              >
                {range.label}
              </button>
            ))}
          </div>
        )}
        
        <Button 
          variant={isSetupMode ? "primary" : "ghost"} 
          size="sm" 
          className="gap-2"
          onClick={onSetupClick}
        >
          <Settings className="w-4 h-4" />
          <span>{isSetupMode ? "Back to Dashboard" : "Setup"}</span>
        </Button>
      </div>
    </header>
  );
}
