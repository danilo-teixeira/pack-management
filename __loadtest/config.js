import { htmlReport } from 'https://raw.githubusercontent.com/benc-uk/k6-reporter/2.4.0/dist/bundle.js';
import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.1/index.js';

export const summaryTrendStats = ['avg', 'min', 'med', 'max', 'p(50)', 'p(90)', 'p(95)', 'p(99)', 'p(99.9)', 'count'];
export const baseURL = 'http://localhost:3300';

export function defaultSummary(data) {
  return {
    'report.html': htmlReport(data),
    'report.json': JSON.stringify(data),
    stdout: textSummary(data, { indent: ' ', enableColors: true })
  };
}
