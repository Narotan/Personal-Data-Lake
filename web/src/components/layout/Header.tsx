import { Calendar, ChevronDown, LayoutDashboard, Settings, X } from "lucide-react";
import { Button } from "../ui/Button";
import { cn } from "../../lib/utils";
import { useState, useEffect, useRef } from "react";


interface HeaderProps {
  dateRangeLabel: string;
  onRangeChange: (range: string) => void;
  currentRange: string;
  onSetupClick: () => void;
  isSetupMode: boolean;
  className?: string;
}

export function Header({ dateRangeLabel, onRangeChange, currentRange, onSetupClick, isSetupMode, className }: HeaderProps) {
  const [showCustomPicker, setShowCustomPicker] = useState(false);
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");
  const pickerRef = useRef<HTMLDivElement>(null);

  const ranges = [
    { label: 'Today', value: 'today' },
    { label: 'Yesterday', value: 'yesterday' },
    { label: '7 Days', value: '7d' },
    { label: '30 Days', value: '30d' },
    { label: 'Month', value: 'month' },
    { label: 'Year', value: 'year' },
    { label: 'All Time', value: 'all_time' },
    { label: 'Custom', value: 'custom' },
  ];

  useEffect(() => {
    if (currentRange.startsWith('custom:')) {
      const parts = currentRange.split(':');
      if (parts.length === 3) {
        setStartDate(parts[1]);
        setEndDate(parts[2]);
      }
    }
  }, [currentRange]);

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (pickerRef.current && !pickerRef.current.contains(event.target as Node)) {
        setShowCustomPicker(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  const handleRangeClick = (value: string) => {
    if (value === 'custom') {
      setShowCustomPicker(!showCustomPicker);
    } else {
      onRangeChange(value);
      setShowCustomPicker(false);
    }
  };

  const applyCustomRange = () => {
    if (startDate && endDate) {
      onRangeChange(`custom:${startDate}:${endDate}`);
      setShowCustomPicker(false);
    }
  };

  const isCustomActive = currentRange.startsWith('custom') || showCustomPicker;

  return (
    <header className={cn("flex flex-col md:flex-row items-start md:items-center justify-between gap-4 mb-8 relative", className)}>
      <div className="flex items-center gap-3 md:hidden">
        <div className="p-2 bg-primary/20 rounded-lg">
          <LayoutDashboard className="w-6 h-6 text-primary" />
        </div>
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-white">Data Lake</h1>
        </div>
      </div>
      
      <div className="hidden md:block">
        <h2 className="text-2xl font-bold tracking-tight text-white">
            {isSetupMode ? "System Setup" : "Dashboard Overview"}
        </h2>
        <p className="text-sm text-slate-400">
            {isSetupMode ? "Configure your data sources" : "Welcome back, here's what's happening"}
        </p>
      </div>

      <div className="flex flex-col items-end gap-2 w-full md:w-auto">
        <div className="flex items-center gap-4 w-full md:w-auto justify-between md:justify-end">
            {!isSetupMode && (
            <div className="flex items-center gap-2 bg-surface/50 p-1 rounded-xl border border-white/5 backdrop-blur-sm overflow-x-auto max-w-full relative">
                {ranges.map((range) => (
                <button
                    key={range.value}
                    onClick={() => handleRangeClick(range.value)}
                    className={cn(
                    "px-3 py-1.5 text-sm font-medium rounded-lg transition-all whitespace-nowrap",
                    (range.value === 'custom' ? isCustomActive : currentRange === range.value)
                        ? "bg-primary text-white shadow-lg shadow-primary/25"
                        : "text-slate-400 hover:text-white hover:bg-white/5"
                    )}
                >
                    {range.label}
                </button>
                ))}
            </div>
            )}
        </div>

        {showCustomPicker && (
            <div ref={pickerRef} className="absolute top-full right-0 mt-2 p-4 bg-surface border border-white/10 rounded-xl shadow-xl z-50 flex flex-col gap-4 w-72 backdrop-blur-xl">
                <div className="flex justify-between items-center">
                    <h3 className="text-sm font-medium text-white">Custom Range</h3>
                    <button onClick={() => setShowCustomPicker(false)} className="text-slate-400 hover:text-white">
                        <X className="w-4 h-4" />
                    </button>
                </div>
                <div className="space-y-2">
                    <div className="flex flex-col gap-1">
                        <label className="text-xs text-slate-400">Start Date</label>
                        <input 
                            type="date" 
                            value={startDate}
                            onChange={(e) => setStartDate(e.target.value)}
                            className="bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-primary"
                        />
                    </div>
                    <div className="flex flex-col gap-1">
                        <label className="text-xs text-slate-400">End Date</label>
                        <input 
                            type="date" 
                            value={endDate}
                            onChange={(e) => setEndDate(e.target.value)}
                            className="bg-white/5 border border-white/10 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-primary"
                        />
                    </div>
                </div>
                <Button onClick={applyCustomRange} disabled={!startDate || !endDate} className="w-full">
                    Apply Range
                </Button>
            </div>
        )}
      </div>
    </header>
  );
}
