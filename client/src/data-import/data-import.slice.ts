import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { ProgressMessage } from '../domain/progress-message';

export interface DataImportState {
  progressMessage?: ProgressMessage;
  uploadProgress?: ProgressMessage;
  support: number;
  modelProcessId?: number;
  dataImportProcessId?: number;
}

const initialState: DataImportState = {
  support: 100
};

export const dataImportSlice = createSlice({
  name: 'dataImport',
  initialState,
  reducers: {
    changeProgressMessage: (state, action: PayloadAction<ProgressMessage>) => {
      state.progressMessage = action.payload;
    },
    changeUploadProgress: (state, action: PayloadAction<ProgressMessage>) => {
      state.uploadProgress = action.payload;
    },
    changeSupport: (state, action: PayloadAction<number>) => {
      state.support = action.payload;
    },
    changeModelProcessId: (state, action: PayloadAction<number>) => {
      state.modelProcessId = action.payload;
    },
    changeDataImportProcessId: (state, action: PayloadAction<number>) => {
      state.dataImportProcessId = action.payload;
    }
  }
});

export const {
  changeProgressMessage,
  changeUploadProgress,
  changeSupport,
  changeModelProcessId,
  changeDataImportProcessId
} = dataImportSlice.actions;

export default dataImportSlice.reducer;