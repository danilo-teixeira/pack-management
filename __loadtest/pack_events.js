import { check, group } from 'k6';
import http from 'k6/http';
import { Trend } from 'k6/metrics';
import { generateRandomString, generateRandomDate } from './helpers.js';
import { defaultSummary, summaryTrendStats, baseURL } from './config.js';

function createPack() {
  const payload = {
    description: "Livros para entrega " + generateRandomString(10),
    sender: "Loja ABC " + generateRandomString(5),
    recipient: "João Silva " + generateRandomString(5),
    estimated_delivery_date: generateRandomDate(),
  };

  const res = http.post(
    baseURL + '/packs',
    JSON.stringify(payload),
    {
      headers: {
        'Content-Type': 'application/json',
      },
    }
  );

  if (res.status !== 201) {
    console.error(`Failed to create pack: ${res.status} ${res.body}`);
  }

  return res.json();
}

const createEventThread = new Trend('_create_event_duration');

export const options = {
  summaryTrendStats,
  scenarios: {
    events: {
      executor: 'constant-arrival-rate',
      startTime: '5s',
      gracefulStop: '5s',
      preAllocatedVUs: 10,
      timeUnit: '1m',
      maxVUs: 100,
      rate: 10000,
      duration: '5m'
    }
  }
};

export function setup() {
  const packIds = Array(10).fill(null).map(() => createPack().id);
  console.log(`Created packs: ${packIds.join('; ')}`);

  return { packIds };
}

export function handleSummary(data) {
  return defaultSummary(data);
}

export default function (data) {
  group('create events', () => {
    const packID = data.packIds[Math.floor(Math.random() * data.packIds.length)];
    if (!packID) {
      console.error('No pack IDs available');
      return;
    }

    const payload = {
      pack_id: packID,
      description: "Pacote chegou ao centro de distribuição",
      location: "Centro de Distribuição São Paulo",
      date: "2025-01-20T15:13:59Z"
    };

    const res = http.post(
      baseURL + '/pack_events',
      JSON.stringify(payload),
      {
        headers: {
          'Content-Type': 'application/json',
        },
      }
    );

    createEventThread.add(res.timings.duration);

    check(res, {
      'status is 204': r => r.status === 204
    });
  });
}
