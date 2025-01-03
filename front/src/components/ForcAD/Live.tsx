import { centrifugeWSURL } from '@/config';
import { useServices, useTeams } from '@/lib/hooks/scoreboard';
import { AttackNotification_Batch } from '@/proto/receiver/receiver';
import { Box } from '@mui/material';
import { styled } from '@mui/material/styles';
import { useWindowSize } from '@uidotdev/usehooks';
import { Centrifuge } from 'centrifuge';
import { useEffect, useMemo, useState } from 'react';
import { FixedSizeList } from 'react-window';

interface Event {
  victimName: string;
  attackerName: string;
  serviceName: string;
  attackerDelta: number;
  victimDelta: number;
}

const EventText = styled('div')({
  fontFamily: 'Roboto Mono',
  fontSize: '0.9rem',
  color: '#00ff00',
});

const Highlight = styled('span')({
  fontWeight: 'bold',
  color: '#ffff00',
});

export default function Live() {
  const [events, setEvents] = useState<Event[]>([]);
  const teams = useTeams();
  const services = useServices();

  const { height: windowHeight } = useWindowSize();

  const teamMap = useMemo(() => {
    return new Map(teams?.map((team) => [team.id, team]));
  }, [teams]);

  const serviceMap = useMemo(() => {
    return new Map(services?.map((service) => [service.id, service]));
  }, [services]);

  useEffect(() => {
    if (!teamMap || !serviceMap) {
      return;
    }

    setEvents(
      Array.from({ length: 1000 }, () => ({
        victimName:
          teamMap.get(
            Array.from(teamMap.keys())[
              Math.floor(Math.random() * teamMap.size)
            ],
          )?.name ?? 'Unknown',
        attackerName:
          teamMap.get(
            Array.from(teamMap.keys())[
              Math.floor(Math.random() * teamMap.size)
            ],
          )?.name ?? 'Unknown',
        serviceName:
          serviceMap.get(
            Array.from(serviceMap.keys())[
              Math.floor(Math.random() * serviceMap.size)
            ],
          )?.name ?? 'Unknown',
        attackerDelta: Math.floor(Math.random() * 100),
        victimDelta: Math.floor(Math.random() * 100),
      })),
    );

    const centrifuge = new Centrifuge(centrifugeWSURL);

    const sub = centrifuge.newSubscription('attacks');

    sub.on('publication', (ctx) => {
      const attackBatch = AttackNotification_Batch.fromJSON(ctx.data);

      const newEvents = attackBatch.attacks.map((attack) => ({
        victimName: teamMap.get(attack.victimId)?.name ?? 'Unknown',
        attackerName: teamMap.get(attack.attackerId)?.name ?? 'Unknown',
        serviceName: serviceMap.get(attack.serviceId)?.name ?? 'Unknown',
        attackerDelta: attack.attackerDelta,
        victimDelta: attack.victimDelta,
      }));

      setEvents((prev) => [...newEvents, ...prev].slice(0, 50));
    });

    sub.subscribe();
    centrifuge.connect();

    return () => {
      sub.unsubscribe();
      centrifuge.disconnect();
    };
  }, [serviceMap, teamMap]);

  const Row = ({
    index,
    style,
  }: {
    index: number;
    style: React.CSSProperties;
  }) => {
    const event = events[index];
    return (
      <EventText style={style}>
        <Highlight>{event.attackerName}</Highlight> attacked{' '}
        <Highlight>{event.victimName}</Highlight>'s service{' '}
        <Highlight>{event.serviceName}</Highlight> and got{' '}
        <Highlight>{event.attackerDelta.toFixed(2)}</Highlight> points
      </EventText>
    );
  };

  if (!windowHeight) {
    return <div>Loading...</div>;
  }

  return (
    <Box sx={{ height: '100vh', width: '100vw', backgroundColor: '#000' }}>
      <FixedSizeList
        height={windowHeight}
        width="100%"
        itemCount={events.length}
        itemSize={20}
      >
        {Row}
      </FixedSizeList>
    </Box>
  );
}
