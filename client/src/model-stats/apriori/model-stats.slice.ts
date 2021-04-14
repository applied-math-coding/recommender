import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { handleError, handleFetchError } from '../../common/error-service';
import { FetchState } from '../../domain/fetch-state';
import { RuleStatistic } from '../../domain/rule-statistic';

export interface ModelStatsState {
  ruleStatistics: RuleStatistic[];
  statsLoading: FetchState;
}

const initialState: ModelStatsState = {
  ruleStatistics: [],
  statsLoading: FetchState.idle
}

export const fetchModelStats = createAsyncThunk('fetchModelStats', async (): Promise<RuleStatistic[]> => {
  try {
    return await fetch('/api/model/apriori/stats')
      .then(r => handleFetchError(r))
      .then(r => r.json());
  } catch (e) {
    handleError(e)
    return [];
  }
});

export const modelStatsSlice = createSlice({
  name: 'modelStats',
  initialState,
  reducers: {
    changeStatsLoading: (state, action: PayloadAction<FetchState>) => {
      state.statsLoading = action.payload;
    }
  },
  extraReducers: builder => {
    builder.addCase(fetchModelStats.pending, (state, action) => {
      state.statsLoading = FetchState.loading;
    });
    builder.addCase(fetchModelStats.fulfilled, (state, action) => {
      state.statsLoading = FetchState.succeeded;
      state.ruleStatistics = action.payload;
    });
  }
});

export const { changeStatsLoading } = modelStatsSlice.actions;
export default modelStatsSlice.reducer;
