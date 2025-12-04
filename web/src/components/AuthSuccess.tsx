import { CheckCircle2, ArrowRight } from 'lucide-react';
import { Button } from './ui/Button';

interface AuthSuccessProps {
  onContinue: () => void;
}

export function AuthSuccess({ onContinue }: AuthSuccessProps) {
  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] text-center p-6">
      <div className="w-20 h-20 bg-green-500/20 rounded-full flex items-center justify-center mb-6 animate-in zoom-in duration-300">
        <CheckCircle2 className="w-10 h-10 text-green-400" />
      </div>
      
      <h1 className="text-3xl font-bold text-white mb-4">Authorization Successful!</h1>
      <p className="text-slate-400 max-w-md mb-8">
        Your account has been successfully connected. Data collection will start automatically in the background.
      </p>

      <Button 
        onClick={onContinue}
        className="gap-2 px-8 py-6 text-lg"
      >
        Continue to Setup <ArrowRight className="w-5 h-5" />
      </Button>
    </div>
  );
}
