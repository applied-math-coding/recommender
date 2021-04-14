import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { handleError, handleFetchError } from '../../common/error-service';
import { CosineStatistic } from '../../domain/cosine-statistic';
import { FetchState } from '../../domain/fetch-state';

export interface CosineModelStatsState {
  cosineStatistic?: CosineStatistic;
  statsLoading: FetchState;
}

const initialState: CosineModelStatsState = {
   statsLoading: FetchState.idle
}

export const fetchCosineModelStats = createAsyncThunk('fetchCosineModelStats', async (): Promise<CosineStatistic> => {
  try {
    return await fetch('/api/model/cosine/stats')
      .then(r => handleFetchError(r))
      .then(r => r.json());
  } catch (e) {
    handleError(e)
    return null;
  }
});

export const cosineModelStatsSlice = createSlice({
  name: 'cosineModelStats',
  initialState,
  reducers: {
    changeStatsLoading: (state, action: PayloadAction<FetchState>) => {
      state.statsLoading = action.payload;
    }
  },
  extraReducers: builder => {
    builder.addCase(fetchCosineModelStats.pending, (state, action) => {
      state.statsLoading = FetchState.loading;
    });
    builder.addCase(fetchCosineModelStats.fulfilled, (state, action) => {
      state.statsLoading = FetchState.succeeded;
      state.cosineStatistic = action.payload;
    });
  }
});

export const { changeStatsLoading } = cosineModelStatsSlice.actions;
export default cosineModelStatsSlice.reducer;
