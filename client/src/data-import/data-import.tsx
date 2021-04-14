import { Button, InputNumber, Select, UploadProps } from 'antd';
import { useEffect } from 'react';
import { useHistory } from 'react-router';
import { changePageTitle, changeSelectedMenuKey, changeSelectedModelType } from '../app.slice';
import { useAppDispatch, useAppSelector } from '../hooks';
import {
  changeModelProcessId, changeProgressMessage,
  changeSupport, changeUploadProgress, changeDataImportProcessId
} from './data-import.slice';
import './data-import.scss';
import { Upload, message } from 'antd';
import { InboxOutlined } from '@ant-design/icons';
import { ProgressMessage } from '../domain/progress-message';
import { ProgressStateType } from '../domain/progress-state-type.enum';
import { MenuKey } from '../domain/menu-key.enum';
import ProgressView from '../common/progress-view';
import { handleError, handleFetchError } from '../common/error-service';
import { ModelType } from '../domain/model-type.enum';
import { Option } from 'antd/lib/mentions';

export default function DataImport() {
  const PAGE_TITLE = 'Data Import';
  const dispatch = useAppDispatch()
  const progressMessage = useAppSelector(state => state.dataImport.progressMessage);
  const uploadProgress = useAppSelector(state => state.dataImport.uploadProgress);
  const pageTitle = useAppSelector(s => s.app.pageTitle);
  const modelProcessId = useAppSelector(s => s.dataImport.modelProcessId);
  const dataImportProcessId = useAppSelector(s => s.dataImport.dataImportProcessId);
  const selectedModelType = useAppSelector(s => s.app.selectedModelType);
  const history = useHistory();
  const { Dragger } = Upload;
  const support = useAppSelector(s => s.dataImport.support);

  useEffect(() => {
    if (pageTitle !== PAGE_TITLE) {
      // this triggers a full re-render
      dispatch(changePageTitle(PAGE_TITLE));
      dispatch(changeSelectedMenuKey(MenuKey.DataImport));
    }
  });

  const fileUploadProps: UploadProps = {
    name: 'file',
    multiple: false,
    action: '/api/data',
    showUploadList: false,
    onChange(info) {
      const { file: { status, response: { id } = {} }, event: { percent } = {} } = info;
      if (status === 'done') {
        message.success(`${info.file.name} file uploaded successfully.`);
        subscribeToUploadProgress(id);
        dispatch(changeDataImportProcessId(id));
      } else if (status === 'error') {
        message.error(`${info.file.name} file upload failed.`);
      } else if (percent && percent <= 1) {
        dispatch(
          changeUploadProgress({
            message: 'Uploading File',
            state: ProgressStateType.Running,
            progress: percent,
            showProgressBar: true
          })
        )
      }
    }
  };

  const handleCancelFit = async () => {
    try {
      await fetch(`/api/model/${encodeURIComponent(selectedModelType)}/cancel/${encodeURIComponent(modelProcessId)}`, {
        method: 'DELETE'
      })
        .then(r => handleFetchError(r))
        .then(() => dispatch(changeModelProcessId(null)));
    } catch (e) {
      handleError(e);
    }
  }

  const handleCancelUpload = async () => {
    try {
      await fetch(`/api/data/cancel/${encodeURIComponent(dataImportProcessId)}`, {
        method: 'DELETE'
      })
        .then(r => handleFetchError(r))
        .then(() => dispatch(changeDataImportProcessId(null)));
    } catch (e) {
      handleError(e);
    }
  }

  const handleFitModel = async (support: number = 3) => {
    try {
      await fetch(`/api/model/${encodeURIComponent(selectedModelType)}/create`, {
        method: 'POST',
        body: JSON.stringify({
          support
        })
      })
        .then(r => handleFetchError(r))
        .then(r => r.json())
        .then(({ id }) => {
          subscribeToProgress(id);
          dispatch(changeModelProcessId(id));
        });
    } catch (e) {
      handleError(e);
    }
  }

  function subscribeToProgress(id: number) {
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
    const url = `${protocol}://${window.location.host}/api/model/${encodeURIComponent(selectedModelType)}/progress/${encodeURIComponent(id)}`;
    const ws = new WebSocket(url);
    ws.addEventListener('error', e => handleError(e));
    ws.addEventListener('message', ({ data }) => {
      const pm: ProgressMessage = JSON.parse(data);
      dispatch(changeProgressMessage(pm));
      if (ProgressStateType.Error === pm.state) {
        ws.close();
      }
      if (ProgressStateType.Finished === pm.state) {
        ws.close();
        history.push(`/model-stats/${selectedModelType}`);
      }
    });
  }

  function subscribeToUploadProgress(id: number) {
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
    const url = `${protocol}://${window.location.host}/api/data/progress/${encodeURIComponent(id)}`;
    const ws = new WebSocket(url);
    ws.addEventListener('error', e => handleError(e));
    ws.addEventListener('message', ({ data }) => {
      const pm: ProgressMessage = JSON.parse(data);
      dispatch(changeUploadProgress(pm));
      if ([ProgressStateType.Finished, ProgressStateType.Error].includes(pm.state)) {
        ws.close();
      }
    });
  }

  return (
    <div className="data-import-container">
      <Dragger {...fileUploadProps}>
        <p className="ant-upload-drag-icon">
          <InboxOutlined />
        </p>
        <p className="ant-upload-text">Click or drag file to this area to upload</p>
        <p className="ant-upload-hint">
          Data must be in comma-separated csv-format. Example:<br></br>
          item1,item2,item3<br></br>
          item1,item3<br></br>
          item2,item3,item1
        </p>
      </Dragger>
      {
        uploadProgress && uploadProgress.state !== ProgressStateType.Finished &&
        <>
          <div className="progress-container">
            <ProgressView progressMessage={uploadProgress}></ProgressView>
          </div>
          <Button type="primary" danger size="large" onClick={() => handleCancelUpload()}>Cancel</Button>
        </>
      }
      {
        uploadProgress && uploadProgress.state === ProgressStateType.Finished &&
        <div className="action-container">
          <div style={{ marginBottom: '10px' }}>
            <div>Select Model:</div>
            <Select defaultValue={selectedModelType} style={{ width: 120 }} onChange={e => dispatch(changeSelectedModelType(e))}>
              {
                Object.values(ModelType).map(v =>
                  <Option value={v} key={v}>{v}</Option>
                )
              }
            </Select>
          </div>
          <div style={{ marginBottom: '10px' }}>
            <div>Support:</div>
            <InputNumber min={1} defaultValue={support} onChange={v => dispatch(changeSupport(v))} />
          </div>
          <div style={{ marginTop: '20px' }}>
            <div>
              {
                progressMessage?.state !== ProgressStateType.Running ?
                  <Button type="primary" size="large" onClick={() => handleFitModel(support)}>Fit Model</Button> :
                  <Button type="primary" danger size="large" onClick={() => handleCancelFit()}>Cancel</Button>
              }
            </div>
            {
              !!progressMessage &&
              <div className="progress-container">
                <ProgressView progressMessage={progressMessage}></ProgressView>
              </div>
            }
          </div>
        </div>
      }
    </div>
  )
}
