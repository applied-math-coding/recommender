import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { handleError, handleFetchError } from '../../common/error-service';
import { FetchState } from '../../domain/fetch-state';
import { Recommendation } from '../../domain/recommendation';
import { Rule } from '../../domain/rule';

export interface ModelUseState {
  exampleRules: Rule[];
  recommendations: Recommendation[];
  loadingExamples: FetchState;
  loadingRecommendations: FetchState;
}

const initialState: ModelUseState = {
  exampleRules: [],
  recommendations: [],
  loadingExamples: FetchState.idle,
  loadingRecommendations: FetchState.idle
}

export const fetchExampleRules = createAsyncThunk(
  'fetchExampleRules',
  async (): Promise<Rule[]> => {
    try {
      return await fetch('/api/model/apriori/examples')
        .then(r => handleFetchError(r))
        .then(r => r.json());
    } catch (e) {
      handleError(e);
      return [];
    }
  });

export const fetchRecommendations = createAsyncThunk(
  'fetchRecommendations',
  async (items: string[]): Promise<Recommendation[]> => {
    try {
      return await fetch('/api/model/apriori/apply', {
        method: 'POST',
        body: JSON.stringify(items)
      }).then(r => handleFetchError(r))
        .then(r => r.json());
    } catch (e) {
      handleError(e);
      return [];
    }
  });

export const modelUseSlice = createSlice({
  name: 'modelUse',
  initialState,
  reducers: {
    changeLoadingExamples: (state, action: PayloadAction<FetchState>) => {
      state.loadingExamples = action.payload;
    }
  },
  extraReducers: builder => {
    builder.addCase(fetchRecommendations.pending, (state, action) => {
      state.loadingRecommendations = FetchState.loading;
    });
    builder.addCase(fetchRecommendations.fulfilled, (state, action) => {
      state.loadingRecommendations = FetchState.succeeded;
      state.recommendations = action.payload;
    });
    builder.addCase(fetchExampleRules.pending, (state, action) => {
      state.loadingExamples = FetchState.loading;
    });
    builder.addCase(fetchExampleRules.fulfilled, (state, action) => {
      state.loadingExamples = FetchState.succeeded;
      state.exampleRules = action.payload;
    });
  }
});

export const { changeLoadingExamples } = modelUseSlice.actions;
export default modelUseSlice.reducer;
