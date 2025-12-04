import { Area, AreaChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { Card } from "../ui/Card";
import { DailyStat } from "../../lib/api";
import { format, parseISO, differenceInDays } from "date-fns";
import { Skeleton } from "../ui/Skeleton";
import { aggregateByMonth } from "../../lib/utils";
import { useTheme } from "../../lib/theme";

interface ProductivityChartProps {
  data: DailyStat[];
  loading: boolean;
}

const CustomTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    return (
      <div className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 p-3 rounded-xl shadow-xl">
        <p className="text-slate-500 dark:text-slate-400 text-sm mb-1">{label}</p>
        <p className="text-purple-600 dark:text-purple-400 font-bold text-lg">
          {payload[0].value}h
          <span className="text-slate-400 text-xs font-normal ml-2">coding</span>
        </p>
      </div>
    );
  }
  return null;
};

export function ProductivityChart({ data, loading }: ProductivityChartProps) {
  const { theme } = useTheme();
  const isDark = theme === 'dark';

  if (loading) {
    return (
      <Card className="col-span-2 min-h-[300px] flex flex-col">
        <div className="flex items-center justify-between mb-6">
          <Skeleton className="h-6 w-48" />
        </div>
        <Skeleton className="flex-1 w-full rounded-lg" />
      </Card>
    );
  }

  // Определяем нужно ли агрегировать по месяцам
  const shouldAggregateByMonth = data.length > 1 && 
    differenceInDays(
      parseISO(data[data.length - 1]?.date || format(new Date(), 'yyyy-MM-dd')), 
      parseISO(data[0]?.date || format(new Date(), 'yyyy-MM-dd'))
    ) > 90;
  
  let chartData;
  if (shouldAggregateByMonth) {
    const aggregated = aggregateByMonth(data);
    chartData = aggregated.map(month => ({
      name: format(parseISO(month.date), 'MMM yyyy'),
      coding: Number((month.total_seconds / 3600).toFixed(2)),
    }));
  } else {
    chartData = data.map(day => ({
      name: format(parseISO(day.date), 'MMM dd'),
      coding: Number((day.total_seconds / 3600).toFixed(2)),
    }));
  }

  const hasData = data.length > 0;

  return (
    <Card className="col-span-2 min-h-[300px] flex flex-col">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h3 className="text-lg font-semibold text-slate-900 dark:text-white">Productivity Trends</h3>
          <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">Daily coding activity over time</p>
        </div>
        <div className="flex gap-2">
          <div className="flex items-center gap-2 text-xs text-slate-500 dark:text-slate-400 bg-slate-100 dark:bg-slate-800 px-2 py-1 rounded-full">
            <div className="w-2 h-2 rounded-full bg-purple-500 shadow-[0_0_8px_rgba(168,85,247,0.5)]" /> Coding
          </div>
        </div>
      </div>
      
      {!hasData ? (
        <div className="flex-1 flex items-center justify-center text-slate-500 dark:text-slate-400">
          <div className="text-center">
            <p className="text-lg mb-2">📊 No coding data available</p>
            <p className="text-sm">Data will appear here once WakaTime starts collecting.</p>
          </div>
        </div>
      ) : (
        <div className="w-full h-[250px]">
          <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={chartData} margin={{ top: 10, right: 10, left: 0, bottom: 0 }}>
            <defs>
              <linearGradient id="colorCoding" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#8b5cf6" stopOpacity={0.4}/>
                <stop offset="95%" stopColor="#8b5cf6" stopOpacity={0}/>
              </linearGradient>
            </defs>
            <CartesianGrid strokeDasharray="3 3" stroke={isDark ? "#334155" : "#e2e8f0"} vertical={false} opacity={0.8} />
            <XAxis 
              dataKey="name" 
              stroke={isDark ? "#94a3b8" : "#94a3b8"} 
              fontSize={12} 
              tickLine={false} 
              axisLine={false}
              dy={10}
            />
            <YAxis 
              stroke="#64748b" 
              fontSize={12} 
              tickLine={false} 
              axisLine={false} 
              tickFormatter={(value) => `${value}h`}
              dx={-10}
            />
            <Tooltip content={<CustomTooltip />} cursor={{ stroke: '#8b5cf6', strokeWidth: 1, strokeDasharray: '4 4' }} />
            <Area 
              type="monotone" 
              dataKey="coding" 
              stroke="#8b5cf6" 
              strokeWidth={3}
              fillOpacity={1} 
              fill="url(#colorCoding)" 
              activeDot={{ r: 6, strokeWidth: 0, fill: '#fff' }}
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
      )}
    </Card>
  );
}
