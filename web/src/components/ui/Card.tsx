import { HTMLMotionProps, motion } from "framer-motion";
import { cn } from "../../lib/utils";

interface CardProps extends HTMLMotionProps<"div"> {
  children: React.ReactNode;
  className?: string;
  noPadding?: boolean;
}

export function Card({ children, className, noPadding = false, ...props }: CardProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4 }}
      className={cn(
        "bg-surface/50 backdrop-blur-md border border-white/10 rounded-2xl shadow-xl overflow-hidden",
        !noPadding && "p-6",
        className
      )}
      {...props}
    >
      {children}
    </motion.div>
  );
}
