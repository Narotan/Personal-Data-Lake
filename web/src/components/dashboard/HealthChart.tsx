import { Bar, BarChart, Cell, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { Card } from "../ui/Card";
import { DailyFitStat } from "../../lib/api";
import { format, parseISO, differenceInDays } from "date-fns";
import { Skeleton } from "../ui/Skeleton";
import { aggregateByMonth } from "../../lib/utils";
import { useTheme } from "../../lib/theme";

interface HealthChartProps {
  data: DailyFitStat[];
  loading: boolean;
}

const CustomTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    return (
      <div className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 p-3 rounded-xl shadow-xl">
        <p className="text-slate-500 dark:text-slate-400 text-sm mb-1">{label}</p>
        <p className="text-emerald-600 dark:text-emerald-400 font-bold text-lg">
          {payload[0].value.toLocaleString()}
          <span className="text-slate-400 text-xs font-normal ml-2">steps</span>
        </p>
      </div>
    );
  }
  return null;
};

export function HealthChart({ data, loading }: HealthChartProps) {
  const { theme } = useTheme();
  const isDark = theme === 'dark';

  if (loading) {
    return (
      <Card className="min-h-[300px] flex flex-col">
        <div className="flex items-center justify-between mb-6">
          <Skeleton className="h-6 w-32" />
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
      steps: Math.round(month.steps),
    }));
  } else {
    chartData = data.map(day => ({
      name: format(parseISO(day.date), 'MMM dd'),
      steps: day.steps,
    }));
  }

  const hasData = data.length > 0;

  return (
    <Card className="min-h-[300px] flex flex-col">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h3 className="text-lg font-semibold text-slate-900 dark:text-white">Steps Activity</h3>
          <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">Daily step count</p>
        </div>
        <span className="text-xs font-medium text-emerald-700 dark:text-emerald-400 bg-emerald-100 dark:bg-emerald-900/30 px-2 py-1 rounded-full border border-emerald-200 dark:border-emerald-800">
          Goal: 10k
        </span>
      </div>
      
      {!hasData ? (
        <div className="flex-1 flex items-center justify-center text-slate-500 dark:text-slate-400">
          <div className="text-center">
            <p className="text-lg mb-2">👟 No activity data</p>
            <p className="text-sm">Connect Google Fit to track your steps.</p>
          </div>
        </div>
      ) : (
        <div className="w-full h-[250px]">
          <ResponsiveContainer width="100%" height="100%">
          <BarChart data={chartData} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
            <XAxis 
              dataKey="name" 
              stroke={isDark ? "#94a3b8" : "#94a3b8"} 
              fontSize={12} 
              tickLine={false} 
              axisLine={false}
              dy={10}
            />
            <YAxis 
              stroke={isDark ? "#94a3b8" : "#94a3b8"} 
              fontSize={12} 
              tickLine={false} 
              axisLine={false} 
              tickFormatter={(value) => `${(value / 1000).toFixed(0)}k`}
            />
            <Tooltip content={<CustomTooltip />} cursor={{ fill: isDark ? '#334155' : '#f1f5f9', opacity: 0.2 }} />
            <Bar dataKey="steps" radius={[4, 4, 0, 0]}>
              {chartData.map((entry, index) => (
                <Cell 
                  key={`cell-${index}`} 
                  fill={entry.steps >= 10000 ? '#10b981' : '#3b82f6'} 
                  fillOpacity={entry.steps >= 10000 ? 1 : 0.6}
                />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </div>
      )}
    </Card>
  );
}
