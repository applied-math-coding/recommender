import { ProgressStateType } from './progress-state-type.enum';

export interface ProgressMessage {
  message: string;
  progress: number;
  state: ProgressStateType;
  showProgressBar: boolean;
}
