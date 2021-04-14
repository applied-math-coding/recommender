import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { handleError, handleFetchError } from '../../common/error-service';
import { CosineRecommendation } from '../../domain/cosine-recommendation';
import { FetchState } from '../../domain/fetch-state';

export interface ModelUseState {
  recommendations: CosineRecommendation[];
  loadingRecommendations: FetchState;
}

const initialState: ModelUseState = {
  recommendations: [],
  loadingRecommendations: FetchState.idle
}

export const fetchCosineRecommendations = createAsyncThunk(
  'fetchCosineRecommendations',
  async (items: string[]): Promise<CosineRecommendation[]> => {
    try {
      return await fetch('/api/model/cosine/apply', {
        method: 'POST',
        body: JSON.stringify(items)
      }).then(r => handleFetchError(r))
        .then(r => r.json());
    } catch (e) {
      handleError(e);
      return [];
    }
  });

export const cosineModelUseSlice = createSlice({
  name: 'cosineModelUse',
  initialState,
  reducers: {},
  extraReducers: builder => {
    builder.addCase(fetchCosineRecommendations.pending, (state, action) => {
      state.loadingRecommendations = FetchState.loading;
    });
    builder.addCase(fetchCosineRecommendations.fulfilled, (state, action) => {
      state.loadingRecommendations = FetchState.succeeded;
      state.recommendations = action.payload;
    });
  }
});

export const { } = cosineModelUseSlice.actions;
export default cosineModelUseSlice.reducer;
