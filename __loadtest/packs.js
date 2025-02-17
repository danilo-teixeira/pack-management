import { check, group } from 'k6';
import http from 'k6/http';
import { Trend } from 'k6/metrics';
import { generateRandomString, generateRandomDate } from './helpers.js';
import { defaultSummary, summaryTrendStats, baseURL } from './config.js';

const createPackThread = new Trend('_create_pack_duration');
const updatePackStatusThread = new Trend('_update_pack_status_duration');
const cancelPackThread = new Trend('_cancel_pack_duration');

export const options = {
  summaryTrendStats,
  scenarios: {
    packs: {
      executor: 'constant-arrival-rate',
      startTime: '5s',
      gracefulStop: '5s',
      preAllocatedVUs: 10,
      timeUnit: '1m',
      maxVUs: 100,
      rate: 6000,
      duration: '5m'
    }
  }
};

let state = {
  packToIntransit: [],
  packToCancelled: []
};

function getState() {
  return state;
}

function updateState(newState) {
  state = newState;
}

export function setup() {
  return {}
}

export function handleSummary(data) {
  return defaultSummary(data);
}

export default function (data) {
  group('create packs', () => {
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

    createPackThread.add(res.timings.duration);

    if (res.status === 201) {
      if (Math.random() < 0.5) {
        state.packToIntransit.push(res.json().id);
      } else {
        state.packToCancelled.push(res.json().id);
      }
    }

    check(res, {
      'status is 201': r => r.status === 201
    });

    updateState(state);
  });

  group('cancel pack', () => {
    const state = getState();
    const packId = state.packToCancelled.pop();
    if (!packId) {
      return;
    }

    const res = http.post(
      baseURL + '/packs/' + packId + '/cancel',
      null,
      {
        headers: {
          'Content-Type': 'application/json',
        },
      }
    );

    cancelPackThread.add(res.timings.duration);

    check(res, {
      'status is 200': r => r.status === 200
    });

    updateState(state);
  });

  group('update pack status', () => {
    const state = getState();
    const packId = state.packToIntransit.pop();
    if (!packId) {
      return;
    }

    const res = http.patch(
      baseURL + '/packs/' + packId,
      JSON.stringify({ status: 'IN_TRANSIT' }),
      {
        headers: {
          'Content-Type': 'application/json',
        },
      }
    );

    updatePackStatusThread.add(res.timings.duration);

    check(res, {
      'status is 200': r => r.status === 200
    });

    // Update status to DELIVERED
    const res2 = http.patch(
      baseURL + '/packs/' + packId,
      JSON.stringify({ status: 'DELIVERED' }),
      {
        headers: {
          'Content-Type': 'application/json',
        },
      }
    );

    updatePackStatusThread.add(res2.timings.duration);

    check(res2, {
      'status is 200': r => r.status === 200
    });

    updateState(state);
  });
}
