import { Area, AreaChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { Card } from "../ui/Card";
import { DailyStat } from "../../lib/api";
import { format, parseISO } from "date-fns";
import { Skeleton } from "../ui/Skeleton";

interface ProductivityChartProps {
  data: DailyStat[];
  loading: boolean;
}

const CustomTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    return (
      <div className="bg-slate-800 border border-slate-700 p-3 rounded-lg shadow-xl">
        <p className="text-slate-300 text-sm mb-1">{label}</p>
        <p className="text-purple-400 font-bold text-lg">
          {payload[0].value}h
          <span className="text-slate-500 text-xs font-normal ml-2">coding</span>
        </p>
      </div>
    );
  }
  return null;
};

export function ProductivityChart({ data, loading }: ProductivityChartProps) {
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

  const chartData = data.map(day => ({
    name: format(parseISO(day.date), 'MMM dd'),
    coding: Number((day.total_seconds / 3600).toFixed(2)),
  })).reverse();

  const hasData = data.length > 0;

  return (
    <Card className="col-span-2 min-h-[300px] flex flex-col bg-slate-900/50 backdrop-blur-sm border-slate-800">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h3 className="text-lg font-semibold text-white">Productivity Trends</h3>
          <p className="text-xs text-slate-400 mt-1">Daily coding activity over time</p>
        </div>
        <div className="flex gap-2">
          <div className="flex items-center gap-2 text-xs text-slate-400 bg-slate-800/50 px-2 py-1 rounded-full">
            <div className="w-2 h-2 rounded-full bg-purple-500 shadow-[0_0_8px_rgba(168,85,247,0.5)]" /> Coding
          </div>
        </div>
      </div>
      
      {!hasData ? (
        <div className="flex-1 flex items-center justify-center text-slate-500">
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
            <CartesianGrid strokeDasharray="3 3" stroke="#334155" vertical={false} opacity={0.3} />
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
              activeDot={{ r: 6, strokeWidth: 0, fill: '#fff', shadow: '0 0 10px #8b5cf6' }}
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
      )}
    </Card>
  );
}
