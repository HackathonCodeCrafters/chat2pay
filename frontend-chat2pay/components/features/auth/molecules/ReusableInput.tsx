import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import React from "react";

interface InputFieldsProps {
  label?: string;
  type: string;
  value: string | number;
  placeholder?: string;
}

const ReusableInput: React.FC<InputFieldsProps> = ({
  label,
  type,
  value,
  placeholder,
}) => {
  return (
    <div>
      <Label>{label}</Label>
      <Input type={type} placeholder={placeholder} value={value} />
    </div>
  );
};

export default ReusableInput;
