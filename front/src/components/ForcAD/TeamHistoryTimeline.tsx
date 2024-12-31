import { createProtoTransform } from '@/lib/clients/common';
import statusColor from '@/lib/styles/statusColor';
import {
  actionToJSON,
  Execution,
  Execution_Batch,
  Status,
  statusToJSON,
} from '@/proto/checker/checker';
import { Service, Service_Batch } from '@/proto/data/services/services';
import Timeline from '@mui/lab/Timeline';
import TimelineConnector from '@mui/lab/TimelineConnector';
import TimelineContent from '@mui/lab/TimelineContent';
import TimelineDot from '@mui/lab/TimelineDot';
import TimelineItem from '@mui/lab/TimelineItem';
import TimelineOppositeContent from '@mui/lab/TimelineOppositeContent';
import TimelineSeparator from '@mui/lab/TimelineSeparator';
import { Typography } from '@mui/material';
import axios from 'axios';
import { useEffect, useState } from 'react';

interface Props {
  teamID: string;
}

export default function TeamHistoryTimeline({ teamID }: Props) {
  const [services, setServices] = useState<Map<string, Service> | null>(null);
  const [history, setHistory] = useState<Execution[]>([]);

  useEffect(() => {
    async function fetchHistory() {
      const { data } = await axios.get<Execution_Batch>(
        `/teams/${teamID}/history`,
        {
          transformResponse: createProtoTransform(Execution_Batch),
        },
      );
      setHistory(data.executions);
    }
    fetchHistory();
  }, [teamID]);

  useEffect(() => {
    async function fetchServices() {
      const { data } = await axios.get<Service_Batch>(`/services`, {
        transformResponse: createProtoTransform(Service_Batch),
      });
      setServices(
        new Map(data.services.map((service) => [service.id, service])),
      );
    }
    fetchServices();
  }, []);

  if (!history || !services) {
    return <div>Loading...</div>;
  }

  return (
    <Timeline position="right">
      {history.map((execution) => (
        <TimelineItem key={execution.createdAt?.toString()}>
          <TimelineOppositeContent color="text.secondary">
            {execution.createdAt?.toLocaleString()}
          </TimelineOppositeContent>
          <TimelineSeparator>
            <TimelineDot
              sx={(theme) => ({
                backgroundColor: statusColor(theme, execution.status).main,
              })}
            />
            <TimelineConnector />
          </TimelineSeparator>
          <TimelineContent>
            <Typography>
              {services.get(execution.serviceId)?.name}:{' '}
              {execution.status == Status.STATUS_UP ? 'OK' : execution.public}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {actionToJSON(execution.action).replace('ACTION_', '')}:{' '}
              {statusToJSON(execution.status).replace('STATUS_', '')}
            </Typography>
          </TimelineContent>
        </TimelineItem>
      ))}
    </Timeline>
  );
}
