import { Card } from "../ui/Card";
import { AggregatedProjectStat } from "../../lib/api";
import { Skeleton } from "../ui/Skeleton";
import { formatDuration } from "../../lib/utils";

interface TopProjectsProps {
  data: AggregatedProjectStat[];
  loading: boolean;
}

export function TopProjects({ data, loading }: TopProjectsProps) {
  if (loading) {
    return (
      <Card className="min-h-[300px] flex flex-col">
        <div className="flex items-center justify-between mb-6">
          <Skeleton className="h-6 w-32" />
        </div>
        <div className="space-y-4">
            {[1, 2, 3, 4].map(i => <Skeleton key={i} className="h-12 w-full" />)}
        </div>
      </Card>
    );
  }

  const maxTime = Math.max(...data.map(p => p.total_seconds), 1);

  return (
    <Card className="min-h-[300px] flex flex-col">
      <h3 className="text-lg font-semibold text-white mb-6">Top Projects</h3>
      
      <div className="space-y-4 overflow-y-auto max-h-[300px] pr-2 custom-scrollbar">
        {data.length === 0 ? (
             <div className="text-center text-slate-500 py-8">No project data available</div>
        ) : (
            data.map((project) => (
            <div key={project.name} className="space-y-2">
                <div className="flex justify-between text-sm">
                <span className="text-slate-200 font-medium truncate max-w-[70%]">{project.name}</span>
                <span className="text-slate-400">{formatDuration(project.total_seconds)}</span>
                </div>
                <div className="h-2 bg-slate-700 rounded-full overflow-hidden">
                <div 
                    className="h-full bg-blue-500 rounded-full transition-all duration-500 ease-out"
                    style={{ width: `${(project.total_seconds / maxTime) * 100}%` }}
                />
                </div>
            </div>
            ))
        )}
      </div>
    </Card>
  );
}
