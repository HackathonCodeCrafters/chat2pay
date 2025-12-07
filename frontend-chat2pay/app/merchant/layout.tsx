import { MerchantAuthProvider } from "@/features/merchant";

export default function MerchantLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <MerchantAuthProvider>{children}</MerchantAuthProvider>;
}
