import { Status } from '@/proto/checker/checker';
import { Theme } from '@mui/material';

export default function statusColor(theme: Theme, status: Status) {
  switch (status) {
    case Status.STATUS_UP:
      return theme.palette.statusUp;
    case Status.STATUS_DOWN:
      return theme.palette.statusDown;
    case Status.STATUS_CORRUPT:
      return theme.palette.statusCorrupt;
    case Status.STATUS_MUMBLE:
      return theme.palette.statusMumble;
    case Status.STATUS_CHECK_FAILED:
      return theme.palette.statusCheckFailed;
    default:
      return theme.palette.statusUnspecified;
  }
}
