import { ChartData } from '../domain/chart-data';

export const listToChartData = (d: number[]): ChartData[] =>
  d.map(e => ({ value: e }));
