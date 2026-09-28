import NetworkCheckIcon from "@mui/icons-material/NetworkCheck";
import CodeIcon from "@mui/icons-material/Code";
import DataObjectIcon from "@mui/icons-material/DataObject";
import AccessTimeIcon from "@mui/icons-material/AccessTime";
import FingerprintIcon from "@mui/icons-material/Fingerprint";
import ScheduleIcon from "@mui/icons-material/Schedule";
import LanIcon from "@mui/icons-material/Lan";
import KeyIcon from "@mui/icons-material/Key";
import CasinoIcon from "@mui/icons-material/Casino";
import PinIcon from "@mui/icons-material/Pin";
import CompareIcon from "@mui/icons-material/Compare";
import LockIcon from "@mui/icons-material/Lock";
import type { ElementType } from "react";

export const toolIconMap: Record<string, ElementType> = {
  NetworkCheck: NetworkCheckIcon,
  Code: CodeIcon,
  RegularExpression: CodeIcon,
  DataObject: DataObjectIcon,
  AccessTime: AccessTimeIcon,
  Fingerprint: FingerprintIcon,
  Schedule: ScheduleIcon,
  Lan: LanIcon,
  Key: KeyIcon,
  Casino: CasinoIcon,
  Pin: PinIcon,
  Compare: CompareIcon,
  Lock: LockIcon,
};

export function getToolIcon(icon: string): ElementType {
  return toolIconMap[icon] || CodeIcon;
}
