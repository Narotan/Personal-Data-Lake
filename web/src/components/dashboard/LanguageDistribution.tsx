import { Cell, Legend, Pie, PieChart, ResponsiveContainer, Tooltip } from "recharts";
import { Card } from "../ui/Card";
import { AggregatedLanguageStat } from "../../lib/api";
import { Skeleton } from "../ui/Skeleton";
import { useState } from "react";

interface LanguageDistributionProps {
  data: AggregatedLanguageStat[];
  loading: boolean;
}

const COLORS = ['#3b82f6', '#8b5cf6', '#10b981', '#f59e0b', '#ef4444', '#ec4899'];

const formatDuration = (seconds: number) => {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  if (hours === 0) return `${minutes}m`;
  return `${hours}h ${minutes}m`;
};

const CustomTooltip = ({ active, payload }: any) => {
  if (active && payload && payload.length) {
    const data = payload[0].payload;
    return (
      <div className="bg-slate-800 border border-slate-700 p-3 rounded-lg shadow-xl">
        <p className="text-slate-300 text-sm mb-1">{data.name}</p>
        <p className="text-white font-bold">
          {formatDuration(data.total_seconds)}
          <span className="text-slate-500 text-xs font-normal ml-2">({data.percent.toFixed(1)}%)</span>
        </p>
      </div>
    );
  }
  return null;
};

export function LanguageDistribution({ data, loading }: LanguageDistributionProps) {
  const [activeIndex, setActiveIndex] = useState<number | undefined>();

  if (loading) {
    return (
      <Card className="min-h-[300px] flex flex-col">
        <div className="flex items-center justify-between mb-6">
          <Skeleton className="h-6 w-48" />
        </div>
        <div className="flex-1 flex items-center justify-center">
           <Skeleton className="w-40 h-40 rounded-full" />
        </div>
      </Card>
    );
  }

  const hasData = data.length > 0;
  const totalSeconds = data.reduce((acc, curr) => acc + curr.total_seconds, 0);
  const totalDuration = formatDuration(totalSeconds);

  const onPieEnter = (_: any, index: number) => {
    setActiveIndex(index);
  };

  const onPieLeave = () => {
    setActiveIndex(undefined);
  };

  return (
    <Card className="min-h-[300px] flex flex-col bg-slate-900/50 backdrop-blur-sm border-slate-800">
      <h3 className="text-lg font-semibold text-white mb-6">Top Languages</h3>
      
      {!hasData ? (
        <div className="flex-1 flex items-center justify-center text-slate-500">
          <p>No language data available</p>
        </div>
      ) : (
        <div className="flex-1 w-full min-h-[200px] relative">
          <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
            <div className="text-center">
              <p className="text-2xl font-bold text-white">{totalDuration}</p>
              <p className="text-xs text-slate-400">Total</p>
            </div>
          </div>
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={data}
                cx="50%"
                cy="50%"
                innerRadius={60}
                outerRadius={80}
                paddingAngle={0}
                dataKey="total_seconds"
                onMouseEnter={onPieEnter}
                onMouseLeave={onPieLeave}
              >
                {data.map((_, index) => (
                  <Cell 
                    key={`cell-${index}`} 
                    fill={COLORS[index % COLORS.length]} 
                    stroke="none"
                    opacity={activeIndex === undefined || activeIndex === index ? 1 : 0.6}
                    style={{ transition: 'opacity 0.3s ease' }}
                  />
                ))}
              </Pie>
              <Tooltip content={<CustomTooltip />} />
              <Legend 
                verticalAlign="bottom" 
                height={36}
                formatter={(value) => {
                    return <span className="text-slate-300 ml-2 text-xs">{value}</span>;
                }}
              />
            </PieChart>
          </ResponsiveContainer>
        </div>
      )}
    </Card>
  );
}
