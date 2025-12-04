import { Bar, BarChart, Cell, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { Card } from "../ui/Card";
import { DailyFitStat } from "../../lib/api";
import { format, parseISO } from "date-fns";
import { Skeleton } from "../ui/Skeleton";

interface HealthChartProps {
  data: DailyFitStat[];
  loading: boolean;
}

const CustomTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    return (
      <div className="bg-slate-800 border border-slate-700 p-3 rounded-lg shadow-xl">
        <p className="text-slate-300 text-sm mb-1">{label}</p>
        <p className="text-emerald-400 font-bold text-lg">
          {payload[0].value.toLocaleString()}
          <span className="text-slate-500 text-xs font-normal ml-2">steps</span>
        </p>
      </div>
    );
  }
  return null;
};

export function HealthChart({ data, loading }: HealthChartProps) {
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

  const chartData = data.map(day => ({
    name: format(parseISO(day.date), 'MMM dd'),
    steps: day.steps,
  })).reverse();

  const hasData = data.length > 0;

  return (
    <Card className="min-h-[300px] flex flex-col bg-slate-900/50 backdrop-blur-sm border-slate-800">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h3 className="text-lg font-semibold text-white">Steps Activity</h3>
          <p className="text-xs text-slate-400 mt-1">Daily step count</p>
        </div>
        <span className="text-xs font-medium text-emerald-400 bg-emerald-400/10 px-2 py-1 rounded-full border border-emerald-400/20">
          Goal: 10k
        </span>
      </div>
      
      {!hasData ? (
        <div className="flex-1 flex items-center justify-center text-slate-500">
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
              stroke="#64748b" 
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
              tickFormatter={(value) => `${(value / 1000).toFixed(0)}k`}
            />
            <Tooltip content={<CustomTooltip />} cursor={{ fill: '#334155', opacity: 0.2 }} />
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
