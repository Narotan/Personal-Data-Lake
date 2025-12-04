import { Card } from "../ui/Card";
import { CalendarEvent } from "../../lib/api";
import { format, parseISO, differenceInMinutes, isWithinInterval } from "date-fns";
import { Skeleton } from "../ui/Skeleton";

interface ScheduleTimelineProps {
  data: CalendarEvent[];
  loading: boolean;
}

export function ScheduleTimeline({ data, loading }: ScheduleTimelineProps) {
  if (loading) {
    return (
      <Card className="min-h-[300px] flex flex-col">
        <h3 className="text-lg font-semibold text-white mb-6">Today's Schedule</h3>
        <div className="space-y-4">
          {[1, 2, 3].map((i) => (
            <div key={i} className="flex items-start gap-4">
              <Skeleton className="w-14 h-4" />
              <Skeleton className="w-3 h-3 rounded-full mt-1" />
              <Skeleton className="flex-1 h-16 rounded-lg" />
            </div>
          ))}
        </div>
      </Card>
    );
  }

  const sortedEvents = [...data].sort((a, b) => 
    new Date(a.start_time).getTime() - new Date(b.start_time).getTime()
  );

  const now = new Date();

  return (
    <Card className="min-h-[300px] flex flex-col">
      <h3 className="text-lg font-semibold text-slate-900 dark:text-white mb-6">Schedule</h3>
      
      <div className="space-y-6 relative pl-2">
        {/* Vertical line */}
        <div className="absolute left-[4.5rem] top-2 bottom-2 w-px bg-slate-200 dark:bg-slate-700" />
        
        {sortedEvents.length === 0 ? (
          <div className="text-slate-500 dark:text-slate-400 text-center py-8">No events scheduled</div>
        ) : (
          sortedEvents.map((event, index) => {
            const start = parseISO(event.start_time);
            const end = parseISO(event.end_time);
            const durationMinutes = differenceInMinutes(end, start);
            const duration = durationMinutes >= 60 
              ? `${Math.floor(durationMinutes / 60)}h ${durationMinutes % 60 > 0 ? ` ${durationMinutes % 60}m` : ''}`
              : `${durationMinutes}m`;

            const isCurrent = isWithinInterval(now, { start, end });
            const isPast = now > end;

            return (
              <div key={event.id || index} className={`flex items-start gap-4 relative z-10 group ${isPast ? 'opacity-50' : ''}`}>
                <div className="w-14 text-sm text-slate-500 dark:text-slate-400 font-mono text-right pt-1">
                  {format(start, 'HH:mm')}
                </div>
                
                <div className={`w-3 h-3 rounded-full mt-2 border-2 border-white dark:border-slate-800 transition-colors ${isCurrent ? 'bg-blue-500 shadow-[0_0_10px_rgba(59,130,246,0.5)]' : 'bg-slate-300 dark:bg-slate-600'}`} />
                
                <div className={`flex-1 p-3 rounded-xl border transition-all ${
                    isCurrent 
                        ? 'border-blue-200 dark:border-blue-800 bg-blue-50 dark:bg-blue-900/20' 
                        : 'border-slate-100 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/50 hover:bg-slate-100 dark:hover:bg-slate-800'
                }`}>
                  <div className={`font-medium text-sm ${isCurrent ? 'text-blue-900 dark:text-blue-100' : 'text-slate-900 dark:text-white'}`}>
                    {event.summary}
                  </div>
                  <div className="text-xs text-slate-500 dark:text-slate-400 mt-1 flex justify-between">
                    <span>{duration}</span>
                    {isCurrent && <span className="text-blue-600 dark:text-blue-400 animate-pulse">Now</span>}
                  </div>
                </div>
              </div>
            );
          })
        )}
      </div>
    </Card>
  );
}
