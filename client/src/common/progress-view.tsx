import { ProgressMessage } from '../domain/progress-message';
import { Progress } from 'antd';
import { ProgressStateType } from '../domain/progress-state-type.enum';

export default function ProgressView(props: { progressMessage?: ProgressMessage }) {
  function translateProgress(p: number = 0): number {
    return Math.round(p * 100);
  }

  function translateProgressState(p: ProgressStateType | undefined): 'success' | 'normal' | 'exception' | 'active' | undefined {
    return ProgressStateType.Error === p ? 'exception' : 'active';
  }

  return (
    props.progressMessage ?
      <div className="progress">
        <div className="message">{props.progressMessage.message}</div>
        {
          props.progressMessage.showProgressBar &&
          <Progress percent={translateProgress(props.progressMessage.progress)}
            status={translateProgressState(props.progressMessage.state)} />
        }
      </div> : <></>
  )
}