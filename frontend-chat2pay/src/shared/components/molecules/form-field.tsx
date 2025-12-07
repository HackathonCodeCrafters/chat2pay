import * as React from "react";
import { Input, Label } from "@/shared/components/atoms";
import { cn } from "@/shared/lib/utils";

interface FormFieldProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label: string;
  error?: string;
  hint?: string;
}

const FormField = React.forwardRef<HTMLInputElement, FormFieldProps>(
  ({ label, error, hint, className, id, ...props }, ref) => {
    const generatedId = React.useId();
    const inputId = id || generatedId;

    return (
      <div className={cn("space-y-2", className)}>
        <Label htmlFor={inputId} className={cn(error && "text-destructive")}>
          {label}
        </Label>
        <Input ref={ref} id={inputId} error={!!error} {...props} />
        {error && <p className="text-sm text-destructive">{error}</p>}
        {hint && !error && <p className="text-sm text-muted-foreground">{hint}</p>}
      </div>
    );
  }
);
FormField.displayName = "FormField";

export { FormField };
