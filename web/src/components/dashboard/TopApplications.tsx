import { Bar, BarChart, ResponsiveContainer, Tooltip, XAxis, YAxis, Cell } from "recharts";
import { Card } from "../ui/Card";
import { AppStat } from "../../lib/api";
import { Skeleton } from "../ui/Skeleton";

interface TopApplicationsProps {
  data: AppStat[];
  loading: boolean;
}

const CustomTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    return (
      <div className="bg-slate-800 border border-slate-700 p-3 rounded-lg shadow-xl">
        <p className="text-slate-300 text-sm mb-1">{label}</p>
        <p className="text-blue-400 font-bold text-lg">
          {payload[0].value}h
          <span className="text-slate-500 text-xs font-normal ml-2">active</span>
        </p>
      </div>
    );
  }
  return null;
};

export function TopApplications({ data, loading }: TopApplicationsProps) {
  if (loading) {
    return (
      <Card className="min-h-[300px] flex flex-col">
        <div className="flex items-center justify-between mb-6">
          <Skeleton className="h-6 w-40" />
        </div>
        <Skeleton className="flex-1 w-full rounded-lg" />
      </Card>
    );
  }

  // Sort by duration and take top 5
  const chartData = [...data]
    .sort((a, b) => b.TotalDuration - a.TotalDuration)
    .slice(0, 5)
    .map(app => ({
      name: app.App,
      duration: Number((app.TotalDuration / 3600).toFixed(2)), // hours
    }));

  return (
    <Card className="min-h-[300px] flex flex-col bg-slate-900/50 backdrop-blur-sm border-slate-800">
      <h3 className="text-lg font-semibold text-white mb-6">Top Applications</h3>
      
      {chartData.length === 0 ? (
        <div className="flex-1 flex items-center justify-center text-slate-500">
          <p>No application data available</p>
        </div>
      ) : (
        <div className="w-full h-[250px]">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={chartData} layout="vertical" margin={{ left: 0, right: 20 }}>
              <XAxis type="number" hide />
              <YAxis 
                dataKey="name" 
                type="category" 
                width={100}
                tick={{ fill: '#94a3b8', fontSize: 12 }}
                tickLine={false}
                axisLine={false}
              />
              <Tooltip content={<CustomTooltip />} cursor={{ fill: '#334155', opacity: 0.2 }} />
              <Bar dataKey="duration" radius={[0, 4, 4, 0]} barSize={24}>
                {chartData.map((_, index) => (
                  <Cell 
                    key={`cell-${index}`} 
                    fill="#3b82f6" 
                    fillOpacity={0.8}
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
